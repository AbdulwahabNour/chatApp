package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


func main(){
      addr := flag.String("addr", ":8080", "server address")
      flag.Parse()
      r := mux.NewRouter()
      r.HandleFunc("/ws", Wshandler )
      log.Fatalln(http.ListenAndServe(*addr, r))
}