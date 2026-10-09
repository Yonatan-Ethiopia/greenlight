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

	v := validator.New()

	v.Check(input.Title != "", "title", "must not be empty")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "can not be in the future")
	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be a positive number")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "Genres", "must contain atleast 1 genre")
	v.Check(len(input.Genres) <= 5, "Genres", "must not be more than 5")

	v.Check(validator.Unique(input.Genres), "Genres", "must not contain duplicate genres")

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
