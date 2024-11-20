package main

import "net/http"


func (app *application) routes() *http.ServeMux{
	mux :=http.NewServeMux()

	fileServer :=http.FileServer(http.Dir("/home/home/Desktop/Golang/LetsGo/snippetbox/ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}