GOPATH=$(shell pwd)
export GOPATH

setup : tools dep
tools :
	@echo ""
	@echo "Installing all tooling and enhancements..."
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/tsenart/deadcode
	go get -u github.com/alecthomas/gocyclo
	go get -u github.com/alexkohler/nakedret
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u github.com/dnephin/govet
	go get -u github.com/golang/lint
	go get -u github.com/gordonklaus/ineffassign
	go get -u github.com/jgautheron/goconst
	go get -u github.com/kisielk/errcheck
	go get -u github.com/kisielk/gotool
	go get -u github.com/mdempsky/maligned
	go get -u github.com/mdempsky/unconvert
	go get -u github.com/mibk/dupl
	go get -u github.com/opennota/check
	go get -u github.com/stripe/safesql
	go get -u github.com/walle/lll/...
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega/...
	go get -u github.com/ramya-rao-a/go-outline
	go get -u github.com/nsf/gocode
	go get -u github.com/newhook/go-symbols
	go get -u github.com/uudashr/gopkgs/cmd/gopkgs
	go get -u honnef.co/go/tools/cmd/gosimple
	go get -u golang.org/x/tools/cmd/guru
	go get -u golang.org/x/tools/cmd/gorename
	go get -u golang.org/x/tools/cmd/gotype
	go get -u github.com/sqs/goreturns
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/rogpeppe/godef
	go get -u github.com/zmb3/gogetdoc
	go get -u golang.org/x/tools/cmd/godoc
	go get -u golang.org/x/lint/golint
	go get -u honnef.co/go/tools/cmd/megacheck
	go get -u github.com/derekparker/delve/cmd/dlv
	go get -u github.com/jgautheron/goconst/cmd/goconst
	go get -u github.com/GoASTScanner/gas/cmd/gas/...
	go get -u github.com/cweill/gotests/...
	go get -u mvdan.cc/interfacer
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u github.com/opennota/check/cmd/aligncheck
	go get -u github.com/opennota/check/cmd/structcheck
	go get -u github.com/opennota/check/cmd/varcheck
	go get -u mvdan.cc/unparam
	go get -u honnef.co/go/tools/cmd/unused
	@echo "Add \$$GOPATH/bin to your \$$PATH to make sure these binaries function correctly!"
dep :
	@echo ""
	@echo "Installing all project dependencies..."
	cd $(GOPATH)/src/service && dep ensure
	@echo "SICK STUFF!"

lint :
	@echo ""
	@echo "Running linters..."
	cd $(GOPATH)/src/service && \
	gometalinter --enable=gofmt --enable=safesql --enable=staticcheck --skip=vendor ./...
	@echo "Looking good!"

fakes :
	@echo ""
	@echo "Generating fresh fakes..."
	cd $(GOPATH)/src/service && go generate \
		./auth ./database ./auth/basic ./auth/token ./identity ./log ./handlers/request \
		./handlers/index

ginkgo :
	@echo ""
	@echo "Running all ginkgo suites..."
	cd $(GOPATH)/src/service && ginkgo -r .

test : fakes ginkgo
	@echo "SUITE SUCCESS!!!"

all : lint test run

run :
	go run src/service/main.go
build : out/app

clean :
	rm -r $(wildcard out/*)

format :
	go fmt ./...
