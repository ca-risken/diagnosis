FROM golang:1.15 AS builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-diagnosis/src/wpscan/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

FROM wpscanteam/wpscan:latest
USER root 
RUN apk add --no-cache ca-certificates tzdata \
  && cd /tmp/ && wget -q https://github.com/gassara-kys/env-injector/releases/download/v0.0.6/env-injector_0.0.6_Linux_x86_64.tar.gz \
  && tar xvfz env-injector_0.0.6_Linux_x86_64.tar.gz && mv env-injector /usr/local/bin/ && rm -rf /tmp/* \
  && mkdir /results
COPY --from=builder /go/bin/wpscan /usr/local/wpscan/bin/
ENV PORT=19003 \
    AWS_REGION= \
    AWS_ACCESS_KEY_ID= \
    AWS_SECRET_ACCESS_KEY= \
    AWS_SESSION_TOKEN= \
    AWS_XRAY_SDK_DISABLED=TRUE \
    AWS_XRAY_CONTEXT_MISSING=LOG_ERROR \
    SQS_ENDPOINT= \
    DIAGNOSIS_WPSCAN_QUEUE_NAME= \
    DIAGNOSIS_WPSCAN_QUEUE_URL= \
    MAX_NUMBER_OF_MESSAGE=5 \
    WAIT_TIME_SECOND=20 \
    FINDING_SVC_ADDR= \
    ALERT_SVC_ADDR= \
    DIAGNOSIS_SVC_ADDR= \
    RESULT_PATH=/results \
    WPSCAN_VULNDB_APIKEY= \
    TZ=Asia/Tokyo

WORKDIR /usr/local/wpscan
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/wpscan"]
