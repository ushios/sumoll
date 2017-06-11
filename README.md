# sumoll
Small sumologic payload sender client

[![CircleCI](https://circleci.com/gh/ushios/sumoll.svg?style=shield&circle-token=1fd034d3a678ace2e03a3f496736c9e8b0093935)](https://circleci.com/gh/ushios/sumoll)

# Installation

```console
$ go get github.com/ushios/sumoll
```

# Usage

### Send string to http source collection

```go
client := NewHTTPSourceClient("http://collectors.au.sumologic.com/receiver/v1/http/...")
client.Send(strings.NewReader("your message here."))
```
