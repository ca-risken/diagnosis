FROM golang:1.18.2 AS builder
WORKDIR /go/src/github.com/ca-risken/diagnosis/cmd/applicationscan/
ADD *.go go.* ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/

WORKDIR /go/src/github.com/gassara-kys
ARG ENV_INJECTOR_VERSION=v0.0.6
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go mod download \
  && CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/

FROM public.ecr.aws/w4a9j2y8/diagnosis/applicationscan-base:v2.10.0
COPY --from=builder /go/bin/env-injector /usr/local/bin/
COPY --from=builder /go/bin/applicationscan /usr/local/bin/
ENV DEBUG= \
  PROFILE_EXPORTER= \
  PROFILE_TYPES= \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
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
WORKDIR /usr/local/
ENTRYPOINT ["bin/env-injector"]
CMD ["bin/applicationscan"]
