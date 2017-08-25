go-apriltag
============

[![GoDoc Reference](https://godoc.org/github.com/twitchyliquid64/go-apriltag?status.svg)](http://godoc.org/github.com/twitchyliquid64/go-apriltag)
[![Coverage Status](https://coveralls.io/repos/twitchyliquid64/go-apriltag/badge.svg?branch=master)](https://coveralls.io/r/twitchyliquid64/go-apriltag?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/twitchyliquid64/go-apriltag)](https://goreportcard.com/report/github.com/twitchyliquid64/go-apriltag)

Description
------------

Apriltags image recognition for Go. Uses Cgo, but does not require any libraries or dependencies.

Installation
------------

This package can be installed with the go get command:

    go get github.com/twitchyliquid64/go-apriltag

**NOTE**: Make sure cgo works on your system.

Only tested on Linux, but should work (or need trivial modifications) to work on other platforms.

Documentation
-------------

API documentation can be found here: http://godoc.org/github.com/twitchyliquid64/go-apriltag

Trivial example, finding apriltags in a PDF:

```go
detector := apriltag.New()
defer detector.Close()
f, err := os.Open("testtags.png")
if err != nil {
  t.Fatal(err)
}
defer f.Close()

img, err := png.Decode(f)
if err != nil {
  t.Fatal(err)
}

findings := detector.Find(ImgToGrayscale(img)) // list of discovered tags
```

License
----------
In this repository, those files are an amalgamation of code that was copied from Apriltag.
The license of that code is the same as the license of Apriltag.
Apriltag copyright notices remain intact as per license requirements.
