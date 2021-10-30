# DIP

[![GoDoc Widget](https://godoc.org/github.com/030/dip?status.svg)](https://godoc.org/github.com/030/dip)
[![Go Report Card](https://goreportcard.com/badge/github.com/030/dip)](https://goreportcard.com/report/github.com/030/dip)
[![Build Status](https://travis-ci.org/030/dip.svg?branch=master)](https://travis-ci.org/030/dip)
[![DevOps SE Questions](https://img.shields.io/stackexchange/devops/t/dip.svg)](https://devops.stackexchange.com/questions/tagged/dip)
![Docker Pulls](https://img.shields.io/docker/pulls/utrecht/dip.svg)
![Issues](https://img.shields.io/github/issues-raw/030/dip.svg)
![Pull requests](https://img.shields.io/github/issues-pr-raw/030/dip.svg)
![Total downloads](https://img.shields.io/github/downloads/030/dip/total.svg)
![License](https://img.shields.io/github/license/030/dip.svg)
![Repository Size](https://img.shields.io/github/repo-size/030/dip.svg)
![Contributors](https://img.shields.io/github/contributors/030/dip.svg)
![Commit activity](https://img.shields.io/github/commit-activity/m/030/dip.svg)
![Last commit](https://img.shields.io/github/last-commit/030/dip.svg)
![Release date](https://img.shields.io/github/release-date/030/dip.svg)
![Latest Production Release Version](https://img.shields.io/github/release/030/dip.svg)
[![codecov](https://codecov.io/gh/030/dip/branch/master/graph/badge.svg)](https://codecov.io/gh/030/dip)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com/r/github.com/030/dip)
[![BCH compliance](https://bettercodehub.com/edge/badge/030/dip?branch=master)](https://bettercodehub.com/results/030/dip)
[![Chocolatey](https://img.shields.io/chocolatey/dt/dip)](https://chocolatey.org/packages/dip)
[![Snap](https://snapcraft.io/kdiutd/badge.svg)](https://snapcraft.io/kdiutd)

<a href="https://dip.releasesoftwaremoreoften.com"><img src="https://github.com/030/dip/raw/master/assets/logo/logo.png" width="100"></a>

Docker Image Patrol (DIP) keeps docker images up-to-date.

## Installation

Keep Docker Images Up To Date (KDIUTD)

```bash
sudo snap install kdiutd
```

## Usage

```bash
Usage:
  dip image [flags]

Flags:
      --dockerfile     Check whether the image that resides in the Dockerfile is outdated
  -h, --help           help for image
  -n, --name string    Name of the Docker image to be checked whether it is up to date
  -r, --regex string   Regex for finding the latest image tag
      --sendSlackMsg   Send message to Slack

Global Flags:
      --configCredHome string   config and cred file home directory (default is $HOME/.dip)
      --debug                   debugging mode
```

## k8s

Create a `~/.dip/config.yml` file:

```bash
dip_images:
  docker.io/alpine: 3\.[0-9]+\.[0-9]+
  elastic/elasticsearch: 7\.[0-9]+\.[0-9]+
  fluent/fluentd-kubernetes-daemonset: v.*-debian-elasticsearch7-.*
  grafana/grafana: 7\.[0-9]+\.[0-9]+
  docker.io/kibana: 7\.[0-9]+\.[0-9]+
  kubernetesui/dashboard: v2\.[0-9]+\.[0-9]+
  kubernetesui/metrics-scraper: v1\.[0-9]+\.[0-9]+
  docker.io/mongo: 4\.[0-9]+\.[0-9]+
  docker.io/postgres: 13\.[0-9]+\.[0-9]+
  prom/alertmanager: v0\.[0-9]+\.[0-9]+
  prom/prometheus: v2\.[0-9]+\.[0-9]+
  prom/pushgateway: v1\.[0-9]+\.[0-9]+
  sonatype/nexus3: 3\.[0-9]+\.[0-9]+
```

and create a `~/.dip/creds.yml` file:

```bash
slack_channel_id: someSlackChannelID
slack_token: some-token
```

or for k8s:

```bash
apiVersion: v1
kind: ConfigMap
metadata:
  name: dip
  namespace: dip
data:
  config.yml: |-
    ---
    dip_images:
      docker.io/alpine: 3\.[0-9]+\.[0-9]+
      elastic/elasticsearch: 7\.[0-9]+\.[0-9]+
```

and

```bash
apiVersion: v1
kind: Secret
metadata:
  name: dip
  namespace: dip
stringData:
  creds.yml: |-
    ---
    slack_channel_id: some-id
    slack_token: some-token
```

Note: follow these steps to create
[a Slack Token](https://github.com/030/sasm#create-a-slack-token).

## latest

### alpine

```bash
dip image --name=alpine --regex="(\d+\.){2}\d"
```

### minio

```bash
dip image --name=minio/minio --regex="RELEASE\.2019.*"
```

### nexus

```bash
dip image --name=sonatype/nexus3 --regex="(\d+\.){2}\d"
```

### nginx

```bash
dip image --name=nginx --regex=".*(\d+\.){2}\d-alpine$"
```

### sonarqube

```bash
dip image --name=sonarqube --regex=".*-community$"
```

### traefik

```bash
dip image --name=traefik --regex="^v(\d+\.){1,2}\d+$"
```

### ubuntu

```bash
dip image --name=ubuntu --regex="^xenial.*"
```

## dockerfile

Use `-dockerfile` to check whether the image that is defined in the `FROM`
should be updated. If the command is run in the Continuous Integration (CI),
the pipeline will fail as an exit 1 is returned if an image is outdated.

### golang

```bash
dip image --name=golang --regex="([0-9]+\.){2}[0-9]+$" --dockerfile
```

### adoptopenjdk

```bash
dip image --name=adoptopenjdk --regex="14.*-jre-hotspot-bionic" --dockerfile
```

## docker

[![dockeri.co](https://dockeri.co/image/utrecht/dip)](https://hub.docker.com/r/utrecht/dip)

```bash
docker run utrecht/dip:4.0.2 dip image --name=grafana/grafana --regex=^7\.5\.7$
```

will return:

```bash
7.5.7
```
