# Copyright (c) Unikraft GmbH
# SPDX-License-Identifier: MPL-2.0

TEST?=$$(go list ./... |grep -v 'vendor')

GIT_VERSION?=$(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)

PROVIDER_HOSTNAME=unikraft.com
PROVIDER_NAMESPACE=unikraft
PROVIDER_TYPE=ukc
PROVIDER_TARGET=$(shell go env GOOS)_$(shell go env GOARCH)
PROVIDER_PATH=~/.terraform.d/plugins/$(PROVIDER_HOSTNAME)/$(PROVIDER_NAMESPACE)/$(PROVIDER_TYPE)/$(VERSION)/$(PROVIDER_TARGET)

default: build

build:
	@echo $(GIT_VERSION)
	@mkdir -p $(PROVIDER_PATH)
	go build \
		-tags release \
		-ldflags '-X $(PACKAGE)/internal/config.Version=$(GIT_VERSION)' \
		-o $(PROVIDER_PATH)/terraform-provider-$(PROVIDER_NAMESPACE)_v$(VERSION)	

test:
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=300s -parallel=4 -count=1