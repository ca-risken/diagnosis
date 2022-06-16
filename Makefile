TARGETS = wpscan portscan applicationscan
BUILD_TARGETS = $(TARGETS:=.build)
BUILD_CI_TARGETS = $(TARGETS:=.build-ci)
IMAGE_PUSH_TARGETS = $(TARGETS:=.push-image)
IMAGE_PULL_TARGETS = $(TARGETS:=.pull-image)
IMAGE_TAG_TARGETS = $(TARGETS:=.tag-image)
MANIFEST_CREATE_TARGETS = $(TARGETS:=.create-manifest)
MANIFEST_PUSH_TARGETS = $(TARGETS:=.push-manifest)
TEST_TARGETS = $(TARGETS:=.go-test)
LINT_TARGETS = $(TARGETS:=.lint)
BUILD_OPT=""
IMAGE_TAG=latest
MANIFEST_TAG=latest
IMAGE_PREFIX=diagnosis
IMAGE_REGISTRY=local

.PHONY: all
all: build

.PHONY: ecr-login
ecr-login:
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

.PHONY: build
build: $(BUILD_TARGETS)
%.build: %.go-test
	TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh

.PHONY: build-ci
build-ci: $(BUILD_CI_TARGETS)
%.build-ci: FAKE
	TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh
	docker tag $(IMAGE_PREFIX)/$(*):$(IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

.PHONY: push-image
push-image: $(IMAGE_PUSH_TARGETS)
%.push-image: FAKE
	docker push $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

PHONY: pull-image $(IMAGE_PULL_TARGETS)
pull-image: $(IMAGE_PULL_TARGETS)
%.pull-image:
	docker pull $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

PHONY: tag-image $(IMAGE_TAG_TARGETS)
tag-image: $(IMAGE_TAG_TARGETS)
%.tag-image:
	docker tag $(SOURCE_IMAGE_PREFIX)/$(*):$(SOURCE_IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

.PHONY: create-manifest
create-manifest: $(MANIFEST_CREATE_TARGETS)
%.create-manifest: FAKE
	docker manifest create $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_amd64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_arm64
	docker manifest annotate --arch amd64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_amd64
	docker manifest annotate --arch arm64 $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG_BASE)_linux_arm64

.PHONY: push-manifest
push-manifest: $(MANIFEST_PUSH_TARGETS)
%.push-manifest: FAKE
	docker manifest push $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG)
	docker manifest inspect $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(MANIFEST_TAG)

.PHONY: go-test pkg-test
go-test: $(TEST_TARGETS) pkg-test
%.go-test: FAKE
	cd cmd/$(*) && GO111MODULE=on go test ./...
pkg-test:
	cd pkg/common && GO111MODULE=on go test ./...

.PHONY: go-mod-update
go-mod-update:
	cd cmd/wpscan \
		&& go get -u \
			github.com/ca-risken/core/... \
			github.com/ca-risken/datasource-api/... \
			github.com/ca-risken/diagnosis/...
	cd cmd/portscan \
		&& go get -u \
			github.com/ca-risken/core/... \
			github.com/ca-risken/datasource-api/... \
			github.com/ca-risken/diagnosis/...
	cd cmd/applicationscan \
		&& go get -u \
			github.com/ca-risken/core/... \
			github.com/ca-risken/datasource-api/... \
			github.com/ca-risken/diagnosis/...

.PHONY: go-mod-tidy
go-mod-tidy:
	cd cmd/wpscan           && go mod tidy
	cd cmd/portscan         && go mod tidy
	cd cmd/applicationscan  && go mod tidy

.PHONY: lint pkg-lint
lint: $(LINT_TARGETS) pkg-lint
%.lint: FAKE
	sh hack/golinter.sh cmd/$(*)
pkg-lint:
	sh hack/golinter.sh pkg/common

FAKE:
