package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (a *API) UpdateWorld(w http.ResponseWriter, r *http.Request) {
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

	updatedData := new(repository.World)
	json.Unmarshal(reqBytes, &updatedData)

	world, err := repository.GetWorldById(a.db, uint(id))
	if world == nil {
		a.ReportNotFoundError(w, r, ErrWorldNotFound)
		return
	}

	if updatedData.Name != "" {
		world.Name = updatedData.Name
	}

	err = world.Update(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, world)
}
