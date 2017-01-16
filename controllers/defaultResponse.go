package controller

const ErrResponseInternalServerError string = "{\"status\" : false, \"message\" : \"Internal server error\", \"errCode\": 500}"
const ErrResponseUnsupportedContentType string = "{\"status\" : false, \"message\" : \"Unsupported Content-Type\", \"errCode\": 415}"

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	ErrCode int    `json:"errCode,omitempty"`
}

func newResponse() Response {
	return Response{}
}
