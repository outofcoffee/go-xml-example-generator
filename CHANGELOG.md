# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [0.4.2] - 2025-04-15
### Changed
- chore(deps): bumps golang.org/x/net to v0.36.0.

## [0.4.1] - 2025-01-29
### Changed
- refactor: removes unneeded protoTree traversal.

## [0.4.0] - 2025-01-29
### Added
- feat: supports element refs to imported schemas.

## [0.3.1] - 2025-01-28
### Changed
- ci: removes unused linter.

## [0.3.0] - 2025-01-28
### Added
- feat: respects default element form.

### Changed
- test: moves simple XSD sample to subdirectory.

## [0.2.5] - 2025-01-24
### Fixed
- fix: use existing complex type prefix.

## [0.2.4] - 2025-01-24
### Fixed
- fix: removes duplicate prefix for complextype closing tag.

## [0.2.3] - 2025-01-24
### Fixed
- fix: removes duplicate prefix for complextype name.

## [0.2.2] - 2025-01-11
### Changed
- chore: bumps since action version.

## [0.2.1] - 2025-01-11
### Changed
- ci: create release on tag push.

## [0.2.0] - 2025-01-11
### Added
- feat: adds namespace and prefix support.

## [0.1.2] - 2025-01-10
### Changed
- build(deps): bump golang.org/x/net (#1)
- docs: adds CI badge.
- docs: simplifies introduction.

## [0.1.1] - 2025-01-10
### Changed
- build: adds since config.
- build: tune linter.
- ci: adds GitHub Actions config.
- docs: adds changelog.
- refactor: don't expose the parser via the API.
- refactor: removes old package.

## [0.1.0] - 2025-01-09
### Other
- initial commit.
