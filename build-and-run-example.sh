#!/bin/bash

set -euxo pipefail

./build.sh

cd examples/provider-install-verification && terraform apply -auto-approve