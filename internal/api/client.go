package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	ServerURL string
	HTTP      *http.Client
	Token     string
}

func New(serverURL string) *Client {
	return &Client{
		ServerURL: strings.TrimRight(serverURL, "/"),
		HTTP:      &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) WithToken(token string) *Client {
	cp := *c
	cp.Token = token
	return &cp
}

type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type TokenResponse struct {
	AccessToken      string    `json:"access_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int       `json:"expires_in"`
	RefreshToken     string    `json:"refresh_token"`
	Scope            string    `json:"scope"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	Error            string    `json:"error"`
	ErrorDescription string    `json:"error_description"`
}

type UserInfo struct {
	Sub         string `json:"sub"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	SystemAdmin bool   `json:"system_admin"`
	ClientID    string `json:"client_id"`
	Scope       string `json:"scope"`
}

type APIError struct {
	Status  int
	Code    string
	Message string
	Body    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Body)
}

func (c *Client) RequestDeviceCode(clientID, scope string) (*DeviceCodeResponse, error) {
	body := map[string]string{"client_id": clientID, "scope": scope}
	var out DeviceCodeResponse
	if err := c.postJSON("/oauth/v1/device/code", body, &out, false); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PollToken(clientID, deviceCode string) (*TokenResponse, error) {
	body := map[string]string{
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		"device_code": deviceCode,
		"client_id":   clientID,
	}
	var out TokenResponse
	status, err := c.doJSON(http.MethodPost, "/oauth/v1/token", body, &out, false)
	if err != nil && status == 0 {
		return nil, err
	}
	if out.Error != "" {
		return &out, &APIError{Status: status, Code: out.Error, Message: out.ErrorDescription}
	}
	if status >= 400 {
		return &out, &APIError{Status: status, Code: out.Error, Message: out.ErrorDescription}
	}
	return &out, nil
}

func (c *Client) Refresh(refreshToken, clientID string) (*TokenResponse, error) {
	body := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"client_id":     clientID,
	}
	var out TokenResponse
	if err := c.postJSON("/oauth/v1/token", body, &out, false); err != nil {
		return nil, err
	}
	if out.Error != "" {
		return nil, &APIError{Code: out.Error, Message: out.ErrorDescription}
	}
	return &out, nil
}

func (c *Client) Revoke(token string) error {
	return c.postJSON("/oauth/v1/revoke", map[string]string{"token": token}, nil, false)
}

func (c *Client) UserInfo() (*UserInfo, error) {
	var out UserInfo
	if err := c.getJSON("/oauth/v1/userinfo", &out, true); err != nil {
		return nil, err
	}
	return &out, nil
}

// Gapi 对 /gapi/v1 发起任意方法请求，body 可为 nil；query 为可选查询参数。
func (c *Client) Gapi(method, path string, body any, query map[string]string) (json.RawMessage, error) {
	return c.requestPrefixed("/gapi/v1", method, path, body, query)
}

// Api 对 /api/v1 发起请求（平台管理；需 agent token 含 platform.* + system_admin）。
func (c *Client) Api(method, path string, body any, query map[string]string) (json.RawMessage, error) {
	return c.requestPrefixed("/api/v1", method, path, body, query)
}

func (c *Client) requestPrefixed(prefix, method, path string, body any, query map[string]string) (json.RawMessage, error) {
	full := prefix + path
	if len(query) > 0 {
		q := url.Values{}
		for k, v := range query {
			if v != "" {
				q.Set(k, v)
			}
		}
		full += "?" + q.Encode()
	}
	var raw json.RawMessage
	status, err := c.doJSON(method, full, body, &raw, true)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNoContent || len(raw) == 0 {
		return json.RawMessage(`{"ok":true}`), nil
	}
	return raw, nil
}

func (c *Client) GapiGET(path string, out any) error {
	return c.getJSON("/gapi/v1"+path, out, true)
}

func (c *Client) postJSON(path string, body any, out any, auth bool) error {
	_, err := c.doJSON(http.MethodPost, path, body, out, auth)
	return err
}

func (c *Client) postJSONStatus(path string, body any, out any, auth bool) (int, error) {
	return c.doJSON(http.MethodPost, path, body, out, auth)
}

func (c *Client) getJSON(path string, out any, auth bool) error {
	_, err := c.doJSON(http.MethodGet, path, nil, out, auth)
	return err
}

func (c *Client) doJSON(method, path string, body any, out any, auth bool) (int, error) {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, c.ServerURL+path, reader)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if auth && c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// OAuth token 端点：错误也在 JSON body 的 error 字段
	if path == "/oauth/v1/token" && out != nil && len(data) > 0 {
		_ = json.Unmarshal(data, out)
		return resp.StatusCode, nil
	}

	if resp.StatusCode >= 400 {
		ae := parseAPIError(resp.StatusCode, data)
		return resp.StatusCode, ae
	}
	if out != nil && len(data) > 0 {
		// 支持直接解到 RawMessage
		if rm, ok := out.(*json.RawMessage); ok {
			*rm = append((*rm)[0:0], data...)
			return resp.StatusCode, nil
		}
		if err := json.Unmarshal(data, out); err != nil {
			return resp.StatusCode, fmt.Errorf("解析响应失败: %w; body=%s", err, truncate(string(data), 200))
		}
	}
	return resp.StatusCode, nil
}

func parseAPIError(status int, data []byte) *APIError {
	var wrap struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		ErrorCode string `json:"error"`
		ErrorDesc string `json:"error_description"`
	}
	_ = json.Unmarshal(data, &wrap)
	if wrap.Error.Code != "" {
		return &APIError{Status: status, Code: wrap.Error.Code, Message: wrap.Error.Message, Body: string(data)}
	}
	if wrap.ErrorCode != "" {
		return &APIError{Status: status, Code: wrap.ErrorCode, Message: wrap.ErrorDesc, Body: string(data)}
	}
	return &APIError{Status: status, Message: fmt.Sprintf("HTTP %d", status), Body: string(data)}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
