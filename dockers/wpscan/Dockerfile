FROM golang:1.18.2 AS builder
WORKDIR /go/src/github.com/ca-risken/diagnosis/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/wpscan cmd/wpscan/main.go

WORKDIR /go/src/github.com/gassara-kys
ARG ENV_INJECTOR_VERSION=v0.0.6
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go mod download \
  && CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/

FROM public.ecr.aws/risken/risken-diagnosis/wpscan-base:v3.8.18
USER root 
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /go/bin/env-injector /usr/local/bin/
COPY --from=builder /go/bin/wpscan /usr/local/bin/
ENV PORT=19003 \
  PROFILE_EXPORTER= \
  PROFILE_TYPES= \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  SQS_ENDPOINT= \
  DIAGNOSIS_WPSCAN_QUEUE_NAME= \
  DIAGNOSIS_WPSCAN_QUEUE_URL= \
  MAX_NUMBER_OF_MESSAGE= \
  WAIT_TIME_SECOND=20 \
  FINDING_SVC_ADDR= \
  ALERT_SVC_ADDR= \
  DIAGNOSIS_SVC_ADDR= \
  RESULT_PATH= \
  WPSCAN_VULNDB_APIKEY= \
  TZ=Asia/Tokyo

WORKDIR /usr/local/
ENTRYPOINT ["bin/env-injector"]
CMD ["bin/wpscan"]
