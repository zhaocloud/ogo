// Ogo

package ogo

import (
	//"fmt"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
	//"strings"
)

type Controller struct {
	Endpoint string
	//Routes   map[string]Method
}

type C struct { // same as goji/web C, reference: http://golang.org/ref/spec#Type_identity
	URLParams map[string]string
	Env       map[string]interface{}
}

type ControllerInterface interface {
	Get(c web.C, w http.ResponseWriter, r *http.Request)
	Post(c web.C, w http.ResponseWriter, r *http.Request)
	Put(c web.C, w http.ResponseWriter, r *http.Request)
	Delete(c web.C, w http.ResponseWriter, r *http.Request)
	Patch(c web.C, w http.ResponseWriter, r *http.Request)
	Head(c web.C, w http.ResponseWriter, r *http.Request)
	DefaultRoutes(c ControllerInterface)
	//DefaultRoutes()
}

func (ctr *Controller) Get(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
func (ctr *Controller) Post(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
func (ctr *Controller) Put(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
func (ctr *Controller) Delete(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
func (ctr *Controller) Patch(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
func (ctr *Controller) Head(c web.C, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// controller default route
//func DefaultRoutes(c ControllerInterface) {
func (ctr *Controller) DefaultRoutes(c ControllerInterface) {
	//func (c *Controller) DefaultRoutes() {
	goji.Get("/"+ctr.Endpoint+"/:id", c.Get)
	goji.Get("/"+ctr.Endpoint, c.Get)
	goji.Post("/"+ctr.Endpoint, c.Post)
	goji.Delete("/"+ctr.Endpoint+"/:id", c.Delete)
	goji.Patch("/"+ctr.Endpoint+"/:id", c.Get)
	goji.Put("/"+ctr.Endpoint+"/:id", c.Put)
}

func RouteGet(p interface{}, h interface{}) {
	goji.Get(p, h)
}
func RoutePost(p interface{}, h interface{}) {
	goji.Post(p, h)
}
func RoutePut(p interface{}, h interface{}) {
	goji.Put(p, h)
}
func RouteDelete(p interface{}, h interface{}) {
	goji.Delete(p, h)
}
func RoutePatch(p interface{}, h interface{}) {
	goji.Patch(p, h)
}
func RouteHead(p interface{}, h interface{}) {
	goji.Head(p, h)
}

func SetDefaultRoutes(c ControllerInterface) {
	c.DefaultRoutes(c)
}
