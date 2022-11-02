FROM golang:1.18.2 AS builder
WORKDIR /go/src/github.com/ca-risken/diagnosis/cmd/portscan/
ADD *.go go.* ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/

FROM public.ecr.aws/risken/base/risken-base:v0.0.1
RUN apk add --no-cache nmap nmap-scripts
COPY --from=builder /go/bin/portscan /usr/local/bin/
ENV DEBUG= \
  PROFILE_EXPORTER= \
  PROFILE_TYPES= \
  AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  SQS_ENDPOINT= \
  PORTSCAN_QUEUE_NAME= \
  PORTSCAN_QUEUE_URL= \
  MAX_NUMBER_OF_MESSAGE= \
  WAIT_TIME_SECOND=20 \
  SCAN_EXCLUDE_NUMBER= \
  FINDING_SVC_ADDR= \
  ALERT_SVC_ADDR= \
  AWS_SVC_ADDR= \
  TZ=Asia/Tokyo
WORKDIR /usr/local/
ENTRYPOINT ["bin/env-injector"]
CMD ["bin/portscan"]
