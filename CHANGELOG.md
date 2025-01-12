# Changelog

## [0.3.0](https://github.com/scottmckendry/cl-parse/compare/v0.2.0...v0.3.0) (2025-01-12)


### Features

* **parser:** add flags for latest and specific releases ([#6](https://github.com/scottmckendry/cl-parse/issues/6)) ([4124c24](https://github.com/scottmckendry/cl-parse/commit/4124c246e90080c836e14b9025a953d7131283c4)), closes [#3](https://github.com/scottmckendry/cl-parse/issues/3)
* **parser:** optionally include commit body ([dc7c33e](https://github.com/scottmckendry/cl-parse/commit/dc7c33e03e0f46091016c0405bfc3ef8ac27d6ee)), closes [#1](https://github.com/scottmckendry/cl-parse/issues/1)


### Bug Fixes

* **cmd:** check if in git repo before trying to fetch commits ([8d2e7bc](https://github.com/scottmckendry/cl-parse/commit/8d2e7bc6e28fbd291c984c1414d7d58c71211469))
* **parser:** more robust SHA detection ([ecceb64](https://github.com/scottmckendry/cl-parse/commit/ecceb64e8a0d0afb695d352293c3c4027ec3ed50))
* **parser:** only remove closing paren if exists ([22822a9](https://github.com/scottmckendry/cl-parse/commit/22822a9f19442b51d952b550e73ad3c229583371))

## [0.2.0](https://github.com/scottmckendry/cl-parse/compare/v0.1.0...v0.2.0) (2025-01-12)


### Features

* **ci:** add release-please config and manifest ([8f35e6e](https://github.com/scottmckendry/cl-parse/commit/8f35e6ee07f85777d590d37fef28a9e8434c0f27))
* **cmd:** add cobra & default cmd definition ([fd63c7c](https://github.com/scottmckendry/cl-parse/commit/fd63c7c7ab30a402c5332e8da2f4e77bcb8d084f))


### Bug Fixes

* **parser:** handle versions without compare urls ([95e9aed](https://github.com/scottmckendry/cl-parse/commit/95e9aedafd0ebbb75256048faf55496e20c4358e))

## 0.1.0 (2025-01-12)


### Features

* add basic parsing logic for standard changelog formats ([163661c](https://github.com/scottmckendry/cl-parse/commit/163661c06dc0d275325f9247bbf42e99400ef909))
* **ci:** add ci workflow ([1256d0d](https://github.com/scottmckendry/cl-parse/commit/1256d0d5d0b3ea1741dbc9e3b352be9b7b1de745))
* **parser:** add some basic test cases ([045e395](https://github.com/scottmckendry/cl-parse/commit/045e395bbe6ef8693c9b64fe094903be73407e39))
