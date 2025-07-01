package response

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}
type Responsefieldempty struct {
	Meta Meta_error `json:"Meta_error"`
}

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Meta_error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Colum   string `json:"colum"`
}

func SuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}

func ErrorResponse(code int, message string) Response {
	return Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: nil,
	}
}
func Errorfieldempty(code int, colum string) Responsefieldempty {
	return Responsefieldempty{
		Meta: Meta_error{
			Code:    code,
			Message: "column cannot be empty",
			Colum:   colum,
		},
	}
}
