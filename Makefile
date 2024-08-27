export GOFLAGS :=-tags=no_cgo
export CGO_ENABLED := 0

.PHONY: test
test:
	$(MAKE) generate -C lib/go
	$(MAKE) test -C lib/go

.PHONY: ci
ci:
	$(MAKE) ci -C lib/go
