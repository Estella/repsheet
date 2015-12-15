# repsheet [![Build Status](https://secure.travis-ci.org/repsheet/repsheet.png)](http://travis-ci.org/repsheet/repsheet?branch=master)

The Repsheet command line utility. If you are looking for the original Repsheet project, please visit the [NGINX](https://github.com/repsheet/repsheet-nginx) or [Apache](https://github.com/repsheet/repsheet-apache) modules. They contain implementations for each webserver, but share a common [core library](https://github.com/repsheet/librepsheet).

## Installation

```
$ go get github.com/repsheet/repsheet
```

## Usage

```
$ repsheet -blacklist=1.1.1.1 -reason=cli
$ repsheet -whitelist=2.2.2.2 -reason=cli
$ repsheet -mark=3.3.3.3 -reason=cli
$ repsheet -remove=4.4.4.4 
$ repsheet -list
Whitelisted Actors
  2.2.2.2:repsheet:ip:whitelisted
Blacklisted Actors
  1.1.1.1:repsheet:ip:blacklisted
Marked Actors
  3.3.3.3:repsheet:ip:marked
$ repsheet -status=1.1.1.1
1.1.1.1 is blacklisted. Reason: cli
$ repsheet -status=8.8.8.8
8.8.8.8 is OK
```
