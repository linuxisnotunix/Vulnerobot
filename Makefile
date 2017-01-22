#Inspired from : https://github.com/littlemanco/boilr-makefile/blob/master/template/Makefile, https://github.com/geetarista/go-boilerplate/blob/master/Makefile, https://github.com/nascii/go-boilerplate/blob/master/GNUmakefile https://github.com/cloudflare/hellogopher/blob/master/Makefile
#PATH=$(PATH:):$(GOPATH)/bin
APP_NAME=vulnerobot

APP_VERSION=$(shell git describe --abbrev=0 || echo "testing")
GIT_HASH=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date -u '+%Y-%m-%d-%H%M-UTC')

ARCHIVE=$(APP_NAME)-$(APP_VERSION)-$(GIT_HASH).tar.gz
#DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
LDFLAGS = \
  -s -w \
  -X main.Version=$(APP_VERSION) -X main.Branch=$(GIT_BRANCH) -X main.Commit=$(GIT_HASH) -X main.BuildTime=$(DATE)

LDFLAGS_RELEASE = -s -w -X main.Version=$(APP_VERSION)

GO15VENDOREXPERIMENT=1
#GOPATH ?= $(GOPATH:):./vendor
DOC_PORT = 6060
#GOOS=linux

ERROR_COLOR=\033[31;01m
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
WARN_COLOR=\033[33;01m

all: build compress done

build: clean deps format compile

compile:
	@echo -e "$(OK_COLOR)==> Building...$(NO_COLOR)"
	go build -o ${APP_NAME} -v -ldflags "$(LDFLAGS)"

release: clean deps format
	@mkdir build || echo -e "$(WARN_COLOR)==> Build dir (/build) allready exist.$(NO_COLOR)"
	@echo -e "$(OK_COLOR)==> Building for linux 32 ...$(NO_COLOR)"
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o build/${APP_NAME}-linux-386 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-386 || upx-ucl --brute  build/${APP_NAME}-linux-386 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for linux 64 ...$(NO_COLOR)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}-linux-amd64 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-amd64 || upx-ucl --brute  build/${APP_NAME}-linux-amd64 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for linux arm ...$(NO_COLOR)"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o build/${APP_NAME}-linux-armv6 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-armv6 || upx-ucl --brute  build/${APP_NAME}-linux-armv6 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for darwin32 ...$(NO_COLOR)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o build/${APP_NAME}-darwin-386 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-darwin-386 || upx-ucl --brute  build/${APP_NAME}-darwin-386 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for darwin64 ...$(NO_COLOR)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}-darwin-amd64 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-darwin-amd64 || upx-ucl --brute  build/${APP_NAME}-darwin-amd64 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for win32 ...$(NO_COLOR)"
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o build/${APP_NAME}-win-386 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-win-386 || upx-ucl --brute  build/${APP_NAME}-win-386 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for win64 ...$(NO_COLOR)"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}-win-amd64 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-win-amd64 || upx-ucl --brute  build/${APP_NAME}-win-amd64 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Archiving ...$(NO_COLOR)"
	@tar -zcvf build/$(ARCHIVE) LICENSE README.md build/

clean:
	@if [ -x ${APP_NAME} ]; then rm ${APP_NAME}; fi
	@if [ -x Vulnerobot ]; then rm Vulnerobot; fi #go build bin
	@if [ -d build ]; then rm -R build; fi

compress:
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute ${APP_NAME} || upx-ucl --brute ${APP_NAME} || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

format:
	@echo -e "$(OK_COLOR)==> Formatting...$(NO_COLOR)"
	go fmt ./ ./modules/... ./cmd/... 
#@find ./ -iname \*.go | grep -v -e "^$$" $(addprefix -e ,/vendor/) | xargs $(GOPATH)/bin/goimports -w
#go fmt ./...

test: deps format
	@echo -e "$(OK_COLOR)==> Running tests...$(NO_COLOR)"
	go vet . ./... || true
	go test -v -race . ./cmd/... ./modules/...

#TODO cover with gocovmerge
#go tool cover -html=coverage.out -o coverage.html
#go test -v -race -coverprofile=coverage.out -covermode=atomic . ./cmd/... ./modules/...
docs: dev-deps
	@echo -e "$(OK_COLOR)==> Serving docs at http://localhost:$(DOC_PORT).$(NO_COLOR)"
	@$(GOPATH)/bin/godoc -http=:$(DOC_PORT)

lint: dev-deps
	$(GOPATH)/bin/gometalinter --deadline=5m --concurrency=2 --vendor ./...

dev-deps:
	@echo -e "$(OK_COLOR)==> Installing developement dependencies...$(NO_COLOR)"
	@go get github.com/nsf/gocode
	@go get golang.org/x/tools/cmd/goimports
	@go get golang.org/x/tools/cmd/godoc
	@go get github.com/alecthomas/gometalinter
	@go get github.com/dpw/vendetta #Vendoring
	@$(GOPATH)/bin/gometalinter --install > /dev/null
#@go get -u https://github.com/wadey/gocovmerge #Test coverage

update-dev-deps:
	@echo -e "$(OK_COLOR)==> Installing/Updating developement dependencies...$(NO_COLOR)"
	go get -u github.com/nsf/gocode
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/tools/cmd/godoc
	#go get -u https://github.com/wadey/gocovmerge #Test coverage
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/dpw/vendetta #Vendoring
	$(GOPATH)/bin/gometalinter --install --update

deps:
	@echo -e "$(OK_COLOR)==> Installing dependencies ...$(NO_COLOR)"
	$(GOPATH)/bin/vendetta
#@go get -d -v ./...

update-deps:
	@echo "$(OK_COLOR)==> Updating all dependencies ...$(NO_COLOR)"
	$(GOPATH)/bin/vendetta -u -p
#@go get -d -v -u ./...

done:
	@echo -e "$(OK_COLOR)==> Done.$(NO_COLOR)"

.PHONY: all build compile clean compress format test docs lint dev-deps deps update-deps update-dev-deps done
