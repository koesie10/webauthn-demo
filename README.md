# webauthn-demo : Web Authentication Demo in Go

## Overview [![Build Status](https://travis-ci.org/koesie10/webauthn.svg?branch=master)](https://travis-ci.org/koesie10/webauthn-demo)

This project is a demo of the [webauthn](https://github.com/koesie10/webauthn) library.

## Install

```
go get github.com/koesie10/webauthn-demo
```

## Running

```
go run .
```

This will start the server at port 9000, reachable at [localhost:9000](http://localhost:9000).

The Web Authentication API is only usable in Secure Contexts, i.e. HTTPS. Fortunately, localhost has also been defined
as a secure context, so you will be able to test the Web Authentication API on localhost.

## License

MIT.
