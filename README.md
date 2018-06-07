# chkit
[![Build Status](https://travis-ci.org/containerum/chkit.svg?branch=master)](https://travis-ci.org/containerum/chkit)

Chkit is desktop CLI client for [Containerum](https://github.com/containerum/containerum).

## Prerequisites
* Containerum
* Windows/Linux/MacOS

## Installation for self-hosted Containerum
To use chkit with your own Kubernetes cluster:

Buid with env:
```bash
export CONTAINERUM_API="YOUR_API_ADDRESS"
make release
```

## Installation for Containerum Online
By default chkit connects to [Containerum Online](https://containerum.com/price/online/) platform. 

To install chkit from Docker run

```bash
export CONTAINERUM_API="https://api.containerum.io:8082"
make docker 
docker run containerum/chkit
```

*or*

* Launch from **[binaries](https://github.com/containerum/chkit/releases)**

## Docs
To learn more about **chkit** commands, please refer to the [Docs section](https://docs.containerum.com/docs/about/) on our website.

## Contributions
Please submit all contributions concerning chkit component to this repository. If you want to make contributions to the project, please start by checking the source code at chkit/pkg/cli/.

## License
Chkit project is licensed under the terms of the MIT license. Please see LICENSE in this repository for more details. 
