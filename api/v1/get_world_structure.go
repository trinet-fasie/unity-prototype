package v1

import (
	"archive/tar"
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type WorldStructure struct {
	WorldId             uint
	WorldName           string                         `validate:"required"`
	WorldLocations      []*WorldStructureWorldLocation `validate:"dive"`
	WorldConfigurations []*repository.WorldConfiguration
	Locations           []*repository.Location
	Objects             []*repository.Object
}

type WorldStructureWorldLocationGroupObject struct {
	Id           uint
	Name         string `validate:"required"`
	InstanceId   uint   `validate:"required"`
	ObjectId     uint
	Data         map[string]interface{}
	Locked       bool
	GroupObjects []*WorldStructureWorldLocationGroupObject `validate:"dive"`
}

func (g *WorldStructureWorldLocationGroupObject) Descendants() (result []*WorldStructureWorldLocationGroupObject) {
	result = g.GroupObjects
	for _, child := range g.GroupObjects {
		result = append(result, child.Descendants()...)
	}
	return result
}

type WorldStructureWorldLocationGroup struct {
	Id           uint                                      `validate:"required"`
	Name         string                                    `validate:"required"`
	GroupObjects []*WorldStructureWorldLocationGroupObject `validate:"dive"`
	Code         string
	EditorData   map[string]interface{}
}

func (g *WorldStructureWorldLocationGroup) AllGroupObjectsData() (result []*WorldStructureWorldLocationGroupObject) {
	result = g.GroupObjects
	for _, child := range g.GroupObjects {
		result = append(result, child.Descendants()...)
	}
	return result
}

type WorldStructureWorldLocation struct {
	Id         uint                                `validate:"required"`
	Sid        string                              `validate:"required"`
	Name       string                              `validate:"required"`
	LocationId uint                                `validate:"required"`
	Groups     []*WorldStructureWorldLocationGroup `validate:"dive"`
}

func (a *API) GetWorldStructure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["worldId"], 10, 0)
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}

	world, err := repository.GetWorldById(a.db, uint(id))
	if err != nil && err == sql.ErrNoRows {
		a.ReportNotFoundError(w, r, ErrWorldNotFound)
		return
	}
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	worldStructure := &WorldStructure{
		WorldId:   world.Id,
		WorldName: world.Name,
		Locations: make([]*repository.Location, 0),
		Objects:   make([]*repository.Object, 0),
	}

	worldLocations, err := repository.GetWorldLocationsByWorldId(a.db, world.Id)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}
	worldStructure.WorldLocations = make([]*WorldStructureWorldLocation, 0, len(worldLocations))

	var locationsMap = make(map[uint]*repository.Location)
	var objectsMap = make(map[uint]*repository.Object)
	for _, worldLocation := range worldLocations {
		if locationsMap[worldLocation.LocationId] == nil {
			locationsMap[worldLocation.LocationId], err = repository.GetLocationById(a.db, worldLocation.LocationId)
			if err != nil {
				a.ReportInternalError(w, r, err)
				return
			}
			worldStructure.Locations = append(worldStructure.Locations, locationsMap[worldLocation.LocationId])
		}

		wl := &WorldStructureWorldLocation{
			Id:         worldLocation.Id,
			Sid:        worldLocation.Sid,
			Name:       worldLocation.Name,
			LocationId: worldLocation.LocationId,
		}

		groups, err := repository.GetGroupsByWorldLocationId(a.db, worldLocation.Id)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
		wl.Groups = make([]*WorldStructureWorldLocationGroup, 0, len(groups))

		for _, group := range groups {
			g := &WorldStructureWorldLocationGroup{
				Id:         group.Id,
				Name:       group.Name,
				Code:       group.Code,
				EditorData: group.EditorData,
			}

			groupObjectsByIdMap, err := repository.GetGroupObjectsMapByGroupId(a.db, g.Id)
			if err != nil {
				a.ReportInternalError(w, r, err)
				return
			}

			groupObjectsTree, err := repository.BuildGroupObjectsTree(groupObjectsByIdMap)
			if err != nil {
				a.ReportInternalError(w, r, err)
				return
			}

			for _, groupObject := range groupObjectsByIdMap {
				if objectsMap[groupObject.ObjectId] == nil {
					objectsMap[groupObject.ObjectId], err = repository.GetObjectById(a.db, groupObject.ObjectId)
					if err != nil {
						a.ReportInternalError(w, r, err)
						return
					}
					worldStructure.Objects = append(worldStructure.Objects, objectsMap[groupObject.ObjectId])
				}
			}

			g.GroupObjects = formWorldStructureWorldLocationGroupObject(groupObjectsTree)

			wl.Groups = append(wl.Groups, g)
		}

		worldStructure.WorldLocations = append(worldStructure.WorldLocations, wl)
	}

	worldStructure.WorldConfigurations, err = repository.GetWorldConfigurationsByWorldId(a.db, uint(id))
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	_, download := r.URL.Query()["export"]
	if download {
		resBytes, err := json.Marshal(worldStructure)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}

		w.Header().Set("Content-type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.owws\"", world.Name))
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(resBytes)), 10))

		w.Write(resBytes)
		return
	}

	_, build := r.URL.Query()["build"]
	if build {
		packageFile, err := ioutil.TempFile(os.TempDir(), "build-package")
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}
		defer os.Remove(packageFile.Name())

		err = a.buildPackage(worldStructure, packageFile)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}

		packageStat, err := packageFile.Stat()
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}

		w.Header().Set("Content-type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.omw\"", world.Name))
		w.Header().Set("Content-Length", strconv.FormatInt(packageStat.Size(), 10))

		packageFile.Seek(0, 0)
		_, err = io.Copy(w, packageFile)
		if err != nil {
			a.ReportInternalError(w, r, err)
			return
		}

		return
	}

	a.ResponseSuccess(w, r, worldStructure)
}

func (a *API) buildPackage(worldStructure *WorldStructure, packageFile *os.File) (err error) {
	tarWriter := tar.NewWriter(packageFile)
	defer tarWriter.Close()

	// write index file
	structureBytes, err := json.Marshal(worldStructure)
	if err != nil {
		return err
	}
	header := &tar.Header{
		Name: "index.json",
		Mode: 0666,
		Size: int64(len(structureBytes)),
	}
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = tarWriter.Write(structureBytes)
	if err != nil {
		return err
	}

	for _, location := range worldStructure.Locations {
		err = a.TarDirectory(location.ResourcesDirectory(), tarWriter)
		if err != nil {
			return err
		}
	}

	for _, object := range worldStructure.Objects {
		err = a.TarDirectory(object.ResourcesDirectory(), tarWriter)
		if err != nil {
			return err
		}
	}

	return nil
}

func formWorldStructureWorldLocationGroupObject(groupObjects []*repository.GroupObject) []*WorldStructureWorldLocationGroupObject {
	result := make([]*WorldStructureWorldLocationGroupObject, 0, len(groupObjects))

	for _, groupObject := range groupObjects {
		result = append(result, &WorldStructureWorldLocationGroupObject{
			Id:           groupObject.Id,
			Name:         groupObject.Name,
			InstanceId:   groupObject.InstanceId,
			ObjectId:     groupObject.ObjectId,
			Locked:       groupObject.Locked,
			Data:         groupObject.Data,
			GroupObjects: formWorldStructureWorldLocationGroupObject(groupObject.GroupObjects),
		})
	}

	return result
}
