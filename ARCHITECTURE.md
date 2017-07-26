`go-gelf` API and Architecture Change Proposal
============================================

**NOTE:** This is a rough draft. Please try to grasp the overall idea
rather than getting too hung up on details I may have glossed over or
gotten wrong. I do, however, welcome feedback on the details if you
agree with the general gist of this idea.

I am monitoring the [go-gelf Gitter channel], so I welcome interactive
discussions there.

[@iporsut]'s comments on [MiLk's review] of go-gelf [PR#9] noted that
the usual thing in the Go standard libraries is to have a `Dial()`
function. Inspired by this, I ([cjs](https://github.com/0cjs)) had a
look through packages such as `net`, `crypto/tls` and `net/textproto`,
and now it seems to me as if we've been going at it rather wrong here.

Oddly enough, the system seems to divide up into layers almost exactly
corresponding to the [OSI model], which is a first in my experience.
So I'll use those names here and in the descriptions mention what they
correspond to in other views of layering.

Here's how it looks to me as if the system divides up.


### Application Layer

This builds GELF messages, i.e., the JSON objects and sends them to a
stream or packet connection via the `io.Writer` interface. This
provides:

    type GELFStream interface {
        Write([]byte) (int, error)
        WriteMessage(*Message) (error)
    }

	type Message struct {
		Host         string                 `json:"host"`
		ShortMessage string                 `json:"short_message"`
		FullMessage  string                 `json:"full_message,omitempty"`
		TimeUnix     float64                `json:"timestamp"`
		Level        int32                  `json:"level,omitempty"`
		Extra        map[string]interface{} `json:"-"`
	}


The fields in `Message` are based on the [GELF Payload Specification].
The following fields are removed from the current version of the
`Message` struct:

* `Version`: `version` in the GELF message is set by the library; the
  library client should not be able to supply it. If a new version
  comes along, we need to modify the library to set it properly.
* `Facility`: `facility` in GELF messages is deprecated; clients
  should send the `_facility` additional field (by setting a
  `facility` key in `Message.Extra`) instead.
* `RawExtra`: This is an internal field used by the library and
  should not be exposed to the library's clients.

We probably ought to allow certain fields to be unspecified and filled
in with a reasonable default value supplied by the library, as
follows:

* `Host`: if not supplied, the name of the host on which the code is
  running.
* `FullMessage`: left out of the GELF message if not supplied
* `TimeUnix`: the current time on the host, if not supplied
* `Level`: The valid values are 0 (Emergency) through 7 (Debug)
  according to [RFC 5254][RFC 5254-PRI]. If not supplied the current
  code uses 1 (Alert); 4 (Error) or 5 (Notice) might be a more
  reasonable default. The Unix `logger` program uses 5 as the default.

This looks to me as if it would require changes to the types of or
interpretations of the fields.


### Presentation Layer

Here we handle how the GELF messages are packaged to be transported
across a stream or packet protocol. (The latter is often refered to as
a datagram protocol, but the Go standard libraries use "packet" so I
use that here.) The ways of doing this can also be thought of as
"CODECs," though in this library we don't decode, only encode.

There are two basic encodings, with variations.

#### Packet

Packet encoding is used for UDP and Unix datagram sockets. Writes to
these sockets use the standard go `io.Writer` interface
(`Write([]bytes) (int, error)`) but each call to `Write` generates a
single packet that either succeeds or fails in its entirety.

Here is a summary of the details described in [GELF via UDP]:

First, the messages may optionally be compressed with GZIP or ZLIB.
(GZIP is the protocol default.)

Then, to handle GELF messages larger than can fit into a packet the
messages can be optionally be divided into chunks. Each chunk is
written as a separate packet with the following header:

* Magic identifier: two bytes `0x1e 0x0f`. Message ID: an arbitrary
* eight byte value that's the same for every
  chunk in a message.
* Sequence number: 1 byte starting at `0` and always less than the
  sequence count.
* Sequence count: 1 byte valued `1` through `128` indicating the total
  number of chunks.

The current library sets a fixed maximum `ChunkSize` of 1420 bytes,
presumably chosen in the hope of avoiding UDP fragmentation (though
one would have to use a chunk size of less than 500 bytes to guarantee
this). There doesn't seem to be any reasonable way to get the optimum
chunk size from a `Writer` but it would probably be good to be able to
pass in an option for this, particular if the `Writer` is writing a
Unix domain datagram socket that can handle much larger sizes.
(`/proc/sys/net/core/wmem_*` on my CentOS 7 system indicate that the
default and maximum sizes are 212992 bytes.)

#### Stream

Stream encoding is used for TCP, TLS, stream-oriented Unix sockets and
similar; in all cases multiple messages are transfered as straight
sequences of bytes. No compression or chunking is allowed, but a
message delimiter must be used. [GELF via TCP] says that the delimiter
must be a "null byte" (`\0`), but Graylog 2.2.3 allows the
adminstrator to select either null byte or newline as the delimiter
for a GELF input.

The specification doesn't indicate whether the delimiter separates or
terminates messages, but experimentation with Graylog 2.2.3 indicates
that data between the last delimter and connection close is ignored,
so in practice the delimiter is a terminator.


### Session Layer

The code in [PR#9] does a lot of work to handle reconnecting broken
TCP connections. This can be split out into a separate,
GELF-independent layer that deals with making multiple sequential
stream connections look like a single connection. Use of this should
be selectable by the library client; for some applications it may be
simpler for the client to handle reconnections. For example,
applications reading from the Systemd journal receive a cursor value
with with each message that it can store to indicate the last message
sent; the easiest way to handle network connectivity problems is
simply to terminate when the connection terminates and let the system
restart the program forwarding the logs. This can easily handle
multi-hour or even multi-day connectivity losses without losing
messages.


### Transport Layer

Connections can be provided over a variety of protocols by the Go
standard libraries; these include stream and datagram Unix domain
sockets, UDP, TCP and TLS connections, and even pipes, files or
standard output. Library clients should have the option of setting up
and using any of these with any of the above layers.


Testing
-------

When split into layers like this, unit testing the individual parts is
not difficult for the most part. We would want an integration test or
two that tested connecting all of the layers, however. Setting up a
Graylog server within the test framework would be an extraordinary
amount of work; instead I suggest that we add a specific output format
where bytestrings written by a bottom layer could be compared with
bytestrings of stored data that have been pre-verified by manually
sending them to a Graylog server. These could be stored in files and
committed to the repo; the test framework would compare the output to
these files and the data in the files themselves would be verified by
a small tool that sends the file data to a Graylog server with the
results on the server verified manually.


-----

[go-gelf Gitter channel]: https://gitter.im/Graylog2/go-gelf
[@iporsut]: https://github.com/iporsut
[MiLk's review]: https://github.com/Graylog2/go-gelf/pull/9#pullrequestreview-51093258
[PR#9]: https://github.com/Graylog2/go-gelf/pull/9
[OSI model]: https://en.wikipedia.org/wiki/OSI_model
[GELF Payload Specification]: http://docs.graylog.org/en/2.2/pages/gelf.html#gelf-payload-specification
[RFC 5254]: https://tools.ietf.org/html/rfc5424
[RFC 5254-PRI]: https://tools.ietf.org/html/rfc5424#section-6.2.1
[GELF via UDP]: http://docs.graylog.org/en/2.2/pages/gelf.html#gelf-via-udp
[GELF via TCP]: http://docs.graylog.org/en/2.2/pages/gelf.html#gelf-via-tcp
