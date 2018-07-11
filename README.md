# chkit
[![Build Status](https://travis-ci.org/containerum/chkit.svg?branch=master)](https://travis-ci.org/containerum/chkit) [![LOC](https://tokei.rs/b1/github/containerum/chkit)](https://github.com/Aaronepower/tokei) [![GitHub release](https://img.shields.io/github/release/containerum/chkit.svg)](https://github.com/containerum/chkit/releases/latest) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Chkit is a desktop CLI client for [Containerum](https://github.com/containerum/containerum).

## Prerequisites
* Containerum
* Windows/Linux/MacOS
* Golang >= 1.8

## Installation

### To get chkit source:

```bash
go get -u -v github.com/containerum/chkit

```
or

```bash
# in $GOPATH/src/containerum
git clone https://github.com/containerum/chkit.git
```

### To build chkit from sources:
```bash
cd $GOPATH/src/containerum/chkit
make single_release CONTAINERUM_API="https://api.containerum.io"  
# or your Containerum API URL
```
then extract executable from tar.gz archive in ./build  to $GOPATH/bin or another $PATH dir

### To delete chkit:
Just delete executable file from your $GOPATH/bin

## Configuring chkit

### Using chkit with your own Kubernetes cluster
Before using chkit to work with Containerum you have to specify the address of your API. You can do that as follows:

* run
```bash
chkit set api YOUR_API_ADDRESS
```

or 

* build chkit with env:
```bash
export CONTAINERUM_API="YOUR_API_ADDRESS"
make release
```

### Using chkit with Containerum Online
By default chkit connects to [Containerum Online](https://containerum.com/price/online/) platform. 

### Configuraton and logs
|    OS   | Config file path | Logs path |
| ------- | ---------------- | --------- |
| Linux   | ~/.config/containerum/config.toml | ~/.config/containerum/support |
| Windows | /Users/$USERNAME/AppData/Local/containerum/config.toml | /Users/$USERNAME/AppData/Local/containerum/suppport |
| Mac OS | /Library/Application Support/containerum/config.toml | /Library/Logs/containerum

## Docs
To learn more about **chkit** commands, please refer to the [Docs section](https://docs.containerum.com/docs/about/) on our website.

## Contributions
Please submit all contributions concerning chkit component to this repository. If you want to make contributions to the project, please start by checking the source code at chkit/pkg/cli/.

## License
Chkit project is licensed under the terms of the MIT license. Please see LICENSE in this repository for more details. 
