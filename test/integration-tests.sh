#!/bin/bash -e

source ./scripts/build.sh
cd cmd/dip

echo -e "\nAdoptopenjdk"
./"${DIP_DELIVERABLE}" image --name=adoptopenjdk --regex="14.*-jre-hotspot-bionic" | grep "[0-9]"

echo -e "\nAdoptopenjdk16"
./"${DIP_DELIVERABLE}" image --name=adoptopenjdk --regex="16.*" | grep "[0-9]"

echo -e "\nAlpine"
./"${DIP_DELIVERABLE}" image --name=alpine --regex "(\d+\.){2}\d" | grep "[0-9]"

echo -e "\nGolang"
./"${DIP_DELIVERABLE}" image --name=golang --regex "^1\..*-alpine[0-9]+\.[0-9]+$" | grep "[0-9]"

echo -e "\nMinio"
./"${DIP_DELIVERABLE}" image --name=minio/minio --regex "RELEASE\.2021.*" | grep "[0-9]"

echo -e "\nN3DR"
n3drTag="^6.0.10$"
./"${DIP_DELIVERABLE}" image --name=utrecht/n3dr --regex=$n3drTag | grep $n3drTag

echo -e "\nNginx"
./"${DIP_DELIVERABLE}" image --name=nginx --regex=1\..* | grep "[0-9]"

echo -e "\nSonatype Nexus3"
./"${DIP_DELIVERABLE}" image --name=sonatype/nexus3 --regex=3\..* | grep "[0-9]"

echo -e "\nSonarqube"
./"${DIP_DELIVERABLE}" image --name=sonarqube --regex ".*-community$" | grep "[0-9]"

echo -e "\nTraefik"
./"${DIP_DELIVERABLE}" image --name=traefik --regex="^v(\d+\.){1,2}\d+$" | grep "[0-9]"

echo -e "\nUbuntu"
./"${DIP_DELIVERABLE}" image --name=ubuntu --regex "^impish.*" | grep "[0-9]"

echo -e "\nUbuntu Hirsute"
ubuntuHirsuteTag="^hirsute-20210522$"
./"${DIP_DELIVERABLE}" image --name=ubuntu --regex $ubuntuHirsuteTag | grep $ubuntuHirsuteTag
