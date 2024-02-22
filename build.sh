set -euxo pipefail

mkdir -p "$(pwd)/functions"
mkdir -p "$(pwd)/public"
GOBIN=$(pwd)/functions go install ./...
chmod +x "$(pwd)"/functions/*
go env