package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// GeneratePKCE 生成 code_verifier 与 S256 code_challenge。
func GeneratePKCE() (verifier, challenge string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return verifier, challenge, nil
}

func RandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ExchangeAuthCode authorization_code + PKCE 换 token。
func (c *Client) ExchangeAuthCode(clientID, code, verifier, redirectURI string) (*TokenResponse, error) {
	body := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"code":          code,
		"code_verifier": verifier,
		"redirect_uri":  redirectURI,
	}
	var out TokenResponse
	if err := c.postJSON("/oauth/v1/token", body, &out, false); err != nil {
		return nil, err
	}
	if out.Error != "" {
		return nil, &APIError{Code: out.Error, Message: out.ErrorDescription}
	}
	if out.AccessToken == "" {
		return nil, fmt.Errorf("token 响应缺少 access_token")
	}
	return &out, nil
}

// LoopbackResult 本地回调结果。
type LoopbackResult struct {
	Code string
	Err  error
}

// StartLoopback 在 127.0.0.1 随机端口监听 /callback，立即返回 redirectURI。
// 调用方应 defer stop()；从 ch 读取 code。
func StartLoopback(state string) (redirectURI string, stop func(), ch <-chan LoopbackResult, err error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, nil, err
	}
	port := ln.Addr().(*net.TCPAddr).Port
	redirectURI = fmt.Sprintf("http://127.0.0.1:%d/callback", port)
	out := make(chan LoopbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if errCode := q.Get("error"); errCode != "" {
			desc := q.Get("error_description")
			if desc == "" {
				desc = errCode
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<html><body><h3>授权失败</h3><p>%s</p><p>可关闭此窗口。</p></body></html>", desc)
			out <- LoopbackResult{Err: fmt.Errorf("%s", desc)}
			return
		}
		gotState := q.Get("state")
		if state != "" && gotState != state {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "state mismatch")
			out <- LoopbackResult{Err: fmt.Errorf("state 不匹配")}
			return
		}
		c := q.Get("code")
		if c == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "missing code")
			out <- LoopbackResult{Err: fmt.Errorf("缺少 code")}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<html><body><h3>登录成功</h3><p>可以关闭此窗口，返回终端。</p></body></html>")
		out <- LoopbackResult{Code: c}
	})

	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()
	stop = func() { _ = srv.Close(); _ = ln.Close() }
	return redirectURI, stop, out, nil
}

// WaitLoopbackCode 兼容封装：阻塞至拿到 code 或超时。
func WaitLoopbackCode(state string, timeout time.Duration) (code, redirectURI string, err error) {
	redir, stop, ch, err := StartLoopback(state)
	if err != nil {
		return "", "", err
	}
	defer stop()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case res := <-ch:
		if res.Err != nil {
			return "", redir, res.Err
		}
		return res.Code, redir, nil
	case <-timer.C:
		return "", redir, fmt.Errorf("等待浏览器授权超时")
	}
}

// BuildAuthorizeURL 组装 Desktop/Web 授权页 URL。
// clientOrigin 为用户前端基址（可与 API 同域）；为空则用 serverURL。
func BuildAuthorizeURL(clientOrigin, serverURL, clientID, redirectURI, scope, challenge, state string) string {
	base := clientOrigin
	if base == "" {
		base = serverURL
	}
	base = trimRightSlash(base)
	u, err := url.Parse(base + "/oauth/authorize")
	if err != nil {
		return base + "/oauth/authorize"
	}
	q := u.Query()
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", scope)
	q.Set("code_challenge", challenge)
	q.Set("code_challenge_method", "S256")
	if state != "" {
		q.Set("state", state)
	}
	if serverURL != "" {
		q.Set("server", serverURL)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func trimRightSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}
