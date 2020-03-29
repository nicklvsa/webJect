package shared

//consts used for the api
const (
	MaxUpload = 10 * 1024 * 1024
	TempDir   = "../TEMP"
)

//Logger - holds the logger structure
type Logger int

//APIResponse - holds the response structure
type APIResponse struct {
	ErrorCode      int         `json:"err_code"`
	ErrorMessage   interface{} `json:"err_msg,omitempty"`
	ContentMessage interface{} `json:"content_msg,omitempty"`
}
