#!/bin/bash -eu

VERSION=6
IMAGE=rutsky/jail-base:${VERSION}

docker build -t ${IMAGE} . && docker push ${IMAGE}
