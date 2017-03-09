package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func GetRouter() *httprouter.Router {
	router := httprouter.New()
	router.PanicHandler = PanicHandler
	router.GET("/v1/ping", Ping)

	return router
}

func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("pong"))
}

func PanicHandler(w http.ResponseWriter, r *http.Request, c interface{}) {
	log.Printf("Recovering from panic, Reason: %+v", c.(error))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(c.(error).Error()))
}
