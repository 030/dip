# DIP

Docker Image Patrol (DIP) keeps docker images up-to-date.

## Usage

```
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

```
dip -image nginx:1.17.5-alpine -registry quay.io/some-org/
```

If the image is absent, true would be returned. False indicates that the image
is available in a docker-registry. Note: do not omit that last forward slash.