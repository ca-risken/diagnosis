TARGETS = diagnosis jira wpscan portscan
BUILD_TARGETS = $(TARGETS:=.build)
BUILD_CI_TARGETS = $(TARGETS:=.build-ci)
IMAGE_PUSH_TARGETS = $(TARGETS:=.push-image)
MANIFEST_CREATE_TARGETS = $(TARGETS:=.create-manifest)
MANIFEST_PUSH_TARGETS = $(TARGETS:=.push-manifest)
TEST_TARGETS = $(TARGETS:=.go-test)
BUILD_OPT=""
IMAGE_TAG=latest
MANIFEST_TAG=latest
IMAGE_PREFIX=diagnosis
IMAGE_REGISTRY=local

.PHONY: all
all: build

.PHONY: install
install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/go-grpc-middleware

.PHONY: clean
clean:
	rm -f proto//**/*.pb.go
	rm -f doc/*.md

.PHONY: fmt
fmt: proto/**/*.proto
	clang-format -i proto/**/*.proto

.PHONY: doc
doc: fmt
	protoc \
		--proto_path=proto \
		--proto_path=${GOPATH}/src \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/**/*.proto;

.PHONY: proto
proto: fmt
	protoc \
		--proto_path=proto \
		--proto_path=${GOPATH}/src \
		--error_format=gcc \
		--go_out=plugins=grpc,paths=source_relative:proto proto/**/*.proto \
		proto/**/*.proto;

.PHONY: build
build: $(BUILD_TARGETS)
%.build: %.go-test
	. env.sh && TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh

.PHONY: build-ci
build-ci: $(BUILD_CI_TARGETS)
%.build-ci: FAKE
	TARGET=$(*) IMAGE_TAG=$(IMAGE_TAG) IMAGE_PREFIX=$(IMAGE_PREFIX) BUILD_OPT="$(BUILD_OPT)" . hack/docker-build.sh
	docker tag $(IMAGE_PREFIX)/$(*):$(IMAGE_TAG) $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

.PHONY: push-image
push-image: $(IMAGE_PUSH_TARGETS)
%.push-image: FAKE
	docker push $(IMAGE_REGISTRY)/$(IMAGE_PREFIX)/$(*):$(IMAGE_TAG)

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

.PHONY: go-test proto-test pkg-test
go-test: $(TEST_TARGETS) proto-test pkg-test
%.go-test: FAKE
	cd cmd/$(*) && go test ./...
proto-test:
	cd proto/diagnosis  && go test ./...
pkg-test:
	cd pkg/message      && go test ./...

.PHONY: go-mod-update
go-mod-update:
	cd cmd/diagnosis \
		&& go get -u \
			github.com/ca-risken/diagnosis/...
	cd cmd/jira \
		&& go get -u \
			github.com/ca-risken/core/... \
			github.com/ca-risken/diagnosis/...
	cd cmd/wpscan \
		&& go get -u \
			github.com/ca-risken/core/... \
			github.com/ca-risken/diagnosis/...
	cd cmd/portscan \
		&& go get -u \
			github.com/ca-risken/core/... \
			github.com/ca-risken/diagnosis/...

.PHONY: go-mod-tidy
go-mod-tidy: proto
	cd pkg/message   && go mod tidy
	cd cmd/diagnosis && go mod tidy
	cd cmd/jira      && go mod tidy
	cd cmd/wpscan    && go mod tidy
	cd cmd/portscan  && go mod tidy

FAKE:
