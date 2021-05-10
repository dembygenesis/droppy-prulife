package response_builder

import (
	"fmt"
	"reflect"
)

type ResponseSuccessOperation struct {
	Response
	Data []string `json:"data,omitempty"`
}

type Response struct {
	HttpCode        int         `json:"httpCode"`
	ResponseMessage string      `json:"responseMessage"`
	OperationStatus string      `json:"operationStatus,omitempty"`
	Data            interface{} `json:"data"`
	Errors          []string    `json:"-"` // Always leave this out in the json response, this is just a container
	Pagination      *Pagination  `json:"pagination,omitempty"`
}

type Pagination struct {
	Pages       []int `json:"pages"`
	RowsPerPage int   `json:"rowsPerPage"`
	Offset      int   `json:"offset"`
	Rows        int   `json:"rows"`
	Page        int   `json:"page"`
	TotalCount  int   `json:"total_count"`
	ResultCount int   `json:"result_count"`
}

func (p *Pagination) SetData(
	rowsPerPage int,
	offset int,
	pages []int,
	rows int,
	page int,
	totalCount int,
	resultCount int,
) {
	p.RowsPerPage = rowsPerPage
	p.Offset = offset
	p.Pages = pages
	p.Rows = rows
	p.Page = page
	p.TotalCount = totalCount
	p.ResultCount = resultCount
}

func (r *Response) SetUpdateSuccess() {
	r.OperationStatus = "UPDATE_SUCCESS"
}

func (r *Response) SetCreateSuccess() {
	r.OperationStatus = "CREATE_SUCCESS"
}

func (r *Response) SetDeleteSuccess() {
	r.OperationStatus = "DELETE_SUCCESS"
}

// SetErrors returns the error list in an array format.
// It checks the error whether it is of type "error", "string", or "array of strings".
// It then converts them into an "array of strings" format except if it is already formatted
// that way.
func (r *Response) SetErrors(i interface{}) {
	var errors []string

	dataType := reflect.TypeOf(i).String()
	fmt.Println("================= SetErrors Data type: ", dataType)


	if dataType == "*errors.errorString" {
		fmt.Println("================= Jay Chou: ")
		val := i.(error)
		errors = append(errors, val.Error())
	} else {
		fmt.Println("================= Else: ", dataType)
	}


	r.Data = struct {
		Errors interface{} `json:"errors"`
	}{
		// Old one - but this time we're going to enforce an array of strings
		// Errors: i,
		Errors: errors,
	}
}

func (r *Response) AddErrors(error string) {
	// Get current value of errors, append a new one - then make that the interface's value
	currentErrors := append(r.Errors, error)

	r.Errors = currentErrors

	errorStruct := make(map[string][]string)
	errorStruct["errors"] = currentErrors

	r.Data = errorStruct
	// r.Data = r.Errors
}

func (r *Response) SetResponseData(data interface{}) {
	r.Data = data
}

type ResponseErrors struct {
	Errors map[string][]string
}

/*func (r *Response) AddErrors(error string) {
	if len(r.Error) == 0 {
		r.Error = make(map[string][]string)
		r.Error["errors"] = []string{}
	}

	r.Error["errors"] = append(r.Error["errors"], error)
}*/

func (r *ResponseSuccessOperation) AddData(text string) {
	/*if len(r.Error) == 0 {
		r.Error = []string
		r.Error["errors"] = []string{}
	}*/

	r.Data = append(r.Data, text)
}

func (r *ResponseErrors) AddErrors(error string) {
	if len(r.Errors) == 0 {
		r.Errors = make(map[string][]string)
		r.Errors["errors"] = []string{}
	}

	r.Errors["errors"] = append(r.Errors["errors"], error)
}

func (r *ResponseErrors) GetErrors() map[string][]string {
	return r.Errors
}
