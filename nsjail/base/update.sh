#!/bin/bash -eu

VERSION=4
IMAGE=rutsky/jail-base:${VERSION}

docker build -t ${IMAGE} . && docker push ${IMAGE}
