# haxeremote
Golang implementation of Haxe remoting server 

Work in progress, not a full implementation.

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

