#!/bin/bash -e

cd "$(dirname "$0")"

# load env
. ../env.sh

# setting remote repository
TAG="local-test-$(date '+%Y%m%d')"
IMAGE_DIAGNOSIS="diagnosis/diagnosis"
IMAGE_WPSCAN="diagnosis/wpscan"
IMAGE_PORTSCAN="diagnosis/portscan"
IMAGE_APPLICATIONSCAN="diagnosis/applicationscan"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
REGISTORY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

# build & push
aws ecr get-login-password --region ${AWS_REGION} \
  | docker login \
    --username AWS \
    --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_DIAGNOSIS}:${TAG} ../cmd/diagnosis/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_WPSCAN}:${TAG} ../cmd/wpscan/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_PORTSCAN}:${TAG} ../cmd/portscan/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_APPLICATIONSCAN}:${TAG} ../cmd/applicationscan/

docker tag ${IMAGE_DIAGNOSIS}:${TAG}      ${REGISTORY}/${IMAGE_DIAGNOSIS}:${TAG}
docker tag ${IMAGE_WPSCAN}:${TAG} ${REGISTORY}/${IMAGE_WPSCAN}:${TAG}
docker tag ${IMAGE_WPSCAN}:${TAG} ${REGISTORY}/${IMAGE_PORTSCAN}:${TAG}
docker tag ${IMAGE_WPSCAN}:${TAG} ${REGISTORY}/${IMAGE_APPLICATIONSCAN}:${TAG}

docker push ${REGISTORY}/${IMAGE_DIAGNOSIS}:${TAG}
docker push ${REGISTORY}/${IMAGE_WPSCAN}:${TAG}
docker push ${REGISTORY}/${IMAGE_PORTSCAN}:${TAG}
docker push ${REGISTORY}/${IMAGE_APPLICATIONSCAN}:${TAG}
