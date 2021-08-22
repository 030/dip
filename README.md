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
Usage of ./dip:
  -debug
        Whether debug mode should be enabled.
  -dockerfile
        Whether dockerfile should be checked.
  -image string
        Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.
  -k8s
        Whether images are up to date in a k8s or openshift cluster.
  -latest string
        The regex to get the latest tag, e.g. "xenial-\d.*".
  -version
        The version of DIP.
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
kind: Secret
metadata:
  name: dip-config
  namespace: dip
stringData:
  config.yml: |-
    ---
    dip_images:
      docker.io/alpine: 3\.[0-9]+\.[0-9]+
      elastic/elasticsearch: 7\.[0-9]+\.[0-9]+
  creds.yml: |-
    ---
    slack_token: some-token
```

## latest

### alpine

```bash
./dip -image alpine -latest "(\d+\.){2}\d"
```

### minio

```bash
./dip -image minio/minio -latest "RELEASE\.2019.*"
```

### nexus

```bash
./dip -image sonatype/nexus3 -latest "(\d+\.){2}\d"
```

### nginx

```bash
./dip -image nginx -latest ".*(\d+\.){2}\d-alpine$"
```

### sonarqube

```bash
./dip -image sonarqube -latest ".*-community$"
```

### traefik

```bash
./dip --image=traefik --latest="^v(\d+\.){1,2}\d+$"
```

### ubuntu

```bash
./dip -image ubuntu -latest "^xenial.*"
```

## dockerfile

Use `-dockerfile` to check whether the image that is defined in the `FROM`
should be updated. If the command is run in the Continuous Integration (CI),
the pipeline will fail as an exit 1 is returned if an image is outdated.

### golang

```bash
./dip -image=golang -latest="([0-9]+\.){2}[0-9]+$" -dockerfile
```

### adoptopenjdk

```bash
./dip -image=adoptopenjdk -latest="14.*-jre-hotspot-bionic" -dockerfile
```

## docker

[![dockeri.co](https://dockeri.co/image/utrecht/dip)](https://hub.docker.com/r/utrecht/dip)

```bash
docker run utrecht/dip:2.2.0 dip -image=grafana/grafana -latest=^7\.5\.7$
```

will return:

```bash
7.5.7
```
