# CHANGELOG v1.0.x <!-- omit in toc -->

**Table of contents**

- [Versions](#versions)
- [Changes](#changes)
  - [Breaking changes](#breaking-changes)
  - [New features](#new-features)
  - [Bug fix, Security fix](#bug-fix-security-fix)
  - [Other changes](#other-changes)
- [Dependencies](#dependencies)
  - [Added](#added)
  - [Changed](#changed)
  - [Removed](#removed)
- [Migration guides](#migration-guides)

## Versions

- AILERON Gateway: `v1.0.0`
- Go: `v1.24`
- protoc: `v29.0`
- protoc-gen-go: `v1.36.4`

## Changes

### Breaking changes

- [#105](https://github.com/aileron-gateway/aileron-gateway/pull/105): Move kernel/testutil package to internal/testutil (@k7a-tomohiro)
- [#108](https://github.com/aileron-gateway/aileron-gateway/pull/108): Move kernel/hash, kernel/encrypt packages to internal/hash, internal/encrypt (@k7a-tomohiro)
- [#109](https://github.com/aileron-gateway/aileron-gateway/pull/109): Move kernel/txtutil package to internal/txtutil (@k7a-tomohiro)
- [#110](https://github.com/aileron-gateway/aileron-gateway/pull/110): Move kernel/kvs package to internal/kvs (@k7a-tomohiro)
- [#111](https://github.com/aileron-gateway/aileron-gateway/pull/111): Move kernel/network package to internal/network (@k7a-tomohiro)
- [#112](https://github.com/aileron-gateway/aileron-gateway/pull/112): Move util/security package to internal/security (@k7a-tomohiro)
- [#114](https://github.com/aileron-gateway/aileron-gateway/pull/114): Move util/session package to kernel/session (@k7a-tomohiro)

### New features

_Nothing has changed._

### Bug fix, Security fix

_Nothing has changed._

### Other changes

_Nothing has changed._

## Dependencies

<!--
Changes are generated with go-modiff like this.
`go-modiff --repository github.com/aileron-gateway/aileron-gateway --from v1.0.4 --to main`
-->

### Added

_Nothing has changed._

### Changed

_Nothing has changed._

### Removed

_Nothing has changed._

## Migration guides

_Migration is not required._
