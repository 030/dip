# DIP

[![CI](https://github.com/030/dip/workflows/Go/badge.svg?event=push)](https://github.com/030/dip/actions?query=workflow%3AGo)
[![GoDoc Widget]][godoc]
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/030/dip)
[![Go Report Card](https://goreportcard.com/badge/github.com/030/dip)](https://goreportcard.com/report/github.com/030/dip)
[![StackOverflow SE Questions](https://img.shields.io/stackexchange/stackoverflow/t/dip.svg?logo=stackoverflow)](https://stackoverflow.com/tags/dip)
[![DevOps SE Questions](https://img.shields.io/stackexchange/devops/t/dip.svg?logo=stackexchange)](https://devops.stackexchange.com/tags/dip)
[![ServerFault SE Questions](https://img.shields.io/stackexchange/serverfault/t/dip.svg?logo=serverfault)](https://serverfault.com/tags/dip)
![Docker Pulls](https://img.shields.io/docker/pulls/utrecht/dip.svg)
[![dip on stackoverflow](https://img.shields.io/badge/stackoverflow-community-orange.svg?longCache=true&logo=stackoverflow)](https://stackoverflow.com/tags/dip)
![Issues](https://img.shields.io/github/issues-raw/030/dip.svg)
![Pull requests](https://img.shields.io/github/issues-pr-raw/030/dip.svg)
![Total downloads](https://img.shields.io/github/downloads/030/dip/total.svg)
![GitHub forks](https://img.shields.io/github/forks/030/dip?label=fork&style=plastic)
![GitHub watchers](https://img.shields.io/github/watchers/030/dip?style=plastic)
![GitHub stars](https://img.shields.io/github/stars/030/dip?style=plastic)
![License](https://img.shields.io/github/license/030/dip.svg)
![Repository Size](https://img.shields.io/github/repo-size/030/dip.svg)
![Contributors](https://img.shields.io/github/contributors/030/dip.svg)
![Commit activity](https://img.shields.io/github/commit-activity/m/030/dip.svg)
![Last commit](https://img.shields.io/github/last-commit/030/dip.svg)
![Release date](https://img.shields.io/github/release-date/030/dip.svg)
![Latest Production Release Version](https://img.shields.io/github/release/030/dip.svg)

<!-- [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=bugs)](https://sonarcloud.io/dashboard?id=030_dip)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=code_smells)](https://sonarcloud.io/dashboard?id=030_dip)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=coverage)](https://sonarcloud.io/dashboard?id=030_dip)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=030_dip)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=ncloc)](https://sonarcloud.io/dashboard?id=030_dip)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=030_dip)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=alert_status)](https://sonarcloud.io/dashboard?id=030_dip)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=030_dip)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=security_rating)](https://sonarcloud.io/dashboard?id=030_dip)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=sqale_index)](https://sonarcloud.io/dashboard?id=030_dip)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=030_dip&metric=vulnerabilities)]\
(https://sonarcloud.io/dashboard?id=030_dip) -->

<!-- [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2810/badge)]\
(https://bestpractices.coreinfrastructure.org/projects/2810)  -->

[![codecov](https://codecov.io/gh/030/dip/branch/main/graph/badge.svg)](https://codecov.io/gh/030/dip)
[![BCH compliance](https://bettercodehub.com/edge/badge/030/dip?branch=main)](https://bettercodehub.com/results/030/dip)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-web.svg)](https://golangci.com/r/github.com/030/dip)
[![Chocolatey](https://img.shields.io/chocolatey/dt/dip)](https://chocolatey.org/packages/dip)
[![kdiutd](https://snapcraft.io/kdiutd/badge.svg)](https://snapcraft.io/kdiutd)

<!-- [![codebeat badge](https://codebeat.co/badges/X)]\
(https://codebeat.co/projects/github-com-030-dip-main) -->

[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)

[godoc]: https://godoc.org/github.com/030/dip
[godoc widget]: https://godoc.org/github.com/030/dip?status.svg

<a href="https://dip.releasesoftwaremoreoften.com">\
<img src="https://github.com/030/dip/raw/main/assets/logo/logo.png" width="100"></a>

Docker Image Patrol (DIP) keeps docker images up-to-date.

## Installation

Keep Docker Images Up To Date (KDIUTD)

```bash
sudo snap install kdiutd
```

## Usage

```bash
Usage:
  dip [flags]
  dip [command]

Available Commands:
  help        Help about any command
  image       A brief description of your command

Flags:
      --configCredHome string   config and cred file home directory (default is $HOME/.dip)
      --debug                   debugging mode
  -h, --help                    help for dip
  -v, --version                 version for dip

Use "dip [command] --help" for more information about a command.
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
docker run utrecht/dip:4.1.0 dip image --name=grafana/grafana --regex=^7\.5\.7$
```

will return:

```bash
7.5.7
```

## updateDockerfile

Use the `--updateDockerfile` to check and update the image that is defined in
the `FROM` inside a Dockerfile.

### golang alpine builder

```bash
dip image --name=golang --regex="^([0-9]+\.){2}[0-9]-alpine([0-9]+\.)[0-9]{2}$" --updateDockerfile
```
