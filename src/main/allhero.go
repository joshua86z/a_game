package main

import (
	"code.google.com/p/go.net/websocket"
	"controllers"
	"fmt"
	"libs/lua"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	var port int
	if Lua, err := lua.NewLua("conf/app.lua"); err != nil {
		panic(err)
	} else {
		port = Lua.GetInt("port")
		Lua.Close()
	}

	http.Handle("/", websocket.Server{Handler: controllers.Handler})
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		fmt.Println("Panic: ", err)
		os.Exit(1)
	}
}
