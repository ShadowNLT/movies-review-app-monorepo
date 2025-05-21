package constants

import "net/http"

var ErrorMessages = map[int]string{
	http.StatusInternalServerError: "The server encountered an unexpected error and could not process your request",
	http.StatusNotFound:            "The requested resource could not be found",
}

var FormattedErrorMessages = map[int]string{
	http.StatusMethodNotAllowed: "The %s method is not supported for this resource",
}
