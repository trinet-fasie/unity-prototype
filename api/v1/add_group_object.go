package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (a *API) AddGroupObject(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	groupObject := new(repository.GroupObject)
	json.Unmarshal(reqBytes, &groupObject)

	err = groupObject.Validate()
	if err != nil {
		a.ResponseFail(w, r, err.(validator.ValidationErrors))
		return
	}

	if groupObject.Name == "" {
		groupObject.Name = strconv.Itoa(int(groupObject.InstanceId))
	}

	err = groupObject.Insert(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, groupObject)
}
