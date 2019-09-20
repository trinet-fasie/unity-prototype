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
	ErrGroupNotFound = errors.New("Group is not found.")
)

func (a *API) GetGroupObjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	groupId, err := strconv.ParseUint(vars["groupId"], 10, 0)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	group, err := repository.GetGroupById(a.db, uint(groupId))
	if err != nil && err == sql.ErrNoRows {
		a.ReportNotFoundError(w, r, ErrGroupNotFound)
		return
	}
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	objects, err := repository.GetGroupObjectsTreeByGroupId(a.db, group.Id)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, objects)
}
