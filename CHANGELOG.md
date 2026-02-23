# Changelog

## [1.17.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.17.0...v1.17.1) (2026-02-23)


### Bug Fixes

* preallocate slices in extractModuleNames and createModulesFromNames ([#111](https://github.com/CloudNationHQ/az-cn-go-validor/issues/111)) ([2d886a4](https://github.com/CloudNationHQ/az-cn-go-validor/commit/2d886a4bce41e560bc3add10605a5103bc054331)), closes [#89](https://github.com/CloudNationHQ/az-cn-go-validor/issues/89)
* remove redundant zero value initialization (fixes [#91](https://github.com/CloudNationHQ/az-cn-go-validor/issues/91)) ([#112](https://github.com/CloudNationHQ/az-cn-go-validor/issues/112)) ([7c505ba](https://github.com/CloudNationHQ/az-cn-go-validor/commit/7c505badce023a42e15921f1a781419e0b5cc285))

## [1.17.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.16.1...v1.17.0) (2026-02-21)


### Features

* small overall refactor ([#109](https://github.com/CloudNationHQ/az-cn-go-validor/issues/109)) ([6cdbec4](https://github.com/CloudNationHQ/az-cn-go-validor/commit/6cdbec4882fc330aa35771eba128b4811b5d1566))

## [1.16.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.16.0...v1.16.1) (2026-02-21)


### Bug Fixes

* align moduleDiscoverer and TestRunner interfaces with implementation ([#99](https://github.com/CloudNationHQ/az-cn-go-validor/issues/99)) ([d3cd9cb](https://github.com/CloudNationHQ/az-cn-go-validor/commit/d3cd9cb59900cd7a442ebf5dec739e9518718289))
* make global config initialization thread-safe with sync.Once ([#97](https://github.com/CloudNationHQ/az-cn-go-validor/issues/97)) ([fce6c0e](https://github.com/CloudNationHQ/az-cn-go-validor/commit/fce6c0ee4b8028b239e929b23f69ca4287ce036b))
* remove unreachable return statement after t.Fatal() ([#100](https://github.com/CloudNationHQ/az-cn-go-validor/issues/100)) ([908d54b](https://github.com/CloudNationHQ/az-cn-go-validor/commit/908d54b4843e1adf5068c2ced598b68c1fa1a2ab))
* remove unused ModuleProcessor interface ([#101](https://github.com/CloudNationHQ/az-cn-go-validor/issues/101)) ([d19f1e1](https://github.com/CloudNationHQ/az-cn-go-validor/commit/d19f1e105ccebbb84c1e5183214f590bf67430a5))
* replace clever bool-to-string map with BoolToStr helper ([#102](https://github.com/CloudNationHQ/az-cn-go-validor/issues/102)) ([7f8ee8c](https://github.com/CloudNationHQ/az-cn-go-validor/commit/7f8ee8cdc8adbc054642cc24f71d5cd41c9d3b44))
* replace regex with idiomatic string operations ([#103](https://github.com/CloudNationHQ/az-cn-go-validor/issues/103)) ([f2ca5e6](https://github.com/CloudNationHQ/az-cn-go-validor/commit/f2ca5e6b8824d6c9a4dcc34825c9793ba5a31a9b))

## [1.16.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.15.0...v1.16.0) (2025-12-05)


### Features

* consolidate test files and improve coverage to 85.7% ([#80](https://github.com/CloudNationHQ/az-cn-go-validor/issues/80)) ([e976560](https://github.com/CloudNationHQ/az-cn-go-validor/commit/e976560b32a4f4d4d95a40497ab9a196adda04d0))


### Bug Fixes

* move GOCACHE env to step level for runner.temp context ([#82](https://github.com/CloudNationHQ/az-cn-go-validor/issues/82)) ([ab60e58](https://github.com/CloudNationHQ/az-cn-go-validor/commit/ab60e587270a4d291f98c7e27c7721b9337f5e0e))

## [1.15.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.14.0...v1.15.0) (2025-12-03)


### Features

* add testing and increase coverage ([#78](https://github.com/CloudNationHQ/az-cn-go-validor/issues/78)) ([6d011a8](https://github.com/CloudNationHQ/az-cn-go-validor/commit/6d011a8e2e56f42f832fabaf5f21985c660372fc))
* **deps:** bump github.com/gruntwork-io/terratest from 0.52.0 to 0.54.0 ([#77](https://github.com/CloudNationHQ/az-cn-go-validor/issues/77)) ([b487ef1](https://github.com/CloudNationHQ/az-cn-go-validor/commit/b487ef1dc0b88b6cecd61e52dc34d9e760f34a8b))
* **deps:** bump golang.org/x/crypto from 0.36.0 to 0.45.0 in /tests ([#75](https://github.com/CloudNationHQ/az-cn-go-validor/issues/75)) ([e792588](https://github.com/CloudNationHQ/az-cn-go-validor/commit/e792588324be017261c77823ed3d31a035effd85))
* **deps:** bump golang.org/x/crypto from 0.41.0 to 0.45.0 ([#76](https://github.com/CloudNationHQ/az-cn-go-validor/issues/76)) ([14b337f](https://github.com/CloudNationHQ/az-cn-go-validor/commit/14b337f34fbdfa87b11d9d61b8d0f94d84e20005))

## [1.14.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.13.0...v1.14.0) (2025-11-13)


### Features

* **deps:** bump github.com/gruntwork-io/terratest from 0.50.0 to 0.52.0 ([#68](https://github.com/CloudNationHQ/az-cn-go-validor/issues/68)) ([9591ce1](https://github.com/CloudNationHQ/az-cn-go-validor/commit/9591ce1047ac95f9c8ec1d9e3a1493e418677f22))


### Bug Fixes

* prevent submodule regex from matching modules outside current repository ([#72](https://github.com/CloudNationHQ/az-cn-go-validor/issues/72)) ([7708222](https://github.com/CloudNationHQ/az-cn-go-validor/commit/7708222bc67e8ab32f51beff694cf714678edc0c))

## [1.13.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.7...v1.13.0) (2025-11-10)


### Features

* add configurable examples path with optional pattern support ([#70](https://github.com/CloudNationHQ/az-cn-go-validor/issues/70)) ([2788db3](https://github.com/CloudNationHQ/az-cn-go-validor/commit/2788db304598ed6f98c295559ce13978b2e28bd1))
* **deps:** bump github.com/ulikunitz/xz from 0.5.10 to 0.5.14 ([#15](https://github.com/CloudNationHQ/az-cn-go-validor/issues/15)) ([47c39e6](https://github.com/CloudNationHQ/az-cn-go-validor/commit/47c39e607172c0527d428457eb7f7ac16017c452))

## [1.12.7](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.6...v1.12.7) (2025-09-22)


### Bug Fixes

* revert to v1.11.0 ([#63](https://github.com/CloudNationHQ/az-cn-go-validor/issues/63)) ([cb9f2a2](https://github.com/CloudNationHQ/az-cn-go-validor/commit/cb9f2a22533ec82bada05b9b8f483450e665a251))

## [1.11.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.10.1...v1.11.0) (2025-09-16)


### Features

* use hclwrite for terraform module source rewrites ([#46](https://github.com/CloudNationHQ/az-cn-go-validor/issues/46)) ([fee8830](https://github.com/CloudNationHQ/az-cn-go-validor/commit/fee8830c3499fc93bc74c43eaee23015c0b7dfe3))

## [1.10.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.10.0...v1.10.1) (2025-09-16)


### Bug Fixes

* fix submodule detection ([#44](https://github.com/CloudNationHQ/az-cn-go-validor/issues/44)) ([91f96bf](https://github.com/CloudNationHQ/az-cn-go-validor/commit/91f96bf81a096933054018d39be9f590d6c793c4))

## [1.10.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.9.0...v1.10.0) (2025-09-16)


### Features

* add submodule support when switching to local paths ([#42](https://github.com/CloudNationHQ/az-cn-go-validor/issues/42)) ([3063941](https://github.com/CloudNationHQ/az-cn-go-validor/commit/3063941951b45aac577c55b4acc6d4b7e2ee4439))

## [1.9.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.8.1...v1.9.0) (2025-09-12)


### Features

* consolidate test logic, modernize codebase, added configurable namespaces and improved error handling ([#40](https://github.com/CloudNationHQ/az-cn-go-validor/issues/40)) ([31b173c](https://github.com/CloudNationHQ/az-cn-go-validor/commit/31b173c8d1e7a9098572837391f63ce22ca4a96b))

## [1.8.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.8.0...v1.8.1) (2025-09-12)


### Bug Fixes

* improve cross-platform repository name extraction using git remote ([#38](https://github.com/CloudNationHQ/az-cn-go-validor/issues/38)) ([ca2fb26](https://github.com/CloudNationHQ/az-cn-go-validor/commit/ca2fb262c7849bed614762483f090835adf26f1c))

## [1.8.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.7.0...v1.8.0) (2025-09-09)


### Features

* refactor test functions with functional options and eliminate duplicate code ([#36](https://github.com/CloudNationHQ/az-cn-go-validor/issues/36)) ([34998a0](https://github.com/CloudNationHQ/az-cn-go-validor/commit/34998a04f3474e154e51bd9d121d344f77ee2af9))

## [1.7.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.6.0...v1.7.0) (2025-08-29)


### Features

* implement interfaces and context aware patterns ([#34](https://github.com/CloudNationHQ/az-cn-go-validor/issues/34)) ([c2a2165](https://github.com/CloudNationHQ/az-cn-go-validor/commit/c2a21658263f39ca7e3f18da2664c8b4d3fc0832))

## [1.6.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.5.2...v1.6.0) (2025-08-29)


### Features

* remove the local plan functionality ([#32](https://github.com/CloudNationHQ/az-cn-go-validor/issues/32)) ([ca1275c](https://github.com/CloudNationHQ/az-cn-go-validor/commit/ca1275c15dec9b2bbf30ffa8925ada935d5b5ad5))

## [1.5.2](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.5.1...v1.5.2) (2025-08-29)


### Bug Fixes

* remove opentofu defaults using terratest ([#30](https://github.com/CloudNationHQ/az-cn-go-validor/issues/30)) ([2dac57e](https://github.com/CloudNationHQ/az-cn-go-validor/commit/2dac57e6a02d360f8a009c90d4ce052b337592f9))

## [1.5.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.5.0...v1.5.1) (2025-08-29)


### Bug Fixes

* fix plan output validation ([#28](https://github.com/CloudNationHQ/az-cn-go-validor/issues/28)) ([aa26be0](https://github.com/CloudNationHQ/az-cn-go-validor/commit/aa26be060e0202e310320a99490a098c722cd524))

## [1.5.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.4.0...v1.5.0) (2025-08-29)


### Features

* add validation for local modules plan only ([#26](https://github.com/CloudNationHQ/az-cn-go-validor/issues/26)) ([73aa284](https://github.com/CloudNationHQ/az-cn-go-validor/commit/73aa2848fad6f219534bea6952fc78ddf855424d))

## [1.4.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.3.0...v1.4.0) (2025-08-29)


### Features

* add issue and PR templates ([#23](https://github.com/CloudNationHQ/az-cn-go-validor/issues/23)) ([900ff55](https://github.com/CloudNationHQ/az-cn-go-validor/commit/900ff55bd69ef58202768924b3180d4740c1986c))
* provider block is dynamic as well now ([#25](https://github.com/CloudNationHQ/az-cn-go-validor/issues/25)) ([1ae5dc8](https://github.com/CloudNationHQ/az-cn-go-validor/commit/1ae5dc8d73af91d91d686b34c42290dfd1906b87))

## [1.3.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.2.1...v1.3.0) (2025-08-29)


### Features

* instead of using file reverts fetch the latest module version through registry api calls ([#21](https://github.com/CloudNationHQ/az-cn-go-validor/issues/21)) ([6e6addc](https://github.com/CloudNationHQ/az-cn-go-validor/commit/6e6addc43d388b3ffca19f51f3864d87333a2ec8))

## [1.2.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.2.0...v1.2.1) (2025-08-29)


### Bug Fixes

* fix path issues ([#19](https://github.com/CloudNationHQ/az-cn-go-validor/issues/19)) ([737289e](https://github.com/CloudNationHQ/az-cn-go-validor/commit/737289ec40156f90e50a671f1c657eec0c3d252d))

## [1.2.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.1.2...v1.2.0) (2025-08-29)


### Features

* add sequential and local test functions ([#16](https://github.com/CloudNationHQ/az-cn-go-validor/issues/16)) ([e370dd4](https://github.com/CloudNationHQ/az-cn-go-validor/commit/e370dd47f4a38ca3344af88efaa932927141057e))

## [1.1.2](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.1.1...v1.1.2) (2025-08-14)


### Bug Fixes

* fix flags issues ([#13](https://github.com/CloudNationHQ/az-cn-go-validor/issues/13)) ([607dce5](https://github.com/CloudNationHQ/az-cn-go-validor/commit/607dce5d7356fa9d09437c07905fc7c7999e410a))

## [1.1.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.1.0...v1.1.1) (2025-07-19)


### Bug Fixes

* fix context issues ([#11](https://github.com/CloudNationHQ/az-cn-go-validor/issues/11)) ([a2600c3](https://github.com/CloudNationHQ/az-cn-go-validor/commit/a2600c3bcdc4334268bf472a692e00c994e0bf8e))

## [1.1.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.0.1...v1.1.0) (2025-07-19)


### Features

* **deps:** bump github.com/gruntwork-io/terratest from 0.48.2 to 0.49.0 ([#7](https://github.com/CloudNationHQ/az-cn-go-validor/issues/7)) ([548dbe7](https://github.com/CloudNationHQ/az-cn-go-validor/commit/548dbe70ebf40bcf0b443c0004863dba56188754))
* **deps:** bump github.com/gruntwork-io/terratest from 0.49.0 to 0.50.0 ([#9](https://github.com/CloudNationHQ/az-cn-go-validor/issues/9)) ([4944c24](https://github.com/CloudNationHQ/az-cn-go-validor/commit/4944c2463e4ab1dab902f04d4d914b6f8544e26a))
* modernize codebase with config struct, context support, and improved error handling ([#10](https://github.com/CloudNationHQ/az-cn-go-validor/issues/10)) ([ce73c58](https://github.com/CloudNationHQ/az-cn-go-validor/commit/ce73c5841166b925102ce66d02027ed2a6df17e9))

## [1.0.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.0.0...v1.0.1) (2025-05-08)


### Bug Fixes

* use fieldfunc for better performance and fix linting issues ([#5](https://github.com/CloudNationHQ/az-cn-go-validor/issues/5)) ([2e60275](https://github.com/CloudNationHQ/az-cn-go-validor/commit/2e60275b89a2b386865e0641a0a48f5cf5eb26ab))

## 1.0.0 (2025-05-06)


### Features

* add initial files ([#1](https://github.com/CloudNationHQ/az-cn-go-validor/issues/1)) ([6de3601](https://github.com/CloudNationHQ/az-cn-go-validor/commit/6de36010c39593152f2f1ba7a294656e38861d8e))
* **deps:** bump golang.org/x/net from 0.34.0 to 0.38.0 ([#3](https://github.com/CloudNationHQ/az-cn-go-validor/issues/3)) ([dbed2d0](https://github.com/CloudNationHQ/az-cn-go-validor/commit/dbed2d085a5e1bc9e3c7726ed7327bf43ef33ae0))
