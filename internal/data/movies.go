package data

import ("time"; "greenligh/internal/validator")



type Movie struct{
	ID    		int64		`json:"id"`
	CreatedAt	time.Time 	`json:"-"`
	Title		string 		`json:"title"`
	Year		int32		`json:"year,omitempty"`
	Runtime		Runtime 		`json:"runtime,omitempty,string"`
	Genres		[]string	`json:"genres, omitempty"`
	Version		int32 		`json:"version"`
}

func ValidateMovie(v *validator.Validator, input *Movie){
	
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
}
