package main 

import ("fmt"; "greenligh/internal/validator"; "io"; "encoding/json"; "errors"; "net/http"; "net/url"; "strconv"; "strings"; "github.com/julienschmidt/httprouter")

type envelope map[string]any

func (app *application) readIDParam( r *http.Request) (int64, error){
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1{
		return 0, errors.New("Invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error{
	js, err := json.MarshalIndent(data, "", "\t");
	if err != nil{
		return err 
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error{
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
			case errors.As(err, &syntaxError):

				return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)

			case errors.As(err, &unmarshalTypeError):

				if unmarshalTypeError.Field != "" {

					return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)

				}

				return fmt.Errorf("body contains badly formed JSON (at character %d)", unmarshalTypeError.Offset)

			case errors.Is(err, io.EOF):

				return errors.New("body must not be empty")

			case errors.As(err, &invalidUnmarshalError):

				panic(err)


			case strings.HasPrefix(err.Error(), "json: unknown field "):

				fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")

				return fmt.Errorf("body contains unknown field %q", fieldName)
			case errors.As(err, &maxBytesError):

				return fmt.Errorf("body must not be longer than %d bytes", maxBytesError.Limit)

			default:
				return err
		}
	}
	err = decoder.Decode(&struct{}{})
	if err != io.EOF{
		return fmt.Errorf("body must contain only a single JSON value")
	}
	return nil
}

func (app *application) readInt(v *validator.Validator,qr url.Values, key string, defaultValue int) int {
	param := qr.Get(key)
	if param == ""{
		return defaultValue
	}

	paramInt, err := strconv.Atoi(param)
	if err != nil{
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return paramInt
}

func (app *application) readString(qr url.Values, key string, defaultValue string) string{
	param := qr.Get(key)
	if param == ""{
		return defaultValue
	}
	return param
}

func (app *application) readCSV(qr url.Values, key string, defaultValue []string) []string{
	csv := qr.Get(key)
	if csv == ""{
		return defaultValue 
	}

	return strings.Split(csv, ",")
}

func (app *application) background(fn func()){
	go func() {
		defer func(){
			if err := recover(); err != nil{
				app.logger.Printf("Error from background work %v\n", err)
			}
		}()

		fn() 
	}()
}
