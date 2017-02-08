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

build: clean generate deps format compile

set-build:
	@if [ ! -d $(shell pwd)/.gopath/src/github.com/linuxisnotunix ]; then mkdir -p $(shell pwd)/.gopath/src/github.com/linuxisnotunix; fi
	@if [ ! -d $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot/ ]; then ln -s $(shell pwd) $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot; fi

compile: set-build
	@echo -e "$(OK_COLOR)==> Building...$(NO_COLOR)"
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath go build -o ${APP_NAME} -v -ldflags "$(LDFLAGS)"
#	@ln -s $(shell pwd) $(shell pwd)/vendor/github.com/linuxisnotunix/Vulnerobot 2&> /dev/null || echo -e "$(OK_COLOR)==> Vendor linking allready exist.$(NO_COLOR)"
#mkdir -p $(shell pwd)/.gopath/src/github.com/linuxisnotunix
#ln -s $(shell pwd) $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot || echo -e "$(OK_COLOR)==> GOPATH linking allready exist.$(NO_COLOR)"
#cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot
#GOPATH=$(shell pwd)/.gopath

generate: set-build dev-deps
	$(GOPATH)/bin/go-embed -input public/ -output modules/assets/main.go

	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://nvd.nist.gov/schema/nvd-cve-feed_2.0.xsd
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://nvd.nist.gov/schema/vulnerability_0.4.1.xsd

	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://nvd.nist.gov/schema/cce_0.1.xsd
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://nvd.nist.gov/schema/cve_0.1.1.xsd
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://nvd.nist.gov/schema/cvss-v2_0.2.xsd
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://nvd.nist.gov/schema/patch_0.1.xsd
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://nvd.nist.gov/schema/scap-core_0.1.xsd
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://scap.nist.gov/schema/cpe/2.2/cpe-language_2.2a.xsd
	cd $(shell pwd)/.gopath/src/github.com/linuxisnotunix/Vulnerobot && GOPATH=$(shell pwd)/.gopath $(GOPATH)/bin/xsd-makepkg -goinst=false --basepath github.com/linuxisnotunix/Vulnerobot/modules/models/xsd/ -uri https://www.w3.org/2009/01/xml.xsd
	echo "type XsdtString xsdt.String" >> modules/models/xsd/nvd.nist.gov/schema/scap-core_0.1.xsd_go/scap-core_0.1.xsd.go
	echo "type XsdtString xsdt.String" >> modules/models/xsd/scap.nist.gov/schema/cpe/2.2/cpe-language_2.2a.xsd_go/cpe-language_2.2a.xsd.go

#rm -R modules/models/xsd/
#find ./modules/models/scap.nist.gov -type f -name "*.go" -exec sed -i 's/XsdtString/xsdt.String/g' {} \;

release: clean generate deps format
	@mkdir build || echo -e "$(WARN_COLOR)==> Build dir (/build) allready exist.$(NO_COLOR)"
	@echo -e "$(OK_COLOR)==> Building for linux 32 ...$(NO_COLOR)"
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o build/${APP_NAME}-linux-386 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-386 || upx-ucl --brute  build/${APP_NAME}-linux-386 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for linux 64 ...$(NO_COLOR)"
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}-linux-amd64 -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-amd64 || upx-ucl --brute  build/${APP_NAME}-linux-amd64 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for linux arm ...$(NO_COLOR)"
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 CC="armv6l-unknown-linux-gnueabihf-gcc" go build -o build/${APP_NAME}-linux-armv6 -ldflags "$(LDFLAGS_RELEASE)"
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
	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC="i686-w64-mingw32-gcc" go build -o build/${APP_NAME}-win-386.exe -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-win-386.exe || upx-ucl --brute  build/${APP_NAME}-win-386.exe || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for win64 ...$(NO_COLOR)"
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="x86_64-w64-mingw32-gcc" go build -o build/${APP_NAME}-win-amd64.exe -ldflags "$(LDFLAGS_RELEASE)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-win-amd64.exe || upx-ucl --brute  build/${APP_NAME}-win-amd64.exe || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Archiving ...$(NO_COLOR)"
	@tar -zcvf $(ARCHIVE) LICENSE README.md build/

clean:
	@if [ -x ${APP_NAME} ]; then rm ${APP_NAME}; fi
	@if [ -x Vulnerobot ]; then rm Vulnerobot; fi #go build bin
	@if [ -d build ]; then rm -R build; fi
	@if [ -d $(shell pwd)/vendor/github.com/linuxisnotunix/Vulnerobot ]; then rm $(shell pwd)/vendor/github.com/linuxisnotunix/Vulnerobot; fi

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
	go test -v -race -covermode=atomic  . ./cmd/... ./modules/...
#	go test -v -race -coverprofile=coverage.out -covermode=atomic ./modules/tools

#TODO cover with gocovmerge
#go tool cover -html=coverage.out -o coverage.html
#go test -v -race -coverprofile=coverage.out -covermode=atomic . ./cmd/... ./modules/...
docs: dev-deps
	@echo -e "$(OK_COLOR)==> Serving docs at http://localhost:$(DOC_PORT).$(NO_COLOR)"
	$(GOPATH)/bin/godoc -http=:$(DOC_PORT)

lint: dev-deps
	$(GOPATH)/bin/gometalinter --deadline=5m --concurrency=2 --vendor ./...

dev-deps:
	@echo -e "$(OK_COLOR)==> Installing developement dependencies...$(NO_COLOR)"
	@go get github.com/nsf/gocode
	@go get golang.org/x/tools/cmd/goimports
	@go get golang.org/x/tools/cmd/godoc
	@go get github.com/alecthomas/gometalinter
	@go get github.com/metaleap/go-xsd/xsd-makepkg
	@go get github.com/pyros2097/go-embed
	@go get github.com/dpw/vendetta #Vendoring
	@$(GOPATH)/bin/gometalinter --install > /dev/null
#@go get -u https://github.com/wadey/gocovmerge #Test coverage

update-dev-deps:
	@echo -e "$(OK_COLOR)==> Installing/Updating developement dependencies...$(NO_COLOR)"
	go get -u github.com/nsf/gocode
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/tools/cmd/godoc
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/metaleap/go-xsd/xsd-makepkg
	go get -u github.com/pyros2097/go-embed
	go get -u github.com/dpw/vendetta #Vendoring
	$(GOPATH)/bin/gometalinter --install --update
#go get -u https://github.com/wadey/gocovmerge #Test coverage

deps:
	@echo -e "$(OK_COLOR)==> Installing dependencies ...$(NO_COLOR)"
	git submodule update --init --recursive
	$(GOPATH)/bin/vendetta -n github.com/linuxisnotunix/Vulnerobot
#@go get -d -v ./...

update-deps: dev-deps
	@echo "$(OK_COLOR)==> Updating all dependencies ...$(NO_COLOR)"
	$(GOPATH)/bin/vendetta -n github.com/linuxisnotunix/Vulnerobot -u -p
#@go get -d -v -u ./...

done:
	@echo -e "$(OK_COLOR)==> Done.$(NO_COLOR)"

.PHONY: all build compile clean compress format test docs lint dev-deps deps update-deps update-dev-deps done
