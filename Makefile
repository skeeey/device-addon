
# Image URL to use all building/pushing image targets
IMG ?= quay.io/skeeey/device-addon:latest

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## generate code and crds
.PHONY: crds-gen
crds-gen: controller-gen
	hack/crds_gen.sh

.PHONY: code-gen
code-gen: code-generator
	hack/code_gen.sh

.PHONY: update
update: code-gen crds-gen

## build binary
.PHONY: build
build:
	go build -o bin/device-addon cmd/main.go

.PHONY: build-thermometer
build-thermometer:
	go build -o bin/thermometer contrib/demo/device/thermometer.go

## build and push image
.PHONY: image
image:
	docker build -t ${IMG} .

.PHONY: image-push
image-push: image
	docker push ${IMG}

## run demo
.PHONY: run-thermometer-a
run-thermometer-a:
	contrib/demo/run-device.sh room-a

.PHONY: run-thermometer-b
run-thermometer-b:
	contrib/demo/run-device.sh room-b

.PHONY: demo
demo:
	contrib/demo/demo.sh

## Tool versions
CONTROLLER_TOOLS_VERSION ?= v0.9.2
CODE_GENERATOR_VERSION ?= v0.26.1

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: controller-gen
controller-gen: $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: code-generator
code-generator: $(LOCALBIN)
	test -s $(LOCALBIN)/client-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/client-gen@$(CODE_GENERATOR_VERSION)
	test -s $(LOCALBIN)/informer-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/informer-gen@$(CODE_GENERATOR_VERSION)
	test -s $(LOCALBIN)/lister-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/lister-gen@$(CODE_GENERATOR_VERSION)
	test -s $(LOCALBIN)/deepcopy-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GENERATOR_VERSION)
	test -s $(LOCALBIN)/openapi-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/openapi-gen@$(CODE_GENERATOR_VERSION)
