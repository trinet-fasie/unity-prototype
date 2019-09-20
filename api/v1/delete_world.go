package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	ErrWorldNotFound = errors.New("World is not found.")
)

func (a *API) DeleteWorld(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	o, err := repository.GetWorldById(a.db, uint(id))
	if err != nil && err == sql.ErrNoRows {
		a.ReportNotFoundError(w, r, ErrWorldNotFound)
		return
	}
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	err = repository.DeleteWorld(a.db, o)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, "")
}
