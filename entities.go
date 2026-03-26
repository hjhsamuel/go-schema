package go_schema

type HttpField struct {
	Url         string `json:"url"`                  // http request url with url params
	Method      string `json:"method"`               // http method, GET/POST, default: GET
	MultiSelect bool   `json:"multi_select"`         // single or multiple selected
	UserData    string `json:"user_data"`            // custom user data
	UserInput   []any  `json:"user_input,omitempty"` // user selected items
}
