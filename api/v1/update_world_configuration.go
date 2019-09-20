package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (a *API) UpdateWorldConfiguration(w http.ResponseWriter, r *http.Request) {
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

	updatedData := new(repository.WorldConfiguration)
	json.Unmarshal(reqBytes, &updatedData)

	wc, err := repository.GetWorldConfigurationById(a.db, uint(id))
	if wc == nil {
		a.ReportNotFoundError(w, r, ErrWorldConfigurationNotFound)
		return
	}

	if updatedData.Name != "" {
		wc.Name = updatedData.Name
	}

	if updatedData.StartWorldLocationId > 0 {
		wc.StartWorldLocationId = updatedData.StartWorldLocationId
	}

	if updatedData.GroupIds != nil {
		wc.GroupIds = updatedData.GroupIds
	}

	err = wc.Update(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, wc)
}
