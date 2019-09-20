package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (a *API) UpdateGroupObject(w http.ResponseWriter, r *http.Request) {
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

	updatedData := new(repository.GroupObject)
	json.Unmarshal(reqBytes, &updatedData)

	o, err := repository.GetGroupObjectById(a.db, uint(id))
	if o == nil {
		a.ReportNotFoundError(w, r, ErrGroupObjectNotFound)
		return
	}

	if updatedData.Name != "" {
		o.Name = updatedData.Name
	}
	if updatedData.InstanceId != 0 {
		o.InstanceId = updatedData.InstanceId
	}
	if updatedData.Data != nil {
		o.Data = updatedData.Data
	}

	err = o.Update(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, o)
}
