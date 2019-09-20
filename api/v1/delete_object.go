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
	ErrObjectNotFound   = errors.New("Object not found.")
	ErrDeleteUsedObject = errors.New("Cannot delete object. Object is used.")
)

func (a *API) DeleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	o, err := repository.GetObjectById(a.db, uint(id))
	if err != nil && err == sql.ErrNoRows {
		a.ReportNotFoundError(w, r, ErrObjectNotFound)
		return
	}
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	if o.Usages > 0 {
		a.ReportBadRequestError(w, r, ErrDeleteUsedObject)
		return
	}

	err = os.RemoveAll(o.ResourcesDirectory())
	if err != nil {
		a.ReportInternalError(w, r, errors.New("Cannot remove resources directory: "+err.Error()))
		return
	}

	err = repository.DeleteObject(a.db, o)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, "")
}
