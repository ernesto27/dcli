#!/bin/bash

echo "Creating tag and push: $1"

git tag -a $1 -m "new release"
git push origin $1

echo "Create and upload release"
goreleaser --clean