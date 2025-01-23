package providers

import (
	ctx "github.com/gophish/gophish/context"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ImportLDAPUsers(r *http.Request) ([]models.Target, error) {
	vars := mux.Vars(r)
	ldapId, _ := strconv.ParseInt(vars["id"], 0, 64)
	l, err := models.GetLDAP(ldapId, ctx.Get(r, "user_id").(int64))

	ts := []models.Target{}
	if err != nil {
		return ts, err
	}

	conn, err := l.Connect()
	if err != nil {
		return ts, err
	}

	defer conn.Close()

	return l.ImportLDAPTargets(conn)
}
