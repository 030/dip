# DIP

[![GoDoc Widget](https://godoc.org/github.com/030/dip?status.svg)](https://godoc.org/github.com/030/dip)
[![Go Report Card](https://goreportcard.com/badge/github.com/030/dip)](https://goreportcard.com/report/github.com/030/dip)
[![Build Status](https://travis-ci.org/030/dip.svg?branch=master)](https://travis-ci.org/030/dip)
[![DevOps SE Questions](https://img.shields.io/stackexchange/devops/t/dip.svg)](https://devops.stackexchange.com/questions/tagged/dip)

Docker Image Patrol (DIP) keeps docker images up-to-date.

## Usage

```bash
Usage of dip:
  -debug
        Whether debug mode should be enabled
  -image string
        The origin of the image, e.g. nginx:1.17.5-alpine
  -registry string
        To what destination the image should be transferred, e.g. quay.io/some-org
exit status 2
```

### Absent

Check whether a docker-image resides in a docker-registry:

```bash
dip -image nginx:1.17.5-alpine -registry quay.io/some-org/
```

An ```exit 0``` will be returned if the image is absent and an ```exit 1``` is
applicable if it already exists to prevent that the tag gets overwritten.

## preserve

### nginx

```bash
go run main.go -image nginx -tagType ".*(\d+\.){2}\d-alpine$"
```

### ubuntu

```bash
go run main.go -image ubuntu -tagType "xenial-\d.*"
```
