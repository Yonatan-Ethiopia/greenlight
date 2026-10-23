package main 

import ("fmt"; "errors"; "net/http"; "greenligh/internal/data"; "greenligh/internal/validator")

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request){
	var input struct{
		Name		string	`json:"name"`
		Email	 	string	`json:"email"`
		Password	string	`json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil{
		app.badRequestError(w, r, err)
		return
	}

	user := &data.User{
		Name:	input.Name,
		Email:	input.Email,
		Activated:	false,
	}

	err = user.Password.Set(input.Password)
	if err != nil{
		app.serverErrorResponse(w, r, err)
		return 
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return 
	}

	err = app.models.Users.Insert(user)
	fmt.Printf("Error is: %v\n", err)

	if err != nil{
				switch{
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return 
	}
	app.background(func(){
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil{
			app.logger.Printf("Error is %v",err)
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil{
		app.serverErrorResponse(w, r, err)
		return 
	}
}
