package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)


func (app *application) newTemplateData(r *http.Request) *templateData{
	// since it is giving me unused error even though i used it in handler.go to shut the compiler up it up i used the r parameter in here even if it's just to assign it to a variable
	_ = r // This line tells the compiler that r is being used
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken: nosurf.Token(r),
	}
}
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData){
	ts, ok := app.templateCache[page]
	if !ok{
		err :=fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// init a new buffer
	buf :=new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err !=nil{
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) serverError(w http.ResponseWriter, err  error){
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	if app.debug{
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

func(app *application) clientError (w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter){
	app.clientError(w, http.StatusNotFound)
}


func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {

	var invalidDecoderError *form.InvalidDecoderError
	if errors.As(err, &invalidDecoderError) {
		panic(err)
	}

	return err
	}
	return nil
	}



func (app *application) isAuthenticated( r *http.Request) bool{
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
