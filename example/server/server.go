package main

import (
	"log"
	"net/http"

	"github.com/tardisgo/haxeremote"
)

func main() {
	haxeremote.AddHaxeRemoteFunc("Server.foo", func(data interface{}) interface{} {
		return data.([]interface{})[0].(int) + data.([]interface{})[1].(int)
	})
	haxeremote.AddHaxeRemoteFunc("Server.bar", func(data interface{}) interface{} {
		return data.([]interface{})[0].(string) + data.([]interface{})[1].(string)
	})
	haxeremote.AddHaxeRemoteFunc("Server.fad", func(data interface{}) interface{} {
		return data.([]interface{})[0].(float64) + data.([]interface{})[1].(float64)
	})
	haxeremote.AddHaxeRemoteFunc("Server.dong", func(data interface{}) interface{} {
		return []interface{}{
			data.([]interface{})[1].(float64), data.([]interface{})[0].(float64)}
	})
	http.HandleFunc("/haxeremote", haxeremote.HaxeRemoteHandler)
	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir("../client"))))

	println("Haxe remote webserver running on port 8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
