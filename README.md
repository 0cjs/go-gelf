go-gelf - GELF Library and Writer for Go
========================================

[GELF] (Graylog Extended Log Format) is an application-level logging
protocol that avoids many of the shortcomings of [syslog]. This
library formats and sends GELF messages.


go-gelf Versions
----------------

The different versions of go-gelf are distinguished using [semantic
versioning] and offered via [gopkg.in]. Thus, the import (and `go
get`) path for the `gelf` package at the latest 1.x.x version (API
version 1) is:

    gopkg.in/Graylog2/go-gelf.v1/gelf

For temporary backwards compatibility, the `github.com/Graylog2/go-gelf`
path currently serves the `v1` branch. However, this will be changed
to the `master` branch in the future, so please ensure any clients
change to the the correct versioned path described above or they will
break.

There are currently two major versions of the API:

* **[v1]** is the original version, offering UDP connectivity only,
  and is initialized with the `NewWriter()` function.
* **[v2]** adds TCP support and changes the constructors to
  `NewTCPWriter()` and `NewUDPWriter()`.


Overview
--------

[GELF] can be run over any stream or datagram transport protocol. When
used over a datagram protocol it has support for compression to
shorten messages and [chunking] to allow long messages to be split
over multiple datagrams.

This implementation supports UDP and TCP as a transport protocol. TLS
and other streams (e.g., to a file) are currently unsupported, but are
planned for a new version of the API.

The library provides an API that applications can use to log messages
directly to a Graylog server and an `io.Writer` that can be used to
redirect the standard library's log messages (`os.Stdout`) to a
Graylog server.


Usage
-----

The easiest way to integrate graylog logging into your go app is by
having your `main` function (or even `init`) call `log.SetOutput()`.
By using an `io.MultiWriter`, we can log to both stdout and graylog -
giving us both centralized and local logs.  (Redundancy is nice).

	package main

	import (
		"flag"
		"gopkg.in/Graylog2/go-gelf.v1/gelf"
		"io"
		"log"
		"os"
	)

	func main() {
		var graylogAddr string

		flag.StringVar(&graylogAddr, "graylog", "", "graylog server addr")
		flag.Parse()

		if graylogAddr != "" {
            // If using UDP
			gelfWriter, err := gelf.NewUDPWriter(graylogAddr)
            // If using TCP
            //gelfWriter, err := gelf.NewTCPWriter(graylogAddr)
			if err != nil {
				log.Fatalf("gelf.NewWriter: %s", err)
			}
			// log to both stderr and graylog2
			log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
			log.Printf("logging to stderr & graylog2@'%s'", graylogAddr)
		}

		// From here on out, any calls to log.Print* functions
		// will appear on stdout, and be sent over UDP or TCP to the
		// specified Graylog2 server.

		log.Printf("Hello gray World")

		// ...
	}

The above program can be invoked as:

	go run test.go -graylog=localhost:12201

When using UDP messages may be dropped or re-ordered. However, Graylog
server availability will not impact application performance; there is
a small, fixed overhead per log call regardless of whether the target
server is reachable or not.


To Do
-----

- WriteMessage example

License
-------

go-gelf is offered under the MIT license, see LICENSE for details.



[GELF]: http://docs.graylog.org/en/2.2/pages/gelf.html
[chunking]: http://docs.graylog.org/en/2.2/pages/gelf.html#chunked-gelf
[gopkg.in]: https://labix.org/gopkg.in
[semantic versioning]: http://semver.org
[syslog]: https://tools.ietf.org/html/rfc5424
[v1]: https://github.com/Graylog2/go-gelf/tree/v1
[v2]: https://github.com/Graylog2/go-gelf/tree/v2
