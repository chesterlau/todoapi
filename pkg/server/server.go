package server

import "todo/pkg/router"

func Server() {
	r := router.Init()
	r.Logger.Fatal(r.Start(":1323"))
}
