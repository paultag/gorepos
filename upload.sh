#!/bin/bash
set -e

go run main.go
aws s3 sync go s3://pault.ag/go
