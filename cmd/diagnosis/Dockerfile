FROM golang:1.15 as builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-diagnosis/diagnosis
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*,github.com/ca-risken/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

RUN mkdir /go/src/diagnosis
WORKDIR /go/src/diagnosis
ADD . .
#ADD ./proto diagnosis/proto

FROM alpine
RUN apk add --no-cache ca-certificates tzdata \
  && cd /tmp/ && wget -q https://github.com/gassara-kys/env-injector/releases/download/v0.0.6/env-injector_0.0.6_Linux_x86_64.tar.gz \
  && tar xvfz env-injector_0.0.6_Linux_x86_64.tar.gz && mv env-injector /usr/local/bin/ && rm -rf /tmp/*
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