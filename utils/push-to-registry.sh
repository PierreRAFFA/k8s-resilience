#!/bin/bash
VERSION=$1
docker build -t pierreraffa/ms-api-gateway:$VERSION ../ms-api-gateway
docker push pierreraffa/ms-api-gateway:$VERSION

docker build -t pierreraffa/ms-payments:$VERSION ../ms-payments
docker push pierreraffa/ms-payments:$VERSION

docker build -t pierreraffa/ms-users:$VERSION ../ms-users
docker push pierreraffa/ms-users:$VERSION