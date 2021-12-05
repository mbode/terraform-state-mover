# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
### Changed
### Deprecated
### Removed
### Fixed
- Fix passing of extra arguments to `terraform plan` call
### Security


## [0.4.2] - 2021-12-05
### Fixed
- Fix deprecation warning about `bottle :unneeded` in homebrew tap

## [0.4.1] - 2021-07-16
### Fixed
- Fix `brew style` issue

## [0.4.0] - 2021-07-16
### Added
- Search in resource addresses ([#20](https://github.com/mbode/terraform-state-mover/pull/20), contributed by [@pascal-hofmann](https://github.com/pascal-hofmann))
- Print terraform state mv commands on interrupt ([#22](https://github.com/mbode/terraform-state-mover/pull/22), contributed by [@pascal-hofmann](https://github.com/pascal-hofmann))

## [0.3.0] - 2020-10-04
### Added
- Delay flag to prevent hitting rate limits with remote state ([#8](https://github.com/mbode/terraform-state-mover/pull/8), contributed by [@xanonid](https://github.com/xanonid))
- Flags for verbose output and dry-run mode ([#8](https://github.com/mbode/terraform-state-mover/pull/8), contributed by [@xanonid](https://github.com/xanonid))

## [0.2.0] - 2020-09-30
### Added
- Publish a homebrew formula

## [0.1.0] - 2020-09-04
### Added
- Initial release

[Unreleased]: https://github.com/mbode/terraform-state-mover/compare/0.4.2...HEAD
[0.4.2]: https://github.com/mbode/terraform-state-mover/compare/0.4.1...0.4.2
[0.4.1]: https://github.com/mbode/terraform-state-mover/compare/0.4.0...0.4.1
[0.4.0]: https://github.com/mbode/terraform-state-mover/compare/0.3.0...0.4.0
[0.3.0]: https://github.com/mbode/terraform-state-mover/compare/0.2.0...0.3.0
[0.2.0]: https://github.com/mbode/terraform-state-mover/compare/0.1.0...0.2.0
[0.1.0]: https://github.com/mbode/terraform-state-mover/releases/tag/0.1.0