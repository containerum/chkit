# chkit
[![Build Status](https://travis-ci.org/containerum/chkit.svg?branch=master)](https://travis-ci.org/containerum/chkit)

Chkit is a desktop CLI client for [Containerum](https://github.com/containerum/containerum).

## Prerequisites
* Containerum
* Windows/Linux/MacOS
* Golang >= 1.8

## Installation
To instal chkit run:

```bash
go get -u -v github.com/containerum/chkit
```
*or*

* in your GOPATH/src run

```bash
git clone https://github.com/containerum/chkit.git
```

*or*

* just launch from **[binaries](https://github.com/containerum/chkit/releases)** 


### To use chkit with your own Kubernetes cluster:

In chkit run
```bash
chkit set api YOUR_API_ADDRESS
```

*or* 

Buid with env:
```bash
export CONTAINERUM_API="YOUR_API_ADDRESS"
make release
```

### To use chkit with Containerum Online
By default chkit connects to [Containerum Online](https://containerum.com/price/online/) platform. 

## Docs
To learn more about **chkit** commands, please refer to the [Docs section](https://docs.containerum.com/docs/about/) on our website.

## Contributions
Please submit all contributions concerning chkit component to this repository. If you want to make contributions to the project, please start by checking the source code at chkit/pkg/cli/.

## License
Chkit project is licensed under the terms of the MIT license. Please see LICENSE in this repository for more details. 
