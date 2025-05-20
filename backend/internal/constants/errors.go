package constants

import "net/http"

var ErrorMessages = map[int]string{
	http.StatusInternalServerError: "The server encountered an unexpected error and could not process your request",
}
