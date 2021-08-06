# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.2.0] - 2021-08-02

### Added

- Check whether docker images in k8s and openshift clusters are outdated.

## [2.1.6] - 2021-06-21

### Fixed

- Consulting <registry.hub.docker.com/> returns 503.

## [2.1.5] - 2021-05-29

### Added

- Docker image.

### Fixed

- Resolve issue in adoptopenjdk sorting.

## [2.1.4] - 2021-03-26

### Added

- Snapcraft package.

### Fixed

- Latest update version was not returned.

### Removed

- `lzo` compression as it results in larger snap packages.

## [2.1.3] - 2021-03-5

### Fixed

- Return latest tag.

[Unreleased]: https://github.com/030/dip/compare/2.2.0...HEAD
[2.2.0]: https://github.com/030/dip/compare/2.1.6...2.2.0
[2.1.6]: https://github.com/030/dip/compare/2.1.5...2.1.6
[2.1.5]: https://github.com/030/dip/compare/2.1.4...2.1.5
[2.1.4]: https://github.com/030/dip/compare/2.1.3...2.1.4
[2.1.3]: https://github.com/030/dip/compare/2.1.2...2.1.3
[2.1.2]: https://github.com/030/dip/releases/tag/2.1.2
