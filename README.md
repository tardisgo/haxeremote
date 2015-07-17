# haxeremote
Golang implementation of Haxe http remoting server.

Work in progress, not a full implementation. Only developed thus far to enable the tgoremote package, which only currently uses String/string values. 

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

To use SSL, each Haxe target and use case has different requirements, so SSL is not used in these examples.  

A socket-based implementation should eventually be possible.

