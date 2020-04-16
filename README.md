# roadmapper

[![Build Status](https://travis-ci.com/peteraba/roadmapper.svg?branch=master)](https://travis-ci.com/peteraba/roadmapper)
[![Docs](https://img.shields.io/badge/docs-current-brightgreen.svg)](https://docs.rdmp.app)
[![codecov](https://codecov.io/gh/peteraba/roadmapper/branch/master/graph/badge.svg)](https://codecov.io/gh/peteraba/roadmapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/peteraba/roadmapper)](https://goreportcard.com/report/github.com/peteraba/roadmapper)
[![GitHub release](https://img.shields.io/github/tag/peteraba/roadmapper.svg?label=release)](https://github.com/peteraba/roadmapper/releases)
[![](https://images.microbadger.com/badges/version/peteraba/roadmapper.svg)](https://microbadger.com/images/peteraba/roadmapper)
[![](https://images.microbadger.com/badges/image/peteraba/roadmapper.svg)](https://microbadger.com/images/peteraba/roadmapper)
[![License](https://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/peteraba/roadmapper/blob/master/LICENSE)

Roadmapper is a CLI tool and webservice designed to help maintaining and tracking roadmaps. Learn more in the [official documentation](https://docs.rdmp.app/).

### Thanks to

- [cli](https://github.com/urfave/cli) - A simple, fast, and fun package for building command line apps in Go
- [echo](https://echo.labstack.com/) - High performance, extensible, minimalist Go web framework
- [canvas](https://github.com/tdewolff/canvas) - Cairo in Go: vector to SVG, PDF, EPS, raster, HTML Canvas, etc.
- [pq](https://github.com/lib/pq) - A pure Go postgres driver for Go's database/sql package
- [uitable](https://github.com/gosuri/uitable) - A go library for representing data as tables for terminal applications.
- [Testify](https://github.com/stretchr/testify) - A set of packages that provide many tools for testifying that your code will behave as you intend
- [dockertest](https://github.com/ory/dockertest) - Integration tests against third party services
- [chromedp](https://github.com/chromedp/chromedp) - A faster, simpler way to drive browsers supporting the Chrome DevTools Protocol.
- [bindata](https://github.com/kevinburke/go-bindata) - A small utility which generates Go code from any file. Useful for embedding binary data in a Go program.

- [x] links added via JS
- [x] editor bug: tabs always select whole lines
- [x] editor bug: indentation errors are not displayed
- [x] editor bug: spaces are not replaced with tabs sometimes
- [x] editor bug: pastes cause the cursor to go to the end of the textarea
- [ ] editor feature: implement / restore edit history
- [ ] backend validation
- [ ] fix displaying multiline project titles
- [ ] spam filtering
- [ ] remove color support
- [ ] stripe overlay for project visualization
