package v1

import (
	"archive/zip"
	"github.com/trinet-fasie/unity-prototype/api/repository"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type InstallLocationInstructions struct {
	Guid string `validate:"required"`
	Name string `validate:"required"`
	Tags []string
}

func (il *InstallLocationInstructions) Validate() error {
	return validate.Struct(il)
}

func (a *API) InstallLocation(w http.ResponseWriter, r *http.Request) {
	requestFile, header, err := r.FormFile("file")
	if err != nil {
		a.ReportBadRequestError(w, r, err)
		return
	}
	defer requestFile.Close()

	if strings.ToLower(filepath.Ext(header.Filename)) != ".zip" {
		a.ReportBadRequestError(w, r, errors.New("File type is not allowed"))
		return
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "upload-location")
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}
	defer os.Remove(tmpFile.Name())

	size, err := io.Copy(tmpFile, requestFile)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	zipReader, err := zip.NewReader(tmpFile, size)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	var installZipFile *zip.File
	for _, zipFile := range zipReader.File {
		if !zipFile.FileInfo().IsDir() && zipFile.Name == "install.json" {
			installZipFile = zipFile
			break
		}
	}
	if installZipFile == nil {
		a.ReportBadRequestError(w, r, errors.New("Cannot find install.json"))
		return
	}

	installFile, err := installZipFile.Open()
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}
	defer installFile.Close()

	installFileContents, err := ioutil.ReadAll(installFile)
	if err != nil {
		a.ReportInternalError(w, r, err)
		return
	}

	var install = new(InstallLocationInstructions)
	err = json.Unmarshal(installFileContents, &install)
	if err != nil {
		a.ReportInternalError(w, r, errors.New("Install config is not a valid json: "+err.Error()))
		return
	}

	err = install.Validate()
	if err != nil {
		a.ReportBadRequestError(w, r, errors.New("Invalid install config: "+err.Error()))
		return
	}

	l, err := repository.GetLocationByGuid(a.db, install.Guid)
	if err != nil && err != sql.ErrNoRows {
		a.ReportInternalError(w, r, errors.New("Cannot get location by guid: "+err.Error()))
		return
	}

	created := false
	if l == nil {
		l = &repository.Location{
			Guid: install.Guid,
			Name: install.Name,
		}
		err = l.Insert(a.db)
		if err != nil {
			a.ReportInternalError(w, r, errors.New("Cannot create location: "+err.Error()))
			return
		}
		created = true

		var tags []*repository.LocationTag
		for _, tagText := range install.Tags {
			tagText = strings.TrimSpace(tagText)
			tagText = strings.ToLower(tagText)
			tag, err := repository.GetLocationTagByText(a.db, tagText)
			if err != nil && err != sql.ErrNoRows {
				a.ReportInternalError(w, r, err)
				return
			}
			if err == sql.ErrNoRows {
				tag = new(repository.LocationTag)
				tag.Text = tagText

				err = tag.Insert(a.db)
				if err != nil {
					a.ReportInternalError(w, r, err)
					return
				}
			}
			tags = append(tags, tag)
		}
		err = repository.UpdateLocationTags(a.db, l, l.Tags, tags)
		if err != nil {
			a.ReportInternalError(w, r, errors.New("Cannot set tags: "+err.Error()))
			return
		}
	} else {
		l.Name = install.Name
		err = l.Update(a.db)
		if err != nil {
			a.ReportInternalError(w, r, errors.New("Cannot updated location: "+err.Error()))
			return
		}
	}

	resourcesDir := l.ResourcesDirectory()

	os.RemoveAll(resourcesDir)
	err = os.MkdirAll(resourcesDir, 0755)
	if err != nil {
		a.ReportInternalError(w, r, errors.New("Cannot create resources directory: "+err.Error()))
	}

	for _, zipFile := range zipReader.File {
		if zipFile.FileInfo().IsDir() {
			continue
		}
		err = a.UnzipResource(zipFile, resourcesDir)
		if err != nil {
			a.ReportInternalError(w, r, errors.New("Cannot unpack from zip: "+err.Error()))
			return
		}
	}

	a.ResponseSuccess(w, r, struct {
		Created  bool
		Location *repository.Location
	}{
		Created:  created,
		Location: l,
	})
}
