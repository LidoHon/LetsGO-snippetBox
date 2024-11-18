package main

import (
	// "flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


func main() {

	fileServer :=http.FileServer(http.Dir("/home/home/Desktop/Golang/LetsGo/snippetbox/ui/static/"))

	// we can always open a file in Go and use it
	// as your log destination.	
	
	/*f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
		log.Fatal(err)
		}
		defer f.Close()
		*/

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)



	err := godotenv.Load("../../.env")
	if err !=nil{
		errorLog.Fatal("error loading .env file")
	}

	port :=os.Getenv("PORT")

	// addr :=flag.String("addr", ":"+port, "HTTP network address")

	srv :=&http.Server{
		Addr: ":"+port,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Print("server is running on port ", port)
	// if err :=http.ListenAndServe(":"+port, mux); err !=nil{
	// 	log.Fatal(err)
	// }
	err =srv.ListenAndServe()
	errorLog.Fatal(err)


}