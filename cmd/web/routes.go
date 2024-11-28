package main

import (
	"net/http"

	"github.com/LidoHon/LetsGO-snippetBox.git/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)


func (app *application) routes() http.Handler{
	// initialize the router
	router :=httprouter.New()


	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		app.notFound(w)
	})

	// fileServer :=http.FileServer(http.Dir("../../ui/static/"))
	//now that we are using the embeded system to serve our files we can use the code bellow
	fileServer :=http.FileServer(http.FS(ui.Files))

	// router.Handler(http.MethodGet,"/static/*filepath", http.StripPrefix("/static", fileServer))

	/*Our static files are contained in the "static" folder of the ui.File embedded filesystem. So, for example, our CSS stylesheet is located at "static/css/main.css". This means that we now longer need to strip the prefix from the request URL -- any requests that start with /static/ can
	just be passed directly to the file server and the correspondingstatic
	file will be served (so long as it exists).*/
	
	router.Handler(http.MethodGet,"/static/*filepath", fileServer)

	// creating a new middleware specific to our dynamic application routes

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)


	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet,"/snippet/view/:id", dynamic.ThenFunc(app.snippetView))


	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))

	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))

	router.Handler(http.MethodGet, "/user/login",
	dynamic.ThenFunc(app.userLogin))

	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))


// protected routes
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc( app.snippetCreate))
	
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))




	// our middleware chain
	standard :=alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}