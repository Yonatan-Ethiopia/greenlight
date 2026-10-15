package data
import ("greenligh/internal/validator") 


type Filters struct{
	Page		int 
	PageSize	int 
	Sort		string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filters){
	v.Check(f.Page < 1, "Page", "Must be greater than 1")
	v.Check(f.Page > 10_000_000, "Page", "Must be less than 10,000,000")
	v.Check(f.PageSize < 1, "Page Size", "Must be greater than 1")
	v.Check(f.PageSize > 100, "Page Size", "Must be less than 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "Sort", "Invalid sort value")

}
