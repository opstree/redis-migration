<p align="left">
  <img src="./img/logo.svg" height="180" width="180">
</p>

[![GitHub Super-Linter](https://github.com/opstree/redis-migration/workflows/CI%20Pipeline/badge.svg)](https://github.com/opstree/redis-migration)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/opstree/redis-migration)
[![Go Report Card](https://goreportcard.com/badge/github.com/opstree/redis-migration)](https://goreportcard.com/report/github.com/opstree/redis-migration)
[![Apache License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/opstree/redis-migration)

Redis migrator is a golang based tool to migrate the database keys from one redis cluster to another. This tool can be used to migrate different types of redis keys from one redis setup to another.

Redis supported keys:-

- String keys
- Hash keys

![](img/architecture.png)

## Quickstart

A quickstart guide for installing, using and managing redis-migrator.

### Installation

redis-migrator installation packages can be found inside [Releases](https://github.com/opstree/redis-migration/releases)

Supported Platforms:-

- Binaries are supported for Linux and Windows platform with these architectures:-
  - Amd 64
  - Arm 64
  - Amd 32
  - Arm 32

For installation on debian and redhat based system, `.deb` and `.rpm` packages can be used.

For installing on MacOS system, use brew:-

```shell
brew install redis-migrator
```
