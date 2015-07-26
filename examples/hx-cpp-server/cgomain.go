package main

/*
// NOTE libraries below are for the test mac, your libraries WILL be different!
// TODO make it easier to use!

#cgo LDFLAGS: -stdlib=libstdc++ server/cpp/libServer.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libstd.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libzlib.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libregexp.a
// /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libsqlite.a /usr/lib/haxe/lib/hxcpp/3,2,102/lib/Mac64/libmysql5.a

#include <stdio.h>
#include <stdlib.h>

// entry point for main haxe
extern int hxmain();

// shared values for haxeremote
char *hxrIn;
char *hxrOut;
int hxrCalling;
*/
import "C"
import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/tardisgo/haxeremote"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		runtime.LockOSThread()
		C.hxmain()
	}()
	r1 := Call([]interface{}{"Server", "foo"}, []interface{}{4, 2})
	fmt.Printf("Server.foo r1=%#v:%T\n", r1, r1)
	r2 := Call([]interface{}{"Server", "foo"}, []interface{}{0, 0})
	fmt.Printf("Server.foo r2=%#v:%T\n", r2, r2)

	r3 := Call([]interface{}{"Server", "bar"}, []interface{}{4.2, 4.2})
	fmt.Printf("Server.bar r3=%#v:%T\n", r3, r3)
	r4 := Call([]interface{}{"Server", "bar"}, []interface{}{float64(0), float64(0)})
	fmt.Printf("Server.bar r4=%#v:%T\n", r4, r4)

}

var callMutex sync.Mutex

// Call as if via Haxe remote to local haxe function
func Call(path, args interface{}) interface{} {
	callMutex.Lock()
	defer callMutex.Unlock()
	s := haxeremote.Serialize(path) + haxeremote.Serialize(args)
	if C.hxrCalling != 0 {
		panic("C.hxrCalling != 0")
	}
	C.hxrIn = C.CString(s)
	C.hxrCalling = 1
	for C.hxrCalling != 0 { // wait for the call to complete
		runtime.Gosched()
	}
	C.free(unsafe.Pointer(C.hxrIn))
	rs := C.GoString(C.hxrOut)
	rs = strings.TrimPrefix(rs, "hxr")
	rv, _, err := haxeremote.Unserialize([]byte(rs))
	if err != nil {
		return interface{}(err)
	}
	return rv
}
