package main

import (
	"log"
	"net/http"

	"github.com/tardisgo/haxeremote/hxrhttp"
)

func main() {
	hxrhttp.AddFunc("Server.foo", func(data interface{}) interface{} {
		return data.([]interface{})[0].(int) + data.([]interface{})[1].(int)
	})
	hxrhttp.AddFunc("Server.bar", func(data interface{}) interface{} {
		return data.([]interface{})[0].(string) + " " + data.([]interface{})[1].(string)
	})
	hxrhttp.AddFunc("Server.fad", func(data interface{}) interface{} {
		return data.([]interface{})[0].(float64) + data.([]interface{})[1].(float64)
	})
	hxrhttp.AddFunc("Server.dong", func(data interface{}) interface{} {
		return []interface{}{
			data.([]interface{})[1].(float64), data.([]interface{})[0].(float64)}
	})
	hxrhttp.AddFunc("Server.dingbat", func(data interface{}) interface{} {
		item := data.([]interface{})[0].([]byte)
		//println("Length of []byte ", len(item))
		for i := range item {
			item[i] = 'A' + byte(i)
		}
		return item
	})
	http.HandleFunc("/_haxeRPC_", hxrhttp.HttpHandler)
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("../client"))))

	println("Haxe remote webserver running on port 8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
