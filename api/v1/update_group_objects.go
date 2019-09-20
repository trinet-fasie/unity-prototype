package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NeowayLabs/wabbit"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"net/http"
)

type UpdateGroupObjectData struct {
	Id           uint
	ObjectId     uint `validate:"required"`
	InstanceId   uint `validate:"required"`
	Name         string
	Data         map[string]interface{}   `validate:"required"`
	GroupObjects []*UpdateGroupObjectData `validate:"dive"`
}

type UpdateGroupObjects struct {
	GroupId      uint                     `validate:"required"`
	GroupObjects []*UpdateGroupObjectData `validate:"required"`
	Objects      []*repository.Object
}

func (a *API) UpdateGroupObjects(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	var dto = new(UpdateGroupObjects)
	json.Unmarshal(reqBytes, &dto)

	err = validate.Struct(dto)
	if err != nil {
		a.ResponseFail(w, r, err.(validator.ValidationErrors))
		return
	}

	g, err := repository.GetGroupById(a.db, dto.GroupId)
	if g == nil {
		a.ReportNotFoundError(w, r, ErrGroupNotFound)
		return
	}

	wl, err := repository.GetWorldLocationById(a.db, g.WorldLocationId)
	if wl == nil {
		a.ReportNotFoundError(w, r, ErrWorldLocationNotFound)
		return
	}

	currentGroupObjectsMap, err := repository.GetGroupObjectsMapByGroupId(a.db, g.Id)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	usedGroupObjectIds := getUsedGroupObjectIds(dto.GroupObjects)
	var groupObjectsToDelete []*repository.GroupObject

nextGroupObject:
	for _, o := range currentGroupObjectsMap {
		for _, usedId := range usedGroupObjectIds {
			if o.Id == usedId {
				continue nextGroupObject
			}
		}

		if o.Locked {
			a.ReportBadRequestError(w, r, errors.New(fmt.Sprintf("Cannot delete group object #%d. Object is locked.", o.Id)))
			return
		}

		groupObjectsToDelete = append(groupObjectsToDelete, o)
		delete(currentGroupObjectsMap, o.Id)
	}
	for _, groupObject := range groupObjectsToDelete {
		err = repository.DeleteGroupObject(a.db, groupObject)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
	}

	err = persistGroupObjectsBranch(g.Id, 0, dto.GroupObjects, currentGroupObjectsMap, a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	ch, err := a.rabbitMq.Channel()
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	var objectsMap = make(map[uint]bool)
	for _, groupObjectData := range flatternGroupObjectsData(dto.GroupObjects) {
		if !objectsMap[groupObjectData.ObjectId] {
			o, err := repository.GetObjectById(a.db, groupObjectData.ObjectId)
			if err != nil {
				a.ReportInternalError(w, r, err)
				return
			}
			objectsMap[groupObjectData.ObjectId] = true
			dto.Objects = append(dto.Objects, o)
		}
	}

	message, err := json.Marshal(dto)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	err = ch.Publish(
		"owd.group_objects.changed",                        // exchange
		fmt.Sprintf("world.%d.group.%d", wl.WorldId, g.Id), // routing key
		message,
		wabbit.Option{
			"contentType": "application/json",
		},
	)
	if err != nil {
		a.ReportInternalError(w, r, errors.New(fmt.Sprintf("Cannot publish event owd.group_objects.changed: %s", err.Error())))
		return
	}

	a.ResponseSuccess(w, r, dto)
}

func persistGroupObjectsBranch(groupId uint, branchId uint, groupObjectsData []*UpdateGroupObjectData, currentGroupObjectsMap map[uint]*repository.GroupObject, db *sql.DB) error {
	for position, data := range groupObjectsData {
		var o *repository.GroupObject
		if data.Id == 0 {
			o = new(repository.GroupObject)
			o.ObjectId = data.ObjectId
			o.GroupId = groupId
		} else {
			o = currentGroupObjectsMap[data.Id]
			if o == nil {
				return errors.New(fmt.Sprintf("Cannot persist objects. Object with id #%d is not found.", data.Id))
			}
		}

		o.InstanceId = data.InstanceId
		o.Data = data.Data
		if branchId > 0 {
			o.ParentId.Scan(branchId)
		} else {
			o.ParentId.Scan(nil)
		}
		o.Position = uint(position)

		err := o.InsertOrUpdate(db)
		if err != nil {
			return err
		}

		data.Id = o.Id
		data.Name = o.Name

		err = persistGroupObjectsBranch(groupId, o.Id, data.GroupObjects, currentGroupObjectsMap, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func getUsedGroupObjectIds(groupObjectsData []*UpdateGroupObjectData) (result []uint) {
	for _, data := range groupObjectsData {
		if data.Id > 0 {
			result = append(result, data.Id)
		}
		result = append(result, getUsedGroupObjectIds(data.GroupObjects)...)
	}

	return result
}

func flatternGroupObjectsData(groupObjectsData []*UpdateGroupObjectData) (result []*UpdateGroupObjectData) {
	for _, data := range groupObjectsData {
		result = append(result, data)
		result = append(result, flatternGroupObjectsData(data.GroupObjects)...)
	}

	return result
}
