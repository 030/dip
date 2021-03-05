#!/bin/bash -e
GITHUB_TAG="${GITHUB_TAG:-local}"
SHA512_CMD="${SHA512_CMD:-sha512sum}"
export DIP_DELIVERABLE="${DIP_DELIVERABLE:-dip}"

echo "GITHUB_TAG: '$GITHUB_TAG' DIP_DELIVERABLE: '${DIP_DELIVERABLE}'"
cd cmd/dip
go build -ldflags "-X main.Version=${GITHUB_TAG}" -o "${DIP_DELIVERABLE}"
$SHA512_CMD "${DIP_DELIVERABLE}" > "${DIP_DELIVERABLE}.sha512.txt"
chmod +x "${DIP_DELIVERABLE}"
cd ../..
