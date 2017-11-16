![Sprise logo](https://raw.githubusercontent.com/nathan-osman/sprise/master/server/static/img/logo.png)

[![GoDoc](https://godoc.org/github.com/nathan-osman/sprise?status.svg)](https://godoc.org/github.com/nathan-osman/sprise)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

This web application provides a web application for storing and organizing photos.

### Building

Assuming you have GNU Make and Docker installed, you can build the application by running:

    make

This will create an executable named `dist/sprise` which can then be run to launch the application. Alternatively, a Docker container can also be built from the binary:

    docker build -t [image name] .
