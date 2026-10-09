package main 

import  ("fmt"; "time"; "net/http"; "greenligh/internal/data"; "greenligh/internal/validator")

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


	fmt.Fprintf(w, "%+v \n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request){
	id, err := app.readIDParam(r) 
	if err != nil || id < 1{
		app.notFoundError(w, r)
	return
	}
	movieData := data.Movie{
		ID: 	id,
		CreatedAt:	time.Now(),
		Title:		"Ace ventura: Pet detective",
		//Year:		1998,
		Runtime:	72,
		Genres:		[]string{"comedy","mystrey"},
		Version:	1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movieData}, nil)
	if err != nil{
		app.logger.Printf("Error encountered: %v \n", err)
		app.serverErrorResponse(w, r, err)
	}
}
