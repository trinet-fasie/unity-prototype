package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"net/http"
)

func (a *API) GetWorlds(w http.ResponseWriter, r *http.Request) {
	worlds, err := repository.GetWorlds(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, worlds)
}
