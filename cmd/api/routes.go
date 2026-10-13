package main 

import ("net/http"; "github.com/julienschmidt/httprouter")

func (app *application) routes() *httprouter.Router{
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundError)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedError)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	return router
}
