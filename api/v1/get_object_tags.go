package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"net/http"
	"strings"
)

const (
	ObjectTagsLimit = 20
)

func (a *API) GetObjectTags(w http.ResponseWriter, r *http.Request) {
	var text = r.URL.Query().Get("search")
	text = strings.ToLower(text)

	tags, err := repository.GetObjectTags(a.db, text, ObjectTagsLimit)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	a.ResponseSuccess(w, r, tags)
}
