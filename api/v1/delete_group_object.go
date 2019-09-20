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
	ErrGroupObjectLocked = errors.New("Group object is locked for deletion.")
)

func (a *API) DeleteGroupObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["groupObjectId"], 10, 0)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	o, err := repository.GetGroupObjectById(a.db, uint(id))
	if err != nil && err == sql.ErrNoRows {
		a.ReportNotFoundError(w, r, ErrGroupObjectNotFound)
		return
	}
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	if o.Locked {
		a.ReportBadRequestError(w, r, ErrGroupObjectLocked)
		return
	}

	err = repository.DeleteGroupObject(a.db, o)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, "")
}
