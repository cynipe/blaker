export GO111MODULE := on

AWS_PROFILE :=
CMD         := echo hey

.PHONY: deps run test fmt
deps:
	go mod vendor

run:
	go run cmd/blaker/main.go --profile $(AWS_PROFILE) $(CMD)

test:
	go test -v -tags=unit $$(go list ./... | grep -v '/vendor/')

fmt:
	go fmt ./...

build: dist/blaker
dist/%: cmd/%/main.go $(shell find . -name *.go) dist fmt
		GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -mod=vendor -o $@ $<
dist:
		mkdir -p $@
