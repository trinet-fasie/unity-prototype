package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"net/http"
)

func (a *API) AddWorld(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	world := new(repository.World)
	json.Unmarshal(reqBytes, &world)

	err = world.Validate()
	if err != nil {
		a.ResponseFail(w, r, err.(validator.ValidationErrors))
		return
	}

	err = world.Insert(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, world)
}
