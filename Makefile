GOPKG ?=	moul.io/logman

include rules.mk

generate:
	GO111MODULE=off go get github.com/campoy/embedmd
	mkdir -p .tmp
	go doc -all > .tmp/godoc.txt
	embedmd -w README.md
	rm -rf .tmp
