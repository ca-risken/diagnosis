FROM golang:1.18.2 AS builder

WORKDIR /go/src/github.com/ca-risken/diagnosis/src/wpscan/

ADD *.go go.* ./
ENV CGO_ENABLED=0
ARG ENV_INJECTOR_VERSION=v0.0.6
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go mod download
RUN GO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /go/bin/

WORKDIR /go/src/github.com/gassara-kys
RUN git clone https://github.com/gassara-kys/env-injector.git -b ${ENV_INJECTOR_VERSION} \
  && cd env-injector \
  && go build \
  && cp env-injector /go/bin/

FROM public.ecr.aws/risken/risken-diagnosis/wpscan-base:v3.8.18
USER root 
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /go/bin/env-injector /usr/local/bin/
COPY --from=builder /go/bin/wpscan /usr/local/wpscan/bin/
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

WORKDIR /usr/local/wpscan
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/wpscan"]
