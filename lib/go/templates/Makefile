.PHONY: test
test:
	CGO_ENABLED=0 go test -tags=no_cgo ./...

.PHONY: generate
generate:
	go generate

.PHONY: check-tidy
check-tidy: generate
	go mod tidy
	git diff --exit-code

.PHONY: ci
ci: check-tidy test
