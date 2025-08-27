## 1.0.0 (2025-08-27)

### ⚠ BREAKING CHANGES

* **go-modules:** Users must update import paths and dependencies to the
new module path. Run `go get github.com/conforma/go-gather@latest` and
`go mod tidy` to resolve.

Signed-off-by: Rob Nester <rnester@redhat.com>

### Features

* **detector:** Add new Detector functionality. ([5b89b1d](https://github.com/robnester-rh/go-gather/commit/5b89b1d25470f5545496aa3965c2a3c69c62992a))

### Bug Fixes

* Added missing "v" in tagFormat in .releaserc ([fa2b652](https://github.com/robnester-rh/go-gather/commit/fa2b652ecb9552efc848631224ea928bc37ea793))
* call ClassifyURI on destinations in copyFile ([aacce9f](https://github.com/robnester-rh/go-gather/commit/aacce9f74ac9f3d151326938a6b12107f4783631))
* **deps:** update module github.com/go-git/go-git/v5 to v5.13.2 ([a34f303](https://github.com/robnester-rh/go-gather/commit/a34f303f7ab8cab26dc2ba8b0a93c7e4e05de698))
* **deps:** update module github.com/stretchr/testify to v1.10.0 ([24a578f](https://github.com/robnester-rh/go-gather/commit/24a578f8b72c419c6d0afa4322792cc4788c2683))
* **detector-tests:** Fixed detector tests. ([00ac820](https://github.com/robnester-rh/go-gather/commit/00ac820fcfebad39bf4c93ddf71e5c32cc954a6e))
* **go-modules:** update module path ([1d67df5](https://github.com/robnester-rh/go-gather/commit/1d67df53a1c8560e9607e4a898c8e268161c87a1))
* **metadata:** fix returns  on `Gather()` methods ([8675a30](https://github.com/robnester-rh/go-gather/commit/8675a3085a3c1b546978cc7de7e99cecb876aeed))
* TLS or not determination ([7ac9200](https://github.com/robnester-rh/go-gather/commit/7ac92008c381e8a198e18df011328e6cb708f657))
* update Matcher() function in `gather/file` ([4fde473](https://github.com/robnester-rh/go-gather/commit/4fde473f7dc657aec64d7145b82f5c5d48912d8d))

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
