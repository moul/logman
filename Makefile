GOPKG ?=	moul.io/logman

include rules.mk

lint:
	cd tool/lint; make
.PHONY: lint
