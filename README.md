hrd [![Build Status](https://secure.travis-ci.org/101loops/hrd.png)](https://travis-ci.org/101loops/hrd)
===

This Go package extends the standard package [appengine.datastore](http://godoc.org/code.google.com/p/appengine-go/appengine/datastore) with useful features:
- caching in local memory as well as memcache
- DSL for read, write and delete
- lifecycle hooks (e.g. beforeSave)
- logging of all datastore actions

The library is used in production and actively worked on. So expect things to change.

**This is still alpha quality. It may have one or two bugs and memory leaks.**

Pull requests are very welcome :)


### Installation
`go get github.com/101loops/hrd`

### Documentation

The [Documentation](http://godoc.org/github.com/101loops/hrd) is still quite sparse.
Please have a look at the E2E tests to see how it is used.

### Credit
- Google: [https://code.google.com/p/appengine-go/]
- OpenVN: [https://github.com/openvn/datastone]
- Jeff Huter: [https://bitbucket.org/SlothNinja/gaelic]
- Matt Jibson: [https://github.com/mjibson/goon]

Without those projects this library would not exist. Thanks!

### License
Apache License 2.0 (see LICENSE).

### Usage

TODO
