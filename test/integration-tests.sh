#!/bin/bash -e

source ./scripts/build.sh
cd cmd/dip

readonly DIP_ERROR="Cannot find the latest tag. Check whether the tags are semantic"

echo -e "\nAdoptopenjdk"
./"${DIP_DELIVERABLE}" -image=adoptopenjdk -latest="14.*-jre-hotspot-bionic" -official

echo -e "\nAlpine"
./"${DIP_DELIVERABLE}" -image alpine -latest "(\d+\.){2}\d" --official | grep "[0-9]"

echo -e "\nGolang"
./"${DIP_DELIVERABLE}" -image=golang -latest "^1\..*-alpine[0-9]+\.[0-9]+$" -official | grep "[0-9]"

echo -e "\nMinio"
./"${DIP_DELIVERABLE}" -image minio/minio -latest "RELEASE\.2019.*" | grep "${DIP_ERROR}" | grep "[0-9]" || true

echo -e "\nNginx"
./"${DIP_DELIVERABLE}" -image=nginx -latest=1\..* -official | grep "[0-9]"

echo -e "\nSonatype Nexus3"
./"${DIP_DELIVERABLE}" -image=sonatype/nexus3 -latest=3\..* | grep "[0-9]"

echo -e "\nSonarqube"
./"${DIP_DELIVERABLE}" -image sonarqube -latest ".*-community$" -official | grep "[0-9]"

echo -e "\nTraefik"
./"${DIP_DELIVERABLE}" --image=traefik --latest="^v(\d+\.){1,2}\d+$" -official | grep "[0-9]"

echo -e "\nUbuntu"
./"${DIP_DELIVERABLE}" -image ubuntu -latest "^xenial.*" -official | grep "${DIP_ERROR}" || true
