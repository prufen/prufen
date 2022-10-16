#!/bin/bash -eu

VERSION=5
IMAGE=rutsky/jail-base:${VERSION}

docker build -t ${IMAGE} . && docker push ${IMAGE}
