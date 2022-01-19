FROM golang:1.17 AS builder

WORKDIR /go/src/github.com/ca-risken/diagnosis/src/jira/

ADD *.go go.* ./
ENV CGO_ENABLED=0 GOPRIVATE="github.com/ca-risken/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go build && cp jira /go/bin/

FROM public.ecr.aws/risken/base/risken-base:v0.0.1
COPY --from=builder /go/bin/jira /usr/local/jira/bin/
ENV AWS_REGION= \
  AWS_ACCESS_KEY_ID= \
  AWS_SECRET_ACCESS_KEY= \
  AWS_SESSION_TOKEN= \
  AWS_XRAY_SDK_DISABLED=TRUE \
  AWS_XRAY_CONTEXT_MISSING=LOG_ERROR \
  SQS_ENDPOINT= \
  DIAGNOSIS_JIRA_QUEUE_URL= \
  DIAGNOSIS_JIRA_QUEUE_NAME= \
  MAX_NUMBER_OF_MESSAGE=10 \
  WAIT_TIME_SECOND=20 \
  FINDING_SVC_ADDR= \
  DIAGNOSIS_SVC_ADDR= \
  DIAGNOSIS_JIRA_URL= \
  DIAGNOSIS_JIRA_USER_ID= \
  DIAGNOSIS_JIRA_USER_PASSWORD= \
  PORT=19002 \
  TZ=Asia/Tokyo

WORKDIR /usr/local/jira
ENTRYPOINT ["/usr/local/bin/env-injector"]
CMD ["bin/jira"]
