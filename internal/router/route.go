package router


import (
	"net/http"
)

func StudentRoute() *http.ServeMux {
	router:=http.NewServeMux()

	return router;
}