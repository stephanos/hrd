hrd [![Build Status](https://secure.travis-ci.org/101loops/hrd.png)](https://travis-ci.org/101loops/hrd)
===

This Go package extends the standard package *appengine.datastore* with useful features:
- lifecycle hooks (e.g. beforeSave)
- fields NOT indexed by default
- omitempty: does not save empty/zero values (and thereby does not index them)
- caching of results in local memory as well as memcache

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

TODO
