.PHONY: all install clean network fmt build doc
all: run

install:
	go get \
		google.golang.org/grpc \
		github.com/golang/protobuf/protoc-gen-go \
		github.com/grpc-ecosystem/go-grpc-middleware

clean:
	rm -f pkg/pb/**/*.pb.go
	rm -f doc/*.md

# @see https://github.com/CyberAgent/mimosa-common/tree/master/local
network:
	@if [ -z "`docker network ls | grep local-shared`" ]; then docker network create local-shared; fi

fmt: proto/*.proto
	clang-format -i proto/*.proto

doc: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--doc_out=markdown,README.md:doc \
		proto/*.proto;

diagnosis-build: fmt
	protoc \
		--proto_path=proto \
		--error_format=gcc \
		--go_out=pkg/pb/diagnosis \
		--go_opt=plugins=grpc,paths=source_relative proto/*.proto \
		proto/*.proto;

go-test: build
	cd pkg/pb/diagnosis && go test ./...
	cd pkg/message      && go test ./...
	cd cmd/diagnosis    && go test ./...
	cd cmd/jira         && go test ./...

go-mod-update:
	cd cmd/diagnosis \
		&& go get -u \
			github.com/CyberAgent/mimosa-diagnosis/...
	cd cmd/jira \
		&& go get -u \
			github.com/CyberAgent/mimosa-core/... \
			github.com/CyberAgent/mimosa-diagnosis/...

go-mod-tidy: build
	cd pkg/message   && go mod tidy
	cd cmd/diagnosis && go mod tidy
	cd cmd/jira      && go mod tidy

run: go-test network
	. env.sh && docker-compose up -d --build

log:
	. env.sh && docker-compose logs -f

stop:
	. env.sh && docker-compose down
