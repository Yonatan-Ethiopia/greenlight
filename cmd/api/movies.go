package main 

import  ("fmt"; "errors"; "net/http"; "greenligh/internal/data"; "greenligh/internal/validator")

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request){
	var input struct{
		Title	string	`json:"title"`
		Year	int32	`json:"year"`
		Runtime	data.Runtime	`json:"runtime"`
		Genres	[]string	`json:"genres"`
	}

	//err := json.NewDecoder(r.Body).Decode(&input)
	err := app.readJSON(w, r, &input)
	if err != nil{
		app.badRequestError(w, r, err)
		return
	}

	movie := data.Movie{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	}
	v := validator.New()
	data.ValidateMovie(v, &movie)


	if !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}


	err = app.models.Movies.Insert(&movie)
	if err != nil{
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%id", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil{
	app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r) 
	if err != nil || id < 1{
		app.notFoundError(w, r)
	return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil{
		if errors.Is(err, data.ErrRecordNotFound){
			app.notFoundError(w, r)
			return
		}else {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil{
		app.logger.Printf("Error encountered: %v \n", err)
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil{
		app.notFoundError(w, r)
		return
	}
	var input struct{
		Title *string	`json:"title"`
		Year *int32		`json:"year"`
		Runtime *data.Runtime	`json:"runtime"`
		Genres []string	`json:"genres"`
	}
	
	err = app.readJSON(w, r, &input)
	if err != nil{
		app.badRequestError(w, r, err)
	}

	movie,err := app.models.Movies.Get(id)
	if err != nil{
		switch{
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundError(w, r)
				return
			default:
				app.serverErrorResponse(w, r, err)
				return
		}
	}

	if input.Title != nil{
		movie.Title = *input.Title 
	}
	if input.Year != nil{
		movie.Year = *input.Year 
	}
	if input.Runtime != nil{
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil{
		movie.Genres = input.Genres
	}
	v := validator.New()
	data.ValidateMovie(v, movie)

	if !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
	}

	err = app.models.Movies.Update(movie)
	if err != nil{
		switch{
			case errors.Is(err, data.ErrEditConflict):
				app.editConflictResponse(w,r)
				return
			default:
				app.serverErrorResponse(w, r, err)
				return
		}
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie":movie}, nil)
	if err != nil{
		app.serverErrorResponse(w, r, err)
		return
	}


}

func (app *application) deleteHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r)
	if err != nil{
		app.notFoundError(w,r)
		return 
	}
	err = app.models.Movies.Delete(id)
	if err != nil{
		switch{
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundError(w, r)
				return 
			default:
				app.serverErrorResponse(w, r, err)
				return
		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Movie deleted succefully"}, nil)
	if err != nil{
		app.serverErrorResponse(w,r, err)
		return
	}
	return
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request){
	var input struct{
		Title string 
		Genres []string 
		Page int 
		PageSize int 
		Sort string
	}
	params := r.URL.Query()
	input.Title = app.readString(params, "title", "")
	input.Genres = app.readCSV(params, "genres", []string{})

	v := validator.New()
	input.Page = app.readInt(v, params, "page", 1)
	input.PageSize = app.readInt(v, params, "page_size", 20)

	input.Sort = app.readString(params, "sort", "id")

	if !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)

}
