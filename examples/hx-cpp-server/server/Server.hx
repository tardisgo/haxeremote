package ;

@:cppFileCode('
  extern "C" char *hxrIn; 
  extern "C" char *hxrOut; 
  extern "C" int hxrCalling; 
  ')


class Server {
  static var ctx:haxe.remoting.Context;

  function new() { }

  function foo(x,y) { return x + y; }
  function bar(x:Float,y:Float):Float { return x + y; }

  public static function main() {
    ctx = new haxe.remoting.Context();
    ctx.addObject("Server",new Server());
    trace("SETUP foo, bar, hxcppRemote",
      new Server().foo(1,2),
      new Server().bar(3,4),
      Server.hxcppRem(null));
    while(true) { 
      var msg:Dynamic = cpp.vm.Thread.readMessage(false); // just a way to do runtime.Gosched()
      if ( msg != null ) {
        throw "cpp.vm.Thread.readMessage returned unexpected value: "+Std.string(msg);
      }     
      if( untyped __cpp__("hxrCalling") != 0 ){
       var r = hxcppRem( untyped __cpp__("hxrIn"));
       untyped __cpp__("hxrOut=r");
       untyped __cpp__("hxrCalling=0");
      }
    }   
  }

  static function hxcppRem(l:cpp.ConstPointer<cpp.Char>):cpp.ConstPointer<cpp.Char>{
      //trace("hxcppRem0");     
      if(l==null) return null;   
      try { 
        //trace("hxcppRem1",l);     
        var str = cpp.NativeString.fromPointer(l);
        //trace("hxcppRem2",str,ctx);     
        var rstr:String=haxe.remoting.HttpConnection.processRequest(str,ctx);
        //trace("hxcppRem3",rstr);     
        var ret:cpp.ConstPointer<cpp.Char> = cpp.NativeString.c_str(rstr);
        //trace("hxcppRem4",ret); 
        return ret;
      } catch (err:Dynamic) {
        trace("Catch err: "+Std.string(err));
        return null;
      }    
  }
}