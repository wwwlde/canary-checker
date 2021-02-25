#!/bin/bash

set -ex

export PLATFORM_CLI_VERSION=v0.28.0
export PLATFORM_CLI="./karina -c test/config.yaml"
export KUBECONFIG=~/.kube/config
export DOCKER_API_VERSION=1.39


if which karina 2>&1 > /dev/null; then
  PLATFORM_CLI="karina -c test/config.yaml"
else
  if [[ "$OSTYPE" == "linux-gnu" ]]; then
    wget -q https://github.com/flanksource/karina/releases/download/$PLATFORM_CLI_VERSION/karina
    chmod +x karina
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    wget -q https://github.com/flanksource/karina/releases/download/$PLATFORM_CLI_VERSION/karina_osx
    cp karina_osx karina
    chmod +x karina
  else
    echo "OS $OSTYPE not supported"
    exit 1
  fi
fi

mkdir -p .bin

$PLATFORM_CLI ca generate --name root-ca --cert-path .certs/root-ca.crt --private-key-path .certs/root-ca.key --password foobar  --expiry 1
$PLATFORM_CLI ca generate --name ingress-ca --cert-path .certs/ingress-ca.crt --private-key-path .certs/ingress-ca.key --password foobar  --expiry 1
$PLATFORM_CLI provision kind-cluster


$PLATFORM_CLI deploy phases --crds --base --stubs --calico --minio
#$PLATFORM_CLI test stubs --wait=480 -v 5
$PLATFORM_CLI apply test/setup.yml

export DOCKER_USERNAME=test
export DOCKER_PASSWORD=password

wget -q https://github.com/atkrad/wait4x/releases/download/v0.3.0/wait4x-linux-amd64  -O ./wait4x
chmod +x ./wait4x

./wait4x tcp 127.0.0.1:30636 -t 120s -i 5s || true
./wait4x tcp 127.0.0.1:30389 || true
./wait4x tcp 127.0.0.1:32432 || true

make vue-dist
cd test
go test ./... -v -c
# ICMP requires privelages so we run the tests with sudo
sudo DOCKER_API_VERSION=1.39 --preserve-env=KUBECONFIG ./test.test  -test.v
