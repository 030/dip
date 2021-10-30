# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [4.0.2] - 2021-10-30

### Changed

- some library versions.

## [4.0.1] - 2021-09-26

### Added

- Let CI fail if image in k8sfile is outdated.

## [4.0.0] - 2021-09-25

### Added

- Send a message to Slack if Dockerfile contains an outdated image.

### Changed

- Different subcommands due to use of Cobra CLI.

## [3.0.3] - 2021-09-07

### Fixed

- Update number was not checked in tags

## [3.0.2] - 2021-08-27

### Changed

- Alpine 3.14.2

## [3.0.1] - 2021-08-22

### Changed

- go v16 to v17.1
- use official go slack library

## [3.0.0] - 2021-08-16

### Changed

- Separate creds and config file.

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

[Unreleased]: https://github.com/030/dip/compare/4.0.2...HEAD
[4.0.2]: https://github.com/030/dip/compare/4.0.1...4.0.2
[4.0.1]: https://github.com/030/dip/compare/4.0.0...4.0.1
[4.0.0]: https://github.com/030/dip/compare/3.0.3...4.0.0
[3.0.3]: https://github.com/030/dip/compare/3.0.2...3.0.3
[3.0.2]: https://github.com/030/dip/compare/3.0.1...3.0.2
[3.0.1]: https://github.com/030/dip/compare/3.0.0...3.0.1
[3.0.0]: https://github.com/030/dip/compare/2.2.0...3.0.0
[2.2.0]: https://github.com/030/dip/compare/2.1.6...2.2.0
[2.1.6]: https://github.com/030/dip/compare/2.1.5...2.1.6
[2.1.5]: https://github.com/030/dip/compare/2.1.4...2.1.5
[2.1.4]: https://github.com/030/dip/compare/2.1.3...2.1.4
[2.1.3]: https://github.com/030/dip/compare/2.1.2...2.1.3
[2.1.2]: https://github.com/030/dip/releases/tag/2.1.2
