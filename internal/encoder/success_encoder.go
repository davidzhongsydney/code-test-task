package encoder

func NewHTTPSuccess(data interface{}) *HTTPSuccess {
	return &HTTPSuccess{
		Code: 200,
		Data: data,
	}
}

type HTTPSuccess struct {
	Code int         `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func FromResponse(data interface{}) *HTTPSuccess {
	return NewHTTPSuccess(data)
}
