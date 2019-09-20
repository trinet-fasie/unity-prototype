package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"net/http"
)

func (a *API) AddWorldConfiguration(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	wc := &repository.WorldConfiguration{}
	json.Unmarshal(reqBytes, &wc)

	err = wc.Validate()
	if err != nil {
		a.ResponseFail(w, r, err.(validator.ValidationErrors))
		return
	}

	err = wc.Insert(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, wc)
}
