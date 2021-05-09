package utils

// Pagination - handles the pagination structure of your database
type Pagination struct {
	Pages       []int `json:"pages"`
	RowsPerPage int   `json:"rowsPerPage"`
	Offset      int   `json:"offset"`
	Rows        int   `json:"rows"`
	Page        int   `json:"page"`
	TotalCount  int   `json:"total_count"`
	ResultCount int   `json:"result_count"`
}

type Response struct {
	HttpCode        int         `json:"httpCode"`
	ResponseMessage string      `json:"responseMessage"`
	OperationStatus string      `json:"operationStatus,omitempty"`
	Data            interface{} `json:"data"`
	Errors          []string    `json:"-"` // Always leave this out in the json response, this is just a container
	Pagination      Pagination  `json:"pagination,omitempty"`
}