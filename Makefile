PROJECT_NAME := "SBWeb"
PKG := "bmstu.codes/developers34/SBWeb"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# TODO cover and coverhtml
.PHONY: all dep build clean test lint

all: build

lint:
	@golint -set_exit_status ${PKG_LIST}

test:
	@go test ${PKG_LIST}

coverage:
	go test -cover ${PKG_LIST}

coverage_html:
	@mkdir coverage
	@go test -coverprofile=./coverage/coverage.cov ${PKG_LIST}
	@go tool cover -html=./coverage/coverage.cov -o ./coverage/coverage.html
	@rm ./coverage/coverage.cov

race: dep
	@go test -race ${PKG_LIST}

msan: dep
	@go test -msan ${PKG_LIST}

dep:
	@go get -v -d ./...

build: dep
	@go build -i -v $(PKG)

clean:
	@rm -f $(PROJECT_NAME)