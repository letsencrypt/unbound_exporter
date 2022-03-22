SHELL := /bin/sh
OUT := $(shell pwd)/_out
BUILDARCH := $(shell uname -m)
GCC := $(OUT)/$(BUILDARCH)-linux-musl-cross/bin/$(BUILDARCH)-linux-musl-gcc
LD := $(OUT)/$(BUILDARCH)-linux-musl-cross/bin/$(BUILDARCH)-linux-musl-ld

clean-compile: musl deps compile

compile:
	CGO_ENABLED=1 \
	CC_FOR_TARGET=$(GCC) \
	CC=$(GCC) \
	go build \
		-ldflags '-linkmode external -extldflags -static' \
		-a -o _out/unbound_explorer .

deps:
	go mod tidy -v
	go mod download

musl: clean
	(cd $(OUT); curl -LOk https://musl.cc/$(BUILDARCH)-linux-musl-cross.tgz)
	tar zxf $(OUT)/$(BUILDARCH)-linux-musl-cross.tgz -C $(OUT)

clean:
	rm -Rf $(OUT) $(BINARY_NAME)
	mkdir -p $(OUT)
	touch $(OUT)/.keep
