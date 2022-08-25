cors-anywhere in Go
===================

A cors-anywhere proxy written in Go which can be used from CLI or as a
library and mounted in your servers directly.

Inspired by https://github.com/Redocly/cors-anywhere

## Usage

**Install:**

```shell
go get -u github.com/goware/corsanywhere/cmd/corsanywhere
```

**Terminal:**

```shell
corsanywhere -port=8080
```

You can now send CORS-enabled requests to http://localhost:8080/<URL>.

For example, say you have a CORS-enabled API at https://api.mydomain.com/accounts,
you can send the same request, but put `http://localhost:8080/` in front of the url,
so you'd make a request to: `http://localhost:8080/https://api.mydomain.com/accounts`

## LICENSE

MIT
