package response

type JSONResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CreateJSONResult(msg string, data interface{}) JSONResult {
	return JSONResult{
		Message: msg,
		Data:    data,
	}
}
