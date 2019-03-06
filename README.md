# Elastos.ELA.Monitor

## Introduction
Elastos.ELA.Monitor is an current dpos information monitor tools and we will add more feature in the future.

## Introduction

## Build the monitor

#### 1. Setup basic workspace
In this instruction we use ~/dev/src/github.com/elastos as our working directory. If you clone the source code to a different directory, please make sure you change other environment variables accordingly (not recommended). 

```shell
$ mkdir -p ~/dev/bin
$ mkdir -p ~/dev/src/github.com/elastos/
```

#### 2. Set correct environment variables

```shell
export GOROOT=/usr/local/opt/go@1.9/libexec
export GOPATH=$HOME/dev
export GOBIN=$GOPATH/bin
export PATH=$GOROOT/bin:$PATH
export PATH=$GOBIN:$PATH
```

#### 3. Check Go version and glide version

Check the golang and glider version. Make sure they are the following version number or above.

```shell
$ go version
go version go1.9.2 darwin/amd64

$ glide --version
glide version 0.13.1
```

If you cannot see the version number, there must be something wrong when install.

#### 4. Clone source code to $GOPATH/src/github/elastos folder
Make sure you are in the folder of $GOPATH/src/github.com/elastos
```shell
$ git clone https://github.com/yiyanwannian/Elastos.ELA.Monitor.git
```

If clone works successfully, you should see folder structure like $GOPATH/src/github.com/elastos/Elastos.ELA/Makefile

#### 5. Install dependencies using Glide

```shell
$ cd $GOPATH/src/github.com/elastos/Elastos.ELA.Monitor
$ glide update && glide install
``` 

#### 6. Make

Build the monitor.
```shell
$ cd $GOPATH/src/github.com/elastos/Elastos.ELA.Monitor
$ make
```

If you did not see any error message, congratulations, you have made the ELA full node.

#### 7. Run the monitor on Mac or Ubuntu

Run the node.
```shell
$ ./monitor
```
