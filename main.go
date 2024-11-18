package main

import (
	"fmt"
	"log"
	"net/http"
)




func home(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello from home"))
}
func main(){
	

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Print("starting server at port 4000")
	err := http.ListenAndServe(":4000", mux)
	if err !=nil{
		fmt.Println(err)
	}
}