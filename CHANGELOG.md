# Changelog

> **This changelog is no longer maintained.** For release notes after v1.1.2,
> see the [GitHub Releases](https://github.com/conforma/go-gather/releases) page.

## [1.1.2](https://github.com/conforma/go-gather/compare/v1.1.1...v1.1.2) (2026-04-13)

### Bug Fixes

* **deps:** update actions/upload-artifact action to v7.0.1 ([3e6fec5](https://github.com/conforma/go-gather/commit/3e6fec5dcfe11f73dbb5cd425b65c7fc98fbeb43))
* **deps:** update github actions ([b2496c3](https://github.com/conforma/go-gather/commit/b2496c3c866489f4ef23341cda553182cb7be952))

## [1.1.1](https://github.com/conforma/go-gather/compare/v1.1.0...v1.1.1) (2026-04-07)

### Bug Fixes

* **ci:** add local Renovate overrides for semantic-release compatibility ([9626eab](https://github.com/conforma/go-gather/commit/9626eabbef20d26e6801c3c36d866fe156b2349f))

## [1.1.0](https://github.com/conforma/go-gather/compare/v1.0.2...v1.1.0) (2026-03-31)

### Features

* document library capabilities and security improvements ([d42d23e](https://github.com/conforma/go-gather/commit/d42d23eb1b271bcf2bb3166ef9e17e0164557e57)), closes [#302](https://github.com/conforma/go-gather/issues/302)

### Bug Fixes

* address CodeRabbit review feedback ([7505688](https://github.com/conforma/go-gather/commit/75056882c68d8e6c7bf3474eb2c99dce459bcb93)), closes [#299](https://github.com/conforma/go-gather/issues/299) [#300](https://github.com/conforma/go-gather/issues/300)
* address CodeRabbit review feedback on CI/CD workflow ([bb12131](https://github.com/conforma/go-gather/commit/bb12131a6de3badac8026d6832df3857897f50df)), closes [#298](https://github.com/conforma/go-gather/issues/298)
* **ci:** pass GITHUB_TOKEN to semantic-release ([f59d267](https://github.com/conforma/go-gather/commit/f59d267ff27da7213f6c5b116609248e5681cf02)), closes [#302](https://github.com/conforma/go-gather/issues/302)
* **ci:** remove changelog/git plugins incompatible with branch protection ([9ab663b](https://github.com/conforma/go-gather/commit/9ab663b2851c7b6b8ed927ac7dd79c36296cef7d)), closes [#302](https://github.com/conforma/go-gather/issues/302)
* **deps:** update go-git to v5.17.1 [SECURITY] ([3a79d5b](https://github.com/conforma/go-gather/commit/3a79d5b81b9c9eb302120c59eaafbef6daac0ea4)), closes [#306](https://github.com/conforma/go-gather/issues/306)
* fix URL routing bugs and improve code quality ([00ac929](https://github.com/conforma/go-gather/commit/00ac9299ced918aaabf73b8fda15959bf5086926)), closes [#300](https://github.com/conforma/go-gather/issues/300) [#302](https://github.com/conforma/go-gather/issues/302)
* harden CI/CD pipeline and release workflow ([6bbf6c2](https://github.com/conforma/go-gather/commit/6bbf6c22a41226e28110a766e78981878c730545)), closes [#298](https://github.com/conforma/go-gather/issues/298)
* harden OCI registry patterns and wrap close errors ([d8639db](https://github.com/conforma/go-gather/commit/d8639dbca1a2028820130f01211e702bc7784350)), closes [#300](https://github.com/conforma/go-gather/issues/300)
* improve error handling, fix SCP URL parsing, and update documentation ([c0f5dc4](https://github.com/conforma/go-gather/commit/c0f5dc4ece34a89819ca5ef9b2fc2b0f5c7ae7b3)), closes [#300](https://github.com/conforma/go-gather/issues/300) [#302](https://github.com/conforma/go-gather/issues/302)
* migrate to golangci-lint v2 and fix Makefile targets ([8991327](https://github.com/conforma/go-gather/commit/8991327a0809f4659e3335ac43ed4aeb16d05e22)), closes [#297](https://github.com/conforma/go-gather/issues/297) [#298](https://github.com/conforma/go-gather/issues/298)
* nil metadata on deferred Close() errors ([9a7b9b7](https://github.com/conforma/go-gather/commit/9a7b9b7640a9bd92bd50910f21bc89fbce6309ff)), closes [#299](https://github.com/conforma/go-gather/issues/299) [#300](https://github.com/conforma/go-gather/issues/300)
* open source before destination in FileSaver and parse URLs in HTTP Matcher ([1f8fe11](https://github.com/conforma/go-gather/commit/1f8fe1170defb70ec68f25484575f9e8475a67d5)), closes [#304](https://github.com/conforma/go-gather/issues/304)
* standardize Gather() to return pointer metadata ([def09ee](https://github.com/conforma/go-gather/commit/def09ee11223013f778cf3b9ce374e6ffff6dabb)), closes [#299](https://github.com/conforma/go-gather/issues/299)

## [1.0.2](https://github.com/conforma/go-gather/compare/v1.0.1...v1.0.2) (2025-03-06)

### Bug Fixes

* **metadata:** fix returns  on `Gather()` methods ([8675a30](https://github.com/conforma/go-gather/commit/8675a3085a3c1b546978cc7de7e99cecb876aeed))

## [1.0.1](https://github.com/conforma/go-gather/compare/v1.0.0...v1.0.1) (2025-02-27)

### Bug Fixes

* update Matcher() function in `gather/file` ([4fde473](https://github.com/conforma/go-gather/commit/4fde473f7dc657aec64d7145b82f5c5d48912d8d))

## [1.0.0](https://github.com/conforma/go-gather/compare/v0.1.2...v1.0.0) (2025-02-12)

### ⚠ BREAKING CHANGES

* **go-modules:** Users must update import paths and dependencies to the
new module path. Run `go get github.com/conforma/go-gather@latest` and
`go mod tidy` to resolve.

Signed-off-by: Rob Nester <rnester@redhat.com>

### Bug Fixes

* **go-modules:** update module path ([1d67df5](https://github.com/conforma/go-gather/commit/1d67df53a1c8560e9607e4a898c8e268161c87a1))

## [0.1.1](https://github.com/conforma/go-gather/compare/v0.1.0...v0.1.1) (2025-02-05)

### Bug Fixes

* **detector-tests:** Fixed detector tests. ([00ac820](https://github.com/conforma/go-gather/commit/00ac820fcfebad39bf4c93ddf71e5c32cc954a6e))

## [0.1.0](https://github.com/conforma/go-gather/compare/v0.0.8...v0.1.0) (2025-02-03)

### Features

* **detector:** Add new Detector functionality. ([5b89b1d](https://github.com/conforma/go-gather/commit/5b89b1d25470f5545496aa3965c2a3c69c62992a))

## [0.0.8](https://github.com/conforma/go-gather/compare/v0.0.7...v0.0.8) (2025-02-03)

### Bug Fixes

* **deps:** update module github.com/go-git/go-git/v5 to v5.13.2 ([a34f303](https://github.com/conforma/go-gather/commit/a34f303f7ab8cab26dc2ba8b0a93c7e4e05de698))

## [0.0.7](https://github.com/conforma/go-gather/compare/v0.0.6...v0.0.7) (2025-01-28)

## [0.0.6](https://github.com/conforma/go-gather/compare/v0.0.5...v0.0.6) (2025-01-10)

### Bug Fixes

* TLS or not determination ([7ac9200](https://github.com/conforma/go-gather/commit/7ac92008c381e8a198e18df011328e6cb708f657))

## [0.0.5](https://github.com/conforma/go-gather/compare/v0.0.4...v0.0.5) (2024-11-25)

### Bug Fixes

* **deps:** update module github.com/stretchr/testify to v1.10.0 ([24a578f](https://github.com/conforma/go-gather/commit/24a578f8b72c419c6d0afa4322792cc4788c2683))

## [0.0.4](https://github.com/conforma/go-gather/compare/v0.0.3...v0.0.4) (2024-10-21)

### Bug Fixes

* Added missing "v" in tagFormat in .releaserc ([fa2b652](https://github.com/conforma/go-gather/commit/fa2b652ecb9552efc848631224ea928bc37ea793))
* call ClassifyURI on destinations in copyFile ([aacce9f](https://github.com/conforma/go-gather/commit/aacce9f74ac9f3d151326938a6b12107f4783631))

## [0.0.4](https://github.com/conforma/go-gather/compare/v0.0.3...0.0.4) (2024-10-21)

### Bug Fixes

* call ClassifyURI on destinations in copyFile ([aacce9f](https://github.com/conforma/go-gather/commit/aacce9f74ac9f3d151326938a6b12107f4783631))
