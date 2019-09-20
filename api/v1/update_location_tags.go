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

func (a *API) UpdateLocationTags(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	var newTags []*repository.LocationTag
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
	l, err := repository.GetLocationById(a.db, uint(id))
	if err != nil && err == sql.ErrNoRows {
		a.ReportNotFoundError(w, r, ErrLocationNotFound)
		return
	}
	if err != nil {
		a.ReportInternalError(w, r, ErrLocationNotFound)
		return
	}
	err = l.LoadTags(a.db)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	err = repository.UpdateLocationTags(a.db, l, l.Tags, newTags)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccessNoContent(w, r)
}
