FROM golang:1.16 as builder

WORKDIR /go/src/github.com/ca-risken/diagnosis/diagnosis

ADD *.go go.* ./
ENV CGO_ENABLED=0 GOPRIVATE="github.com/ca-risken/*"
ARG ENV_INJECTOR_VERSION=v0.0.6
RUN echo "machine github.com" > ~/.netrc
RUN --mount=type=secret,id=GITHUB_USER echo "login $(cat /run/secrets/GITHUB_USER)" >> ~/.netrc
RUN --mount=type=secret,id=GITHUB_TOKEN echo "password $(cat /run/secrets/GITHUB_TOKEN)" >> ~/.netrc
RUN go build && cp diagnosis /go/bin/

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM alpine
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /go/bin/env-injector /usr/local/bin/
COPY --from=builder /go/bin/diagnosis /usr/local/diagnosis/bin/
ENV PORT=19001 \
  DB_MASTER_HOST= \
  DB_MASTER_USER= \
  DB_MASTER_PASSWORD= \
  DB_SLAVE_HOST= \
  DB_SLAVE_USER= \
  DB_SLAVE_PASSWORD= \
  DB_SCHEMA=mimosa \
  DB_PORT=3306 \
  DB_LOG_MODE=false \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  AWS_XRAY_SDK_DISABLED=TRUE \
  AWS_XRAY_CONTEXT_MISSING=LOG_ERROR \
  SQS_ENDPOINT= \
  DIAGNOSIS_JIRA_QUEUE_URL= \
  DIAGNOSIS_JIRA_QUEUE_NAME= \
  DIAGNOSIS_WPSCAN_QUEUE_URL= \
  DIAGNOSIS_WPSCAN_QUEUE_NAME= \
  TZ=Asia/Tokyo
WORKDIR /usr/local/diagnosis
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/diagnosis"]