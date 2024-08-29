#!/bin/sh

set -euxo pipefail

go test -v -cover ./internal/validators/
TF_ACC=1 go test -v -cover ./internal/provider/