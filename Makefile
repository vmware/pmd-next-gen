HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
VERSION := ${HASH} (${COMMIT_DATE})

BUILDDIR ?= .
SRCDIR ?= .

.PHONY: help
help:
	@echo "make [TARGETS...]"
	@echo
	@echo "This is the maintenance makefile of photon-mgmtd. The following"
	@echo "targets are available:"
	@echo
	@echo "    help:               Print this usage information."
	@echo "    build:              Builds project"
	@echo "    install:            Installs binary, configuration and unit files"
	@echo "    clean:              Cleans the build"

$(BUILDDIR)/:
	mkdir -p "$@"

$(BUILDDIR)/%/:
	mkdir -p "$@"

.PHONY: build
build:
	- mkdir -p bin
	go build -buildmode=pie -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'" -o bin/photon-mgmtd ./cmd/photon-mgmt
	go build -buildmode=pie -ldflags="-X 'main.buildVersion=${VERSION}' -X 'main.buildDate=${BUILD_DATE}'" -o bin/pmctl ./cmd/pmctl

.PHONY: install
install:
	install bin/photon-mgmtd /usr/bin/
	install bin/pmctl /usr/bin/

	install -vdm 755 /etc/photon-mgmt
	install -m 755 conf/photon-mgmt.toml /etc/photon-mgmt
	install -m 755 conf/photon-mgmt-auth.conf /etc/photon-mgmt

	install -m 0644 units/photon-mgmtd.service /lib/systemd/system/
	systemctl daemon-reload

PHONY: clean
clean:
	go clean
	rm -rf bin
