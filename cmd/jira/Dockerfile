FROM golang AS builder
WORKDIR /go/src/github.com/CyberAgent/mimosa-diagnosis/src/jira/
ADD *.go go.* ./
ENV CGO_ENABLED=0 GO111MODULE=on GOPRIVATE="github.com/CyberAgent/*"
ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN echo "machine github.com" > ~/.netrc
RUN echo "login $GITHUB_USER" >> ~/.netrc
RUN echo "password $GITHUB_TOKEN" >> ~/.netrc
RUN go install -v

FROM alpine
RUN apk add --no-cache ca-certificates tzdata \
  && cd /tmp/ && wget -q https://github.com/okzk/env-injector/releases/download/v0.0.5/env-injector_0.0.5_linux_amd64.tar.gz \
  && tar xvfz env-injector_0.0.5_linux_amd64.tar.gz && mv env-injector /usr/local/bin/ && rm /tmp/*
COPY --from=builder /go/bin/jira /usr/local/jira/bin/
ENV AWS_REGION= \
    AWS_ACCESS_KEY_ID= \
    AWS_SECRET_ACCESS_KEY= \
    AWS_SESSION_TOKEN= \
    SQS_ENDPOINT= \
    JIRA_QUEUE_NAME= \
    JIRA_QUEUE_URL= \
    MAX_NUMBER_OF_MESSAGE=10 \
    WAIT_TIME_SECOND=20 \
    FINDING_SVC_ADDR= \
    DIAGNOSIS_SVC_ADDR= \
    JIRA_URL= \
    JIRA_USER_ID= \
    JIRA_USER_PASSWORD= \
    TZ=Asia/Tokyo
WORKDIR /usr/local/jira
ENTRYPOINT ["/usr/local/bin/env-injector"]
#CMD ["bin/jira"]
CMD ["sh"]
