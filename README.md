# DIP

[![GoDoc Widget](https://godoc.org/github.com/030/dip?status.svg)](https://godoc.org/github.com/030/dip)
[![Go Report Card](https://goreportcard.com/badge/github.com/030/dip)](https://goreportcard.com/report/github.com/030/dip)
[![Build Status](https://travis-ci.org/030/dip.svg?branch=master)](https://travis-ci.org/030/dip)
[![DevOps SE Questions](https://img.shields.io/stackexchange/devops/t/dip.svg)](https://devops.stackexchange.com/questions/tagged/dip)
![Issues](https://img.shields.io/github/issues-raw/030/n3dr.svg)
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

<a href="https://dip.releasesoftwaremoreoften.com"><img src="https://github.com/030/dip/raw/master/assets/logo/logo.png" width="100"></a>

Docker Image Patrol (DIP) keeps docker images up-to-date.

Note: use `-official` if an image is `official` according to dockerhub.

## Usage

```bash
Usage of ./dip:
  -debug
        Whether debug mode should be enabled.
  -dockerfile
        Whether dockerfile should be checked.
  -image string
        Find an image on dockerhub, e.g. nginx:1.17.5-alpine or nginx.
  -latest string
        The regex to get the latest tag, e.g. "xenial-\d.*".
  -official
        Use this parameter if an image is official according to dockerhub.
  -version
        The version of DIP.
```

## latest

### alpine

```bash
./dip -image alpine -latest "(\d+\.){2}\d" --official
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
./dip -image nginx -latest ".*(\d+\.){2}\d-alpine$" -official
```

### sonarqube

```bash
./dip -image sonarqube -latest ".*-community$"
```

### traefik

```bash
./dip --image=traefik --latest="^v(\d+\.){1,2}\d+$" --official
```

### ubuntu

```bash
./dip -image ubuntu -latest "^xenial.*" -official
```

## dockerfile

Use `-dockerfile` to check whether the image that is defined in the `FROM`
should be updated. If the command is run in the Continuous Integration (CI),
the pipeline will fail as an exit 1 is returned if an image is outdated.

### golang

```bash
./dip -image=golang -latest="([0-9]+\.){2}[0-9]+$" -dockerfile -official
```

### adoptopenjdk

```bash
./dip -image=adoptopenjdk -latest="14.*-jre-hotspot-bionic" -dockerfile -official
```
