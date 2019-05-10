package account

import (
	"log"

	"github.com/gorilla/mux"
)

func newRouter(logger *log.Logger) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		handler := route.HandlerFunc

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}
