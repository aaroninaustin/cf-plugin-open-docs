Overview [![Build Status](https://travis-ci.org/cloudfoundry-community/cf-plugin-open.svg?branch=master)](https://travis-ci.org/cloudfoundry-community/cf-plugin-open)
========

Open app url in browser

Installation
------------

```
$ go get github.com/cloudfoundry-community/cf-plugin-open-docs
$ cf install-plugin $GOPATH/bin/cf-plugin-open-docs
```

Usage
-----

```
$ cf open-docs <appname>
```

Development
-----------

```
cf uninstall-plugin open-docs; go get ./...; cf install-plugin $GOPATH/bin/cf-plugin-open
```
