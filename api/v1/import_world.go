package v1

import (
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v8"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func (a *API) ImportWorld(w http.ResponseWriter, r *http.Request) {
	requestFile, header, err := r.FormFile("file")
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}
	defer requestFile.Close()

	if strings.ToLower(filepath.Ext(header.Filename)) != ".owws" {
		a.ReportBadRequestError(w, r, errors.New("File type is not allowed"))
		return
	}

	reqBytes, err := ioutil.ReadAll(requestFile)

	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	var s *WorldStructure
	err = json.Unmarshal(reqBytes, &s)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	err = validate.Struct(s)
	if err != nil {
		a.ReportBadRequestError(w, r, err.(validator.ValidationErrors))
		return
	}

	locationsByExternalId, err := getUsedLocationsByExternalId(s, a.db)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	objectsByExternalId, err := getUsedObjectsByExternalId(s, a.db)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	// data integrity validation
	for _, worldLocationData := range s.WorldLocations {
		if locationsByExternalId[worldLocationData.LocationId] == nil {
			a.ReportBadRequestError(w, r, errors.New(fmt.Sprintf("Location with id #%d is not declared in world structure.", worldLocationData.LocationId)))
			return
		}
		for _, groupData := range worldLocationData.Groups {
			for _, groupObjectData := range groupData.AllGroupObjectsData() {
				if objectsByExternalId[groupObjectData.ObjectId] == nil {
					a.ReportBadRequestError(w, r, errors.New(fmt.Sprintf("Object with id #%d is not declared in world structure.", groupObjectData.ObjectId)))
					return
				}
			}
		}
	}

	world := new(repository.World)
	world.Name = s.WorldName

	err = world.Insert(a.db)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	var locationIdByExternalId = make(map[uint]uint)
	var groupIdByExternalId = make(map[uint]uint)
	for _, worldLocationData := range s.WorldLocations {
		worldLocation := new(repository.WorldLocation)
		worldLocation.Sid = worldLocationData.Sid
		worldLocation.Name = worldLocationData.Name
		worldLocation.WorldId = world.Id
		worldLocation.LocationId = locationsByExternalId[worldLocationData.LocationId].Id

		err = worldLocation.Insert(a.db)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
		locationIdByExternalId[worldLocationData.Id] = worldLocation.Id

		for _, groupData := range worldLocationData.Groups {
			group := repository.NewGroup()
			group.WorldLocationId = worldLocation.Id
			group.Name = groupData.Name
			group.Code = groupData.Code
			group.EditorData = groupData.EditorData

			err = group.Insert(a.db)
			if err != nil {
				a.ReportInternalError(w, r, err)
				return
			}
			groupIdByExternalId[groupData.Id] = group.Id

			err = importGroupObjectsBranch(group.Id, 0, groupData.GroupObjects, objectsByExternalId, a.db)
			if err != nil {
				a.ReportInternalError(w, r, err)
				return
			}
		}
	}

	for _, worldConfigurationData := range s.WorldConfigurations {
		worldConfiguration := new(repository.WorldConfiguration)
		worldConfiguration.WorldId = world.Id
		worldConfiguration.Sid = worldConfigurationData.Sid
		worldConfiguration.Name = worldConfigurationData.Name
		worldConfiguration.StartWorldLocationId = locationIdByExternalId[worldConfigurationData.StartWorldLocationId]
		for _, externalId := range worldConfigurationData.GroupIds {
			worldConfiguration.GroupIds = append(worldConfiguration.GroupIds, groupIdByExternalId[externalId])
		}

		err = worldConfiguration.Insert(a.db)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
	}

	world.Configurations = uint(len(s.WorldConfigurations))

	a.ResponseSuccess(w, r, world)
}

func getUsedLocationsByExternalId(s *WorldStructure, db *sql.DB) (result map[uint]*repository.Location, err error) {
	result = make(map[uint]*repository.Location)
	for _, locationData := range s.Locations {
		result[locationData.Id], err = repository.GetLocationByGuid(db, locationData.Guid)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Location %s is not found in library.", locationData.Guid))
		}
	}

	return result, nil
}

func getUsedObjectsByExternalId(s *WorldStructure, db *sql.DB) (result map[uint]*repository.Object, err error) {
	result = make(map[uint]*repository.Object)
	for _, objectData := range s.Objects {
		result[objectData.Id], err = repository.GetObjectByGuid(db, objectData.Guid)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Object %s is not found in library.", objectData.Guid))
		}
	}

	return result, nil
}

func importGroupObjectsBranch(groupId uint, branchId uint, groupObjectsData []*WorldStructureWorldLocationGroupObject, objectsByExternalId map[uint]*repository.Object, db *sql.DB) error {
	for position, data := range groupObjectsData {
		o := new(repository.GroupObject)
		o.GroupId = groupId
		o.ObjectId = objectsByExternalId[data.ObjectId].Id
		o.InstanceId = data.InstanceId
		o.Name = data.Name
		o.Data = data.Data
		o.Locked = data.Locked

		if branchId > 0 {
			o.ParentId.Scan(branchId)
		} else {
			o.ParentId.Scan(nil)
		}
		o.Position = uint(position)

		err := o.Insert(db)
		if err != nil {
			return err
		}

		err = importGroupObjectsBranch(groupId, o.Id, data.GroupObjects, objectsByExternalId, db)
		if err != nil {
			return err
		}
	}

	return nil
}
