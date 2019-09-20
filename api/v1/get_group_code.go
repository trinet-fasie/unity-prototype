package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (a *API) GetGroupCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	groupId, err := strconv.ParseUint(vars["id"], 10, 0)
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

	w.Header().Set("Content-type", "text/plain; charset=utf8")
	w.Write([]byte(group.Code))
}
