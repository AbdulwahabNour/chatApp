package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
)


func main(){
      runtime.GOMAXPROCS(1)
      addr := flag.String("addr", ":8080", "server address")
      flag.Parse()

      wsServer := NewWsServer()
      go wsServer.Run()

      //tells fileserve serves local directory
      f := http.FileServer(http.Dir("./view/"))
      r := mux.NewRouter()
 
      r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
            Wshandler(w, r, wsServer)
      } )
     
       r.PathPrefix("/").Handler(f)
  
      
       srv := &http.Server{
             Handler: r,
             Addr: fmt.Sprint("127.0.0.1",*addr),
          
       }
      
      
 
 
   
      log.Fatalln(srv.ListenAndServe())
} 


 