hrd [![Build Status](https://secure.travis-ci.org/101loops/hrd.png)](https://travis-ci.org/101loops/hrd) [![Coverage Status](https://coveralls.io/repos/101loops/hrd/badge.png?branch=master)](https://coveralls.io/r/101loops/hrd?branch=master)  [![GoDoc](https://camo.githubusercontent.com/6bae67c5189d085c05271a127da5a4bbb1e8eb2c/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f736d61727479737472656574732f676f636f6e7665793f7374617475732e706e67)](http://godoc.org/github.com/101loops/hrd)
===

This Go package extends the standard package [appengine.datastore](http://godoc.org/code.google.com/p/appengine-go/appengine/datastore) with useful features:
- caching in local memory as well as memcache
- DSL for read, write and delete
- lifecycle hooks (e.g. beforeSave)
- logging of all datastore actions

**This is still alpha quality.**

Pull requests are very welcome :)


### Installation
`go get github.com/101loops/hrd`

### Documentation
[godoc.org](http://godoc.org/github.com/101loops/hrd)

### Credit
- Google: [https://code.google.com/p/appengine-go/]
- OpenVN: [https://github.com/openvn/datastone]
- Jeff Huter: [https://bitbucket.org/SlothNinja/gaelic]
- Matt Jibson: [https://github.com/mjibson/goon]

Without those projects this library would not exist. Thanks!

### License
Apache License 2.0 (see LICENSE).

### Usage
I suggest having a look at the E2E tests to see how it is used.