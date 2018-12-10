# sumoll
Small sumologic payload sender client

[![CircleCI](https://circleci.com/gh/ushios/sumoll.svg?style=shield&circle-token=1fd034d3a678ace2e03a3f496736c9e8b0093935)](https://circleci.com/gh/ushios/sumoll)
[![Build Status](https://travis-ci.org/ushios/sumoll.svg?branch=master)](https://travis-ci.org/ushios/sumoll)
[![Coverage Status](https://coveralls.io/repos/github/ushios/sumoll/badge.svg?branch=feature%2Fci)](https://coveralls.io/github/ushios/sumoll?branch=feature%2Fci)
[![Go Report Card](https://goreportcard.com/badge/github.com/ushios/sumoll)](https://goreportcard.com/report/github.com/ushios/sumoll)
[![GoDoc](https://godoc.org/github.com/ushios/sumoll?status.svg)](https://godoc.org/github.com/ushios/sumoll)

# Installation

```console
$ go get github.com/ushios/sumoll
```

# Usage

## Send string to http source collection

```go
client := NewHTTPSourceClient("http://collectors.au.sumologic.com/receiver/v1/http/...")
client.Send(strings.NewReader("your message here."))
```

# Developing

## Integration test

There's a test actually injecting data to sumologic. In order that to run, you need to set the env variable 
`SUMOLL_TEST_HTTP_SOURCE_URL` like this:

```bash
SUMOLL_TEST_HTTP_SOURCE_URL="https://collectors.eu.sumologic.com/receiver/v1/http/randomCollectorURL" go test --count=1 . -v
```
