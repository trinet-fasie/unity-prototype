package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"net/http"
)

func (a *API) GetLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := repository.GetLocations(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	for _, location := range locations {
		err = location.LoadTags(a.db)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
	}

	a.ResponseSuccess(w, r, locations)
}
