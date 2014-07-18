package base

type (
	// Type TmplData represents data in the templates.
	TmplData map[string]interface{}
	TplName  string

	ApiJsonErr struct {
		Message string `json:"message"`
		DocUrl  string `json:"documentation_url"`
	}
)

var GoGetMetas = make(map[string]bool)
