# repsheet

The Repsheet command line utility. If you are looking for the original Repsheet project, please visit the [NGINX](https://github.com/repsheet/repsheet-nginx) or [Apache](https://github.com/repsheet/repsheet-apache) modules. They contain implementations for each webserver, but share a common [core library](https://github.com/repsheet/librepsheet).

## Installation

```
$ go get github.com/repsheet/repsheet
```

## Usage

```
$ repsheet -blacklist=1.1.1.1
Blacklisting 1.1.1.1
$ repsheet -whitelist=2.2.2.2
Whitelisting 2.2.2.2
$ repsheet -mark=3.3.3.3
Marking 3.3.3.3
$ repsheet -list
Whitelisted Actors
  2.2.2.2:repsheet:ip:whitelisted
Blacklisted Actors
  1.1.1.1:repsheet:ip:blacklisted
Marked Actors
  3.3.3.3:repsheet:ip:marked
```
