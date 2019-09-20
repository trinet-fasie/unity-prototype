package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"encoding/json"
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UpdateGroupRequest struct {
	Name            string
	Code            string
	EditorData      map[string]interface{}
	LockedInstances []uint
}

func (a *API) UpdateGroup(w http.ResponseWriter, r *http.Request) {
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

	request := new(UpdateGroupRequest)
	json.Unmarshal(reqBytes, &request)

	err = validate.Struct(request)
	if err != nil {
		a.ResponseFail(w, r, err.(validator.ValidationErrors))
		return
	}

	g, err := repository.GetGroupById(a.db, uint(id))
	if g == nil {
		a.ReportNotFoundError(w, r, ErrGroupNotFound)
		return
	}

	if request.Name != "" {
		g.Name = request.Name
	}

	var codeChanged = false
	if request.Code != "" && g.Code != request.Code {
		g.Code = request.Code
		codeChanged = true
	}

	if request.EditorData != nil {
		g.EditorData = request.EditorData
	}

	err = g.Update(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	if codeChanged {
		ch, err := a.rabbitMq.Channel()
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
		err = ch.Publish(
			"owd.logic.changed",           // exchange
			fmt.Sprintf("group.%d", g.Id), // routing key
			[]byte(g.Code),
			wabbit.Option{
				"contentType": "text/plain",
			},
		)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
	}

	if request.LockedInstances != nil {
		err = repository.UpdateGroupObjectsLocks(a.db, g.Id, request.LockedInstances)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
	}

	a.ResponseSuccess(w, r, g)
}
