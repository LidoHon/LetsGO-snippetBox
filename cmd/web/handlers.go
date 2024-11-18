package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w,r)
		return
	}

	// files := []string{
	// 	"../../ui/html/base.tmpl",
	// 	"../../ui/html/pages/home.tmpl",
	// 	}
	// ts, err := template.ParseFiles(files...)
	ts, err := template.ParseFiles("/home/home/Desktop/Golang/LetsGo/snippetbox/ui/html/pages/home.tmpl")
	if err !=nil{
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err !=nil{
		log.Print(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}


func snippetView(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	 if err !=nil || id<1{
		http.NotFound(w, r)
		return

	 }
	 fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)	 
}


func snippetCreate(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}


