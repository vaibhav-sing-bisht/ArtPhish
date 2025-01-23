package api

import (
	"encoding/json"
	"errors"
	ctx "github.com/gophish/gophish/context"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// LdapProfiles handles requests for the /api/ldap/ endpoint
func (as *Server) LdapProfiles(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		ls, err := models.GetLDAPs(ctx.Get(r, "user_id").(int64))
		if err != nil {
			log.Error(err)
		}
		JSONResponse(w, ls, http.StatusOK)
	//POST: Create a new LDAP and return it as JSON
	case r.Method == "POST":
		l := models.LDAP{}
		// Put the request into a page
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
			return
		}
		// Check to make sure the name is unique
		_, err = models.GetLDAPByName(l.Name, ctx.Get(r, "user_id").(int64))
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			JSONResponse(w, models.Response{Success: false, Message: "LDAP name already in use"}, http.StatusConflict)
			log.Error(err)
			return
		}
		l.ModifiedDate = time.Now().UTC()
		l.UserId = ctx.Get(r, "user_id").(int64)
		err = models.PostLDAP(&l)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, l, http.StatusCreated)
	}
}

// LdapProfile contains functions to handle the GET'ing, DELETE'ing, and PUT'ing
// of a LDAP object
func (as *Server) LdapProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 0, 64)
	l, err := models.GetLDAP(id, ctx.Get(r, "user_id").(int64))
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "SMTP not found"}, http.StatusNotFound)
		return
	}
	switch {
	case r.Method == "GET":
		JSONResponse(w, l, http.StatusOK)
	case r.Method == "DELETE":
		err = models.DeleteLDAP(id, ctx.Get(r, "user_id").(int64))
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error deleting LDAP"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "LDAP Deleted Successfully"}, http.StatusOK)
	case r.Method == "PUT":
		l = models.LDAP{}
		err = json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			log.Error(err)
		}
		if l.Id != id {
			JSONResponse(w, models.Response{Success: false, Message: "/:id and /:ldap_id mismatch"}, http.StatusBadRequest)
			return
		}
		err = l.Validate()
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		l.ModifiedDate = time.Now().UTC()
		l.UserId = ctx.Get(r, "user_id").(int64)
		err = models.PutLDAP(&l)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Error updating page"}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, l, http.StatusOK)
	}
}
