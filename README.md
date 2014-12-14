hrd [![Build Status](https://secure.travis-ci.org/101loops/hrd.png)](https://travis-ci.org/101loops/hrd) [![Coverage Status](https://coveralls.io/repos/101loops/hrd/badge.png?branch=master)](https://coveralls.io/r/101loops/hrd?branch=master)  [![GoDoc](https://camo.githubusercontent.com/6bae67c5189d085c05271a127da5a4bbb1e8eb2c/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f736d61727479737472656574732f676f636f6e7665793f7374617475732e706e67)](http://godoc.org/github.com/101loops/hrd)
===

This Go package extends the package [appengine.datastore](http://godoc.org/code.google.com/p/appengine-go/appengine/datastore)
with very useful additional features.

## Disclaimer

**This package is currently undergoing a massive restructuring. Use at own risk**

## Features
- **caching:** great performance through memcache 
- **fluent API:** concise code for read, query, write and delete actions
- **hybrid query:** queries that have strong consistency and use memcache 
- **lifecycle hooks:** BeforeLoad/AfterLoad and BeforeSave/AfterSave
- **caching control:** turn caching on/off for queries and entities
- **logging:** every datastore action is logged for debugging

Internally it uses [nds](https://github.com/qedus/nds),
 [structor](https://github.com/101loops/structor) and
 [iszero](github.com/101loops/iszero).

## ToDos
- validated projection query
- in-memory cache
- field name & name transformer
- allow to pass-in logger
- RPC listener
- catch more errors at codec creation
- delete from query
- namespace support

## Install
```bash
go get github.com/101loops/hrd
```

## Documentation
[godoc.org](http://godoc.org/github.com/101loops/hrd)

## Credit
- Google: [https://code.google.com/p/appengine-go/]
- OpenVN: [https://github.com/openvn/datastone]
- Jeff Huter: [https://bitbucket.org/SlothNinja/gaelic]
- Matt Jibson: [https://github.com/mjibson/goon]

Without those projects this library would not exist. Thanks!

## License
Apache License 2.0 (see LICENSE).

## Usage
I suggest having a look at the E2E tests to see how it is used.
