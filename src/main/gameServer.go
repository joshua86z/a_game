package main

import (
	"code.google.com/p/go.net/websocket"
	"controllers"
	"libs/lua"
	"net/http"
	"runtime"
	"strconv"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	L, err := lua.NewLua("conf/app.lua")
	if err != nil {
		panic(err)
	}

	port := L.GetInt("port")

	L.Close()

	http.Handle("/", websocket.Server{Handler: controllers.Handler})
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		panic(err)
	}
}
