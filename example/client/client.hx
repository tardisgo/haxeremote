package;
class Client {
  static function main() {
    var cnx = haxe.remoting.HttpAsyncConnection.urlConnect("http://localhost:8088/haxeremote");
    cnx.setErrorHandler( function(err) { trace('Error: $err'); } );
    cnx.Server.foo.call([1,2], function(data) { trace('Result: $data'); } );
    cnx.Server.foo.call([11,22], function(data) { trace('Result: $data'); } );
    cnx.Server.foo.call([111,222], function(data) { trace('Result: $data'); } );
    cnx.Server.foo.call([1111,2222], function(data) { trace('Result: $data'); } );
    cnx.Server.bar.call(["ding希腊危机bat","doo欧元区dah"], function(data) { trace('Result: $data'); } );
    cnx.Server.fad.call([1111.1111,2222.2222], function(data) { trace('Result: $data'); } );
    cnx.Server.dong.call([1111.1111,2222.2222], function(data) { trace('Result: $data'); } );
    cnx.Server.dingbat.call([haxe.io.Bytes.alloc(42)], 
        function(data) { trace('Result: $data'); });
  }
}
