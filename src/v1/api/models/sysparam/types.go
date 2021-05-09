package sysparam

type SysParam struct {
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
}

/**
Get
*/
type ResponseSysParam []SysParam

/**
Params
*/

type ParamsUpdateSysParam struct {
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
}
