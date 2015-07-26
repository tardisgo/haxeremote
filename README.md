# haxeremote
Go implementation of the Haxe http remoting server.

Work in progress, not a full implementation. Only developed thus far to enable the tgoremote package (which only currently uses String/string values) and a little local testing. 

The only Haxe types that can be (Un-)Serialized to/from Go are:

| Haxe           | Go               |
| -------------- | ---------------- |
| Array<Dynamic> | []interface{}    |
| Bool           | bool             |
| Float          | float64          |
| Int            | int              |
| Null           | interface{}(nil) |
| String         | string           |
| 0              | int(0)           |
| haxe.io.Bytes  | []byte           |

See the package preamble (and code) in serialization.go for the the full current state of implementation.

To use SSL, each Haxe target and use case has different requirements, so SSL is not used in these examples.  

The examples folder contains both a straightforward go-http-server; and a rather more complex and experimental hx-cpp-server (where the Go code calls a linked Haxe/C++ code library utilizing the Haxe remote protocol, using shared variables as the transport rather than http).
