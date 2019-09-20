package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (a *API) UpdateWorldLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	updatedData := new(repository.WorldLocation)
	json.Unmarshal(reqBytes, &updatedData)

	wl, err := repository.GetWorldLocationById(a.db, uint(id))
	if wl == nil {
		a.ReportNotFoundError(w, r, ErrWorldLocationNotFound)
		return
	}

	if updatedData.LocationId > 0 {
		wl.LocationId = updatedData.LocationId
	}

	if updatedData.Name != "" {
		wl.Name = updatedData.Name
	}

	err = wl.Update(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, wl)
}
