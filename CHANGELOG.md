# Changelog

## [1.12.6](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.5...v1.12.6) (2025-09-22)


### Bug Fixes

* fix drift detections ([#61](https://github.com/CloudNationHQ/az-cn-go-validor/issues/61)) ([b3ac829](https://github.com/CloudNationHQ/az-cn-go-validor/commit/b3ac8293d8ea6f943c9727a63c9acba17aaaef0e))

## [1.12.5](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.4...v1.12.5) (2025-09-22)


### Bug Fixes

* small refactor ([#59](https://github.com/CloudNationHQ/az-cn-go-validor/issues/59)) ([2392a78](https://github.com/CloudNationHQ/az-cn-go-validor/commit/2392a78a919fa5132536a62b5766bf4c1a6715db))

## [1.12.4](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.3...v1.12.4) (2025-09-22)


### Bug Fixes

* fix early destroy on wrong path ([#57](https://github.com/CloudNationHQ/az-cn-go-validor/issues/57)) ([ca970e5](https://github.com/CloudNationHQ/az-cn-go-validor/commit/ca970e587ad697f886f6722c870f0c6b0eec2c68))

## [1.12.3](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.2...v1.12.3) (2025-09-22)


### Bug Fixes

* fix order of execution again ([#55](https://github.com/CloudNationHQ/az-cn-go-validor/issues/55)) ([f6c764c](https://github.com/CloudNationHQ/az-cn-go-validor/commit/f6c764c796e2eafe9fc56e7d3e9920cb048edd7c))

## [1.12.2](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.1...v1.12.2) (2025-09-22)


### Bug Fixes

* fix faulty order parity testing ([#52](https://github.com/CloudNationHQ/az-cn-go-validor/issues/52)) ([14851b3](https://github.com/CloudNationHQ/az-cn-go-validor/commit/14851b39cd0f888c5d461cf267528b6382d04105))
* fix logic abit ([#54](https://github.com/CloudNationHQ/az-cn-go-validor/issues/54)) ([a48b219](https://github.com/CloudNationHQ/az-cn-go-validor/commit/a48b219afe15828117fd77c8319a8560661a4d18))

## [1.12.1](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.12.0...v1.12.1) (2025-09-22)


### Bug Fixes

* fix parallel processing ([#50](https://github.com/CloudNationHQ/az-cn-go-validor/issues/50)) ([3a9a79b](https://github.com/CloudNationHQ/az-cn-go-validor/commit/3a9a79b2c4274f90286c9ea4a0439b741c7bcc93))

## [1.12.0](https://github.com/CloudNationHQ/az-cn-go-validor/compare/v1.11.0...v1.12.0) (2025-09-22)


### Features

* implement parity testing between registry and local module sources ([#48](https://github.com/CloudNationHQ/az-cn-go-validor/issues/48)) ([4db1184](https://github.com/CloudNationHQ/az-cn-go-validor/commit/4db118496affbe26f543fc3a3d61317f5aaf9fc6))

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
