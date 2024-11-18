package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)




func home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w, r)
	}else{
		w.Write([]byte("Hello from home"))
	}
	
}

func snippetView(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("here is a single view"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte ("method is not allowed"))
		// or we could do it like this looks better and clean
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		
		
		
	}else{
		w.Write([]byte("here is for creating a single view"))
	}
	
}
func main(){
	err := godotenv.Load()
	if err !=nil{
		log.Fatal("error loading the .env")
	}

	port :=os.Getenv("PORT")
	if port == ""{
		port="4000"
	}

	
// route definition
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/view", snippetView)
	mux.HandleFunc("/create", snippetCreate)

// listen and serve
	fmt.Printf("Starting server at port %v\n", port)
	err = http.ListenAndServe(":"+port, mux)
	if err !=nil{
		fmt.Println(err)
	}
}