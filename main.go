package main

import (
	"net/http"
	"reptile-go/handler"
	"reptile-go/middleware"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.Use(middleware.Cors, mux.CORSMethodMiddleware(r))
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
