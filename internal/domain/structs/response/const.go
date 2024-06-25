package response

type HttpResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	HANDLER_INTERR = &Error{
		Code:    500,
		Message: "HandlerError: Terjadi Kesalahan Sistem",
	}

	SERVICE_INTERR = &Error{
		Code:    500,
		Message: "ServiceError: Terjadi Kesalahan Sistem",
	}

	VIOLATED_ERR = &Error{
		Code:    406,
		Message: "Data Tidak Dapat Dihapus Karena Berelasi Dengan Data Lain",
	}
)

func BADREQ_ERR(message string) error {
	return &Error{
		Code:    400,
		Message: message,
	}
}

func NOTFOUND_ERR(message string) error {
	return &Error{
		Code:    404,
		Message: message,
	}
}

func UNAUTH_ERR(message string) error {
	return &Error{
		Code:    401,
		Message: message,
	}
}

func SUCCESS_RES(message string) *HttpResponse {
	return &HttpResponse{
		Success: true,
		Message: message,
	}
}

func ERROR_RES(message string) *HttpResponse {
	return &HttpResponse{
		Success: false,
		Message: message,
	}
}

func DATA_RES(data interface{}) *HttpResponse {
	return &HttpResponse{
		Success: true,
		Data:    data,
	}
}
