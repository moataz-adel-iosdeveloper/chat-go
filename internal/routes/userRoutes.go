package routes

import (
	"chat-go/internal/controllers"
	"net/http"
)

func (r *Route) loginHandler(w http.ResponseWriter, req *http.Request) {
	response, statusCode := controllers.Login(req)
	r.WriterResponse(w, req, response, statusCode)
}

func (r *Route) registerHandler(w http.ResponseWriter, req *http.Request) {
	response, statusCode := controllers.Register(req)
	r.WriterResponse(w, req, response, statusCode)
}

func (r *Route) allUserHandler(w http.ResponseWriter, req *http.Request) {
	response, statusCode := controllers.GetAllUsers(req)
	r.WriterResponse(w, req, response, statusCode)
}
