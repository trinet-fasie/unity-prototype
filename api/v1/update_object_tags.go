package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (a *API) UpdateObjectTags(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	var newTags []*repository.ObjectTag
	err = json.Unmarshal(reqBytes, &newTags)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

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
		a.ReportInternalError(w, r, ErrObjectNotFound)
		return
	}
	err = o.LoadTags(a.db)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	err = repository.UpdateObjectTags(a.db, o, o.Tags, newTags)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccessNoContent(w, r)
}
