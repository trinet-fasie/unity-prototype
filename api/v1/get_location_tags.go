package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"net/http"
	"strings"
)

const (
	LocationTagsLimit = 20
)

func (a *API) GetLocationTags(w http.ResponseWriter, r *http.Request) {
	var text = r.URL.Query().Get("search")
	text = strings.ToLower(text)

	tags, err := repository.GetLocationTags(a.db, text, LocationTagsLimit)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, tags)
}
