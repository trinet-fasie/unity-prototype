package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"net/http"
)

func (a *API) GetObjects(w http.ResponseWriter, r *http.Request) {
	objects, err := repository.GetObjects(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	for _, object := range objects {
		err = object.LoadTags(a.db)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
	}

	a.ResponseSuccess(w, r, objects)
}
