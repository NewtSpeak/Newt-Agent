.PHONY: build test doctor install release-snapshot clean

VERSION ?= 0.4.1
LDFLAGS := -X github.com/OwlSpeak/Owl-Agent/internal/cmd.Version=$(VERSION) \
	-X github.com/OwlSpeak/Owl-Agent/internal/cmd.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo dev) \
	-X github.com/OwlSpeak/Owl-Agent/internal/cmd.BuildDate=$$(date -u +%Y-%m-%dT%H:%M:%SZ)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/owl ./cmd/owl

test:
	go test ./... -count=1

doctor: build
	./bin/owl doctor || true

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/owl

# 需要 goreleaser（可选本地快照）
release-snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -rf bin dist owl owl.exe
