package response

type JSONResult struct {
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

func BuildJSONResponse(msg string, data interface{}) JSONResult {
	return JSONResult{
		Message: msg,
		Body:    data,
	}
}
