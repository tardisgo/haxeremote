# haxeremote
Golang implementation of Haxe remoting server 

Work in progress, not a full implementation.

The only Haxe types that can be (Un-)Serialized to/from Go are:

| Haxe | Go |
| ---- | -- |
| Array<Dynamic> | []interface{} |
| Bool | bool |
| Float | float64 |
| Int  | int |
| Null | interface{}{nil} |
| String | string |
| 0 | int(0) |
