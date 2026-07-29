.PHONY: build test doctor install release-snapshot clean

VERSION ?= 0.4.1
LDFLAGS := -X github.com/NewtSpeak/Newt-Agent/internal/cmd.Version=$(VERSION) \
	-X github.com/NewtSpeak/Newt-Agent/internal/cmd.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo dev) \
	-X github.com/NewtSpeak/Newt-Agent/internal/cmd.BuildDate=$$(date -u +%Y-%m-%dT%H:%M:%SZ)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/newt ./cmd/newt

test:
	go test ./... -count=1

doctor: build
	./bin/newt doctor || true

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/newt

# 需要 goreleaser（可选本地快照）
release-snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -rf bin dist newt newt.exe
