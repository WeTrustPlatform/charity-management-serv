#! /bin/sh

# This script builds and publish docker images
# Usage: ./docker.sh [tagname]
# Example: ./docker.sh v1.1.0
# will build and publish 2 tags
# sihoang/charity-management-serv:latest
# and sihoang/charity-management-serv:v1.1.0

DOCKER_REPO="sihoang/charity-management-serv"
DOCKER_IMAGE="$DOCKER_REPO:latest"

GIT_COMMIT=`git rev-parse HEAD`

docker build . \
  -t $DOCKER_IMAGE \
  --build-arg GIT_COMMIT=$GIT_COMMIT

echo ">> Pushing image $DOCKER_IMAGE to registry"
docker push $DOCKER_IMAGE

if [ -z "$1" ]; then
  echo ">> No tag specified"
  exit 0
else
  echo ">> Tag $1"
  DOCKER_TAG="$DOCKER_REPO:$1"
  docker tag $DOCKER_IMAGE $DOCKER_TAG
  echo ">> Pushing tag $DOCKER_TAG to registry"
  docker push $DOCKER_TAG
fi
