package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

var (
	ErrDeleteUsedLocation = errors.New("Cannot delete location. Location is used.")
)

func (a *API) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	o, err := repository.GetLocationById(a.db, uint(id))
	if err != nil && err == sql.ErrNoRows {
		a.ReportNotFoundError(w, r, ErrLocationNotFound)
		return
	}
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	if o.Usages > 0 {
		a.ReportBadRequestError(w, r, ErrDeleteUsedLocation)
		return
	}

	err = os.RemoveAll(o.ResourcesDirectory())
	if err != nil {
		a.ReportInternalError(w, r, errors.New("Cannot remove resources directory: "+err.Error()))
	}

	err = repository.DeleteLocation(a.db, o)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, "")
}
