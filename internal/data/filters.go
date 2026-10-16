package data
import ("fmt";  "math"; "strings";  "greenligh/internal/validator") 


type Filters struct{
	Page		int 
	PageSize	int 
	Sort		string
	SortSafeList []string
}

type Metadata struct{
	CurrentPage		int `json:"current_page,omitempty"`
	PageSize		int `json:"page_size,omitempty"`
	FirstPage		int `json:"first_page,omitempty"`
	LastPage		int `json:"last_page,omitempty"`
	TotalRecords	int `json:"total_records,omitempty"`
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList{
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}
func (f Filters) sortDirection() string{
	if strings.HasPrefix(f.Sort, "-"){
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int{
	return f.PageSize 
}
func (f Filters) offset() int{
	return (f.Page -1)*f.PageSize 
}
func ValidateFilters(v *validator.Validator, f Filters){
	v.Check(f.Page < 1, "Page", "Must be greater than 1")
	v.Check(f.Page > 10_000_000, "Page", "Must be less than 10,000,000")
	v.Check(f.PageSize < 1, "Page Size", "Must be greater than 1")
	v.Check(f.PageSize > 100, "Page Size", "Must be less than 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "Sort", "Invalid sort value")

}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0{
		fmt.Printf("totalRecords is %d\n", totalRecords)
		return Metadata{}

	}
	return Metadata{
		CurrentPage:	page,
		PageSize:		pageSize,
		FirstPage:		1,
		LastPage:		int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords:	totalRecords,
	}
}
