FROM golang:1.16 AS builder

WORKDIR /go/src/github.com/CyberAgent/mimosa-diagnosis/cmd/application_scan/

ADD *.go go.* ./
ENV CGO_ENABLED=0 GOPRIVATE="github.com/CyberAgent/*,github.com/ca-risken/*"
ARG ENV_INJECTOR_VERSION=v0.0.6
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go build && cp application_scan /go/bin/

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM owasp/zap2docker-bare:latest
COPY --from=builder /go/bin/env-injector /usr/local/bin/
COPY --from=builder /go/bin/application_scan /usr/local/application_scan/bin/
ENV DEBUG= \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  AWS_XRAY_SDK_DISABLED=TRUE \
  AWS_XRAY_CONTEXT_MISSING=LOG_ERROR \
  SQS_ENDPOINT= \
  APPLICATION_SCAN_QUEUE_NAME= \
  APPLICATION_SCAN_QUEUE_URL= \
  MAX_NUMBER_OF_MESSAGE= \
  WAIT_TIME_SECOND=20 \
  SCAN_EXCLUDE_NUMBER= \
  FINDING_SVC_ADDR= \
  ALERT_SVC_ADDR= \
  AWS_SVC_ADDR= \
  ZAP_PORT= \
  TZ=Asia/Tokyo
WORKDIR /usr/local/application_scan
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/application_scan"]
