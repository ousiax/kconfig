#!/bin/sh
set -e

gitVersion=$(git describe)
gitCommit=$(git rev-parse HEAD)
buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
ldflags="\
-X 'k8s.io/client-go/pkg/version.gitVersion=$gitVersion' \
-X 'k8s.io/client-go/pkg/version.gitCommit=$gitCommit' \
-X 'k8s.io/client-go/pkg/version.buildDate=$buildDate' \
"

go build -ldflags="$ldflags"
