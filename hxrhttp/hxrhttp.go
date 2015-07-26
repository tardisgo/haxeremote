// package hxrhttp enables a Go webserver to act as a Haxe remoting server
package hxrhttp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/tardisgo/haxeremote"
)

var haxeRemoteFuncs = map[string]func(interface{}) interface{}{}
var haxeRemoteFuncsMutex sync.Mutex

// AddFunc adds a remote function available for execution
func AddFunc(name string, fn func(interface{}) interface{}) {
	haxeRemoteFuncsMutex.Lock()
	haxeRemoteFuncs[name] = fn
	haxeRemoteFuncsMutex.Unlock()
}

func callHaxeRemoteFunc(name string, arg interface{}) (interface{}, error) {
	haxeRemoteFuncsMutex.Lock()
	fn, ok := haxeRemoteFuncs[name]
	haxeRemoteFuncsMutex.Unlock()
	if !ok {
		return nil, errors.New("could not find remote func " + name + " in CallHaxeRemoteFunc()")
	}
	return fn(arg), nil
}

// HttpHandler provides an http handler for Haxe remoting calls
func HttpHandler(rw http.ResponseWriter, req *http.Request) {
	// TODO add further validaton of the request as from a valid Haxe user here?

	if req.Header.Get("X-Haxe-Remoting") != "1" {
		rw.WriteHeader(http.StatusBadRequest)
		log.Panic("header does not contain X-Haxe-Remoting=1")
		return
	}

	p, e := ioutil.ReadAll(req.Body)
	//fmt.Printf("DEBUG Request: %#v\nBody: %s, err=%v\n", *req, p, e)
	if e != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Panic(e)
		return
	}
	if len(p) < 4 {
		rw.WriteHeader(http.StatusBadRequest)
		log.Panic("no body")
		return
	}
	if string(p[:4]) != "__x=" {
		rw.WriteHeader(http.StatusBadRequest)
		log.Panic("__x= not found")
		return
	}
	une, err := url.QueryUnescape(string(p[4:]))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Panic(err)
		return
	}
	buf := []byte(une)
	//fmt.Printf("DEBUG URL unescaped = %s\n", string(buf))
	var targetA interface{}
	target := ""
	targetA, buf, err = haxeremote.Unserialize(buf)
	for i, t := range targetA.([]interface{}) {
		if i > 0 {
			target += "."
		}
		target += t.(string)
	}
	//fmt.Printf("DEBUG Unserialized Target decoded=%s, remaining=%s, error=%v\n", target, buf, err)
	var args interface{}
	args, buf, err = haxeremote.Unserialize(buf)
	//fmt.Printf("DEBUG Unserialized Args decoded=%v, remaining=%s, error=%v\n", args, buf, err)

	results, err := callHaxeRemoteFunc(target, args)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Panic(err)
		return
	}
	reply := "hxr" + haxeremote.Serialize(results)
	fmt.Fprintln(rw, reply)
	//fmt.Printf("DEBUG haxe http remote results: %v serialized-reply: %s\n", results, reply)
}
