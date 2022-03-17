FROM golang:1.17 AS builder

WORKDIR /go/src/github.com/ca-risken/diagnosis/cmd/applicationscan/

ADD *.go go.* ./
ENV CGO_ENABLED=0
ARG ENV_INJECTOR_VERSION=v0.0.6
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go build && cp applicationscan /go/bin/

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM public.ecr.aws/w4a9j2y8/diagnosis/applicationscan-base:v2.10.0
COPY --from=builder /go/bin/env-injector /usr/local/bin/
COPY --from=builder /go/bin/applicationscan /usr/local/applicationscan/bin/
ENV DEBUG= \
  PROFILE_EXPORTER= \
  PROFILE_TYPES= \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  AWS_XRAY_SDK_DISABLED=TRUE \
  AWS_XRAY_CONTEXT_MISSING=LOG_ERROR \
  SQS_ENDPOINT= \
  DIAGNOSIS_APPLICATION_SCAN_QUEUE_NAME= \
  DIAGNOSIS_APPLICATION_SCAN_QUEUE_URL= \
  MAX_NUMBER_OF_MESSAGE= \
  WAIT_TIME_SECOND=20 \
  SCAN_EXCLUDE_NUMBER= \
  FINDING_SVC_ADDR= \
  ALERT_SVC_ADDR= \
  AWS_SVC_ADDR= \
  ZAP_PORT= \
  TZ=Asia/Tokyo
WORKDIR /usr/local/applicationscan
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/applicationscan"]
