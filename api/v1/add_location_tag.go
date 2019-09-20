package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"encoding/json"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"net/http"
	"strings"
)

type AddLocationTagDto struct {
	Id        uint
	CreatedAt string
	UpdatedAt string
	Text      string `validate:"required,max=25"`
}

func (a *API) AddLocationTag(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	dto := new(AddLocationTagDto)
	json.Unmarshal(reqBytes, &dto)

	dto.Text = strings.TrimSpace(dto.Text)
	dto.Text = strings.ToLower(dto.Text)

	err = validate.Struct(dto)
	if err != nil {
		a.ResponseFail(w, r, err.(validator.ValidationErrors))
		return
	}

	t, err := repository.GetLocationTagByText(a.db, dto.Text)
	if err != nil && err != sql.ErrNoRows {
		a.ReportInternalError(w, r, err)
		return
	}
	if err == sql.ErrNoRows {
		t = new(repository.LocationTag)
		t.Text = dto.Text

		err = t.Insert(a.db)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
	}

	dto.Id = t.Id
	dto.Text = t.Text
	dto.CreatedAt = t.CreatedAt
	dto.UpdatedAt = t.UpdatedAt

	a.ResponseSuccess(w, r, dto)
}
