# Changelog

## [0.5.0](https://github.com/scottmckendry/cl-parse/compare/v0.4.0...v0.5.0) (2025-06-15)


### Features

* **deps:** update module github.com/go-git/go-git/v5 ( v5.13.1 → v5.16.2 ) ([7b725c0](https://github.com/scottmckendry/cl-parse/commit/7b725c0da3bb54f90904d88d1d0ecc1bcba22c42))
* **deps:** update module github.com/spf13/cobra ( v1.8.1 → v1.9.1 ) ([846213d](https://github.com/scottmckendry/cl-parse/commit/846213d4582d81de32e57747cb6afd021ae6892c))
* support nixos with flake ([1a8be47](https://github.com/scottmckendry/cl-parse/commit/1a8be47b621ba938e5c9d72e348bd7f2e388a943))


### Bug Fixes

* **deps:** update module github.com/pelletier/go-toml/v2 ( v2.2.3 → v2.2.4 ) ([#14](https://github.com/scottmckendry/cl-parse/issues/14)) ([08faab1](https://github.com/scottmckendry/cl-parse/commit/08faab129f95132a3c0b646158f36178a0706f3a))
* **origin:** support ssh git urls ([d09f5de](https://github.com/scottmckendry/cl-parse/commit/d09f5de99763f0f461cd0842587f7dc66f8ebb05))
* **parser:** support azure devops style prs ([1a33648](https://github.com/scottmckendry/cl-parse/commit/1a33648112c7a9f0f72d3446fb8a9c23cd4f4d02))
* update pinned deps ([f3c3fb0](https://github.com/scottmckendry/cl-parse/commit/f3c3fb05f6ee597a907bc6a5b451aba56bb84f7d))

## [0.4.0](https://github.com/scottmckendry/cl-parse/compare/v0.3.0...v0.4.0) (2025-01-14)


### Features

* **cmd:** `format` option with new YAML & TOML outputs ([7ffb283](https://github.com/scottmckendry/cl-parse/commit/7ffb28361ceb950ebdb5483cfdc9f800181f214f))
* **origin:** add gitlab support ([282dc8a](https://github.com/scottmckendry/cl-parse/commit/282dc8a4e502a2901f2597baf65fd88a4b5147a7))
* **origin:** add support for Azure DevOps workitems ([e297c7a](https://github.com/scottmckendry/cl-parse/commit/e297c7a5f1bba68e0f0a489870c7e5410abaa7c4)), closes [#2](https://github.com/scottmckendry/cl-parse/issues/2)
* **origin:** add support for github issue lookup ([539c4cd](https://github.com/scottmckendry/cl-parse/commit/539c4cdf5fabcdef93dbf3eca6200b09f6c68683)), closes [#2](https://github.com/scottmckendry/cl-parse/issues/2)
* **parser:** add RelatedItems property ([fed2a7b](https://github.com/scottmckendry/cl-parse/commit/fed2a7b4d0824ac6d04f96cda2f37f2cd80e9d31))


### Bug Fixes

* **parser:** handle "closes #X" strings ([e44ef80](https://github.com/scottmckendry/cl-parse/commit/e44ef80328f7284436ab7f81562df7aa9e11d6af))

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
