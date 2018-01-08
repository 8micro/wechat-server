package api

import "github.com/8micro/wechat-server/etc"
import "github.com/8micro/wechat-server/server"
import "github.com/gorilla/mux"

import (
	"net/http"
)

type handler func(c *Context) error

var routes = map[string]map[string]handler{
	"GET": {
		"/v1/_ping": ping,
	},
}

func NewRouter(enableCors bool, store Store) *mux.Router {

	router := mux.NewRouter()
	for method, mappings := range routes {
		for route, handler := range mappings {
			routemethod := method
			routepattern := route
			routehandler := handler
			wrap := func(w http.ResponseWriter, r *http.Request) {
				if enableCors {
					writeCorsHeaders(w, r)
				}
				c := NewContext(w, r, store)
				routehandler(c)
			}
			router.Path(routepattern).Methods(routemethod).HandlerFunc(wrap)
			if enableCors {
				optionsmethod := "OPTIONS"
				optionshandler := optionsHandler
				wrap := func(w http.ResponseWriter, r *http.Request) {
					if enableCors {
						writeCorsHeaders(w, r)
					}
					c := NewContext(w, r, store)
					optionshandler(c)
				}
				router.Path(routepattern).Methods(optionsmethod).HandlerFunc(wrap)
			}
		}
	}

	router.HandleFunc("/v1/wechat_callback", func(w http.ResponseWriter, r *http.Request) {
		coreServer := store["CoreServer"].(*server.CoreServer)
		coreServer.Server().ServeHTTP(w, r, nil)
	})
	return router
}

func ping(c *Context) error {

	pangData := struct {
		Key          string             `json:"key"`
		SystemConfig *etc.Configuration `json:"systemconfig"`
	}{
		Key:          c.Get("Key").(string),
		SystemConfig: c.Get("SystemConfig").(*etc.Configuration),
	}
	return c.JSON(http.StatusOK, pangData)
}

func optionsHandler(c *Context) error {

	//allowHeaders := r.Header.Get("Access-Control-Request-Headers")

	c.WriteHeader(http.StatusOK)
	return nil
}
