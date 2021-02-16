#!/bin/bash -e

cd "$(dirname "$0")"

# load env
. ../env.sh

# setting remote repository
TAG="local-test-$(date '+%Y%m%d')"
IMAGE_DIAGNOSIS="diagnosis/diagnosis"
IMAGE_JIRA="diagnosis/jira"
IMAGE_WPSCAN="diagnosis/wpscan"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
REGISTORY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

# build & push
aws ecr get-login-password --region ${AWS_REGION} \
  | docker login \
    --username AWS \
    --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_DIAGNOSIS}:${TAG} ../cmd/diagnosis/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_JIRA}:${TAG} ../cmd/jira/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_WPSCAN}:${TAG} ../cmd/wpscan/

docker tag ${IMAGE_DIAGNOSIS}:${TAG}      ${REGISTORY}/${IMAGE_DIAGNOSIS}:${TAG}
docker tag ${IMAGE_JIRA}:${TAG}       ${REGISTORY}/${IMAGE_JIRA}:${TAG}
docker tag ${IMAGE_WPSCAN}:${TAG} ${REGISTORY}/${IMAGE_WPSCAN}:${TAG}

docker push ${REGISTORY}/${IMAGE_DIAGNOSIS}:${TAG}
docker push ${REGISTORY}/${IMAGE_JIRA}:${TAG}
docker push ${REGISTORY}/${IMAGE_WPSCAN}:${TAG}
