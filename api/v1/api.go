package v1

import (
	"archive/tar"
	"archive/zip"
	"database/sql"
	"encoding/json"
	"github.com/NeowayLabs/wabbit"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v8"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type API struct {
	rabbitMq wabbit.Conn
	db       *sql.DB
}

var validate = validator.New(&validator.Config{TagName: "validate"})

func New(router *mux.Router, rabbitMq wabbit.Conn, db *sql.DB) *API {
	a := new(API)
	a.rabbitMq = rabbitMq
	a.db = db

	ch, err := rabbitMq.Channel()
	if err != nil {
		log.Fatalf("Cannot create channel: %s", err)
	}

	err = ch.ExchangeDeclare(
		"owd.logic.changed",
		"topic",
		wabbit.Option{
			"durable":  true,
			"delete":   false,
			"internal": false,
			"noWait":   false,
		},
	)
	if err != nil {
		log.Fatalf("Cannot declare exchange owd.logic.changed: %s", err)
	}

	err = ch.ExchangeDeclare(
		"owd.group_objects.changed",
		"topic",
		wabbit.Option{
			"durable":  true,
			"delete":   false,
			"internal": false,
			"noWait":   false,
		},
	)
	if err != nil {
		log.Fatalf("Cannot declare exchange owd.group_objects.changed: %s", err)
	}

	r := router.PathPrefix("/v1").Subrouter()
	r.Use(CorsMiddleware)

	r.HandleFunc("/objects", a.GetObjects).Methods("GET", "OPTIONS")
	r.HandleFunc("/install-object", a.InstallObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/delete-object/{id:[0-9]+}", a.DeleteObject).Methods("DELETE", "POST", "OPTIONS")

	r.HandleFunc("/locations", a.GetLocations).Methods("GET", "OPTIONS")
	r.HandleFunc("/install-location", a.InstallLocation).Methods("POST", "OPTIONS")
	r.HandleFunc("/locations/{id:[0-9]+}", a.GetLocation).Methods("GET", "OPTIONS")
	r.HandleFunc("/delete-location/{id:[0-9]+}", a.DeleteLocation).Methods("DELETE", "POST", "OPTIONS")

	r.HandleFunc("/add-world-location", a.AddWorldLocation).Methods("PUT", "POST", "OPTIONS")
	r.HandleFunc("/delete-world-location/{id:[0-9]+}", a.DeleteWorldLocation).Methods("DELETE", "POST", "OPTIONS")
	r.HandleFunc("/update-world-location/{id:[0-9]+}", a.UpdateWorldLocation).Methods("POST", "OPTIONS")

	r.HandleFunc("/add-group", a.AddGroup).Methods("PUT", "POST", "OPTIONS")
	r.HandleFunc("/delete-group/{id:[0-9]+}", a.DeleteGroup).Methods("DELETE", "POST", "OPTIONS")
	r.HandleFunc("/update-group/{id:[0-9]+}", a.UpdateGroup).Methods("POST", "OPTIONS")

	r.HandleFunc("/group-objects/{groupId:[0-9]+}", a.GetGroupObjects).Methods("GET", "OPTIONS")
	r.HandleFunc("/update-group-objects", a.UpdateGroupObjects).Methods("POST", "OPTIONS")
	r.HandleFunc("/update-group-object/{id:[0-9]+}", a.UpdateGroupObject).Methods("POST", "OPTIONS")
	r.HandleFunc("/can-delete-group-object/{groupObjectId:[0-9]+}", a.CanDeleteGroupObject).Methods("GET", "OPTIONS")

	r.HandleFunc("/worlds", a.GetWorlds).Methods("GET", "OPTIONS")
	r.HandleFunc("/world-structure/{worldId:[0-9]+}", a.GetWorldStructure).Methods("GET", "OPTIONS")
	r.HandleFunc("/add-world", a.AddWorld).Methods("PUT", "POST", "OPTIONS")
	r.HandleFunc("/import-world", a.ImportWorld).Methods("PUT", "POST", "OPTIONS")
	r.HandleFunc("/update-world/{id:[0-9]+}", a.UpdateWorld).Methods("POST", "OPTIONS")
	r.HandleFunc("/delete-world/{id:[0-9]+}", a.DeleteWorld).Methods("DELETE", "POST", "OPTIONS")

	r.HandleFunc("/add-world-configuration", a.AddWorldConfiguration).Methods("PUT", "POST", "OPTIONS")
	r.HandleFunc("/delete-world-configuration/{id:[0-9]+}", a.DeleteWorldConfiguration).Methods("DELETE", "POST", "OPTIONS")
	r.HandleFunc("/update-world-configuration/{id:[0-9]+}", a.UpdateWorldConfiguration).Methods("POST", "OPTIONS")

	r.HandleFunc("/object-tags", a.GetObjectTags).Methods("GET", "OPTIONS")
	r.HandleFunc("/add-object-tag", a.AddObjectTag).Methods("PUT", "POST", "OPTIONS")
	r.HandleFunc("/update-object-tags/{id:[0-9]+}", a.UpdateObjectTags).Methods("POST", "OPTIONS")

	r.HandleFunc("/location-tags", a.GetLocationTags).Methods("GET", "OPTIONS")
	r.HandleFunc("/add-location-tag", a.AddLocationTag).Methods("PUT", "POST", "OPTIONS")
	r.HandleFunc("/update-location-tags/{id:[0-9]+}", a.UpdateLocationTags).Methods("POST", "OPTIONS")

	dataRouter := router.PathPrefix("/data").Subrouter()

	fs := http.FileServer(http.Dir("data/"))
	dataRouter.PathPrefix("/objects/resources/").Handler(
		http.StripPrefix(
			"/data/",
			a.ObjectResourcesPermissionsWrapper(fs)))

	dataRouter.PathPrefix("/locations/resources/").Handler(
		http.StripPrefix(
			"/data/",
			a.ObjectResourcesPermissionsWrapper(fs)))

	return a
}

func (a *API) ObjectResourcesPermissionsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (a *API) ReportInternalError(w http.ResponseWriter, req *http.Request, err error) {
	dump, _ := httputil.DumpRequest(req, false)
	log.Printf("api: %s error: %s ", string(dump), err)
	w.WriteHeader(http.StatusInternalServerError)

	errBytes, _ := json.Marshal(struct {
		Status  string
		Code    string
		Message string
	}{
		Status:  "error",
		Code:    strconv.Itoa(http.StatusInternalServerError),
		Message: err.Error(),
	})

	w.Header().Set("Content-type", "application/json; charset=utf8")
	w.Write(errBytes)
}

func (a *API) ReportNotFoundError(w http.ResponseWriter, req *http.Request, err error) {
	dump, _ := httputil.DumpRequest(req, false)
	log.Printf("api: %s error: %s ", string(dump), err)
	w.WriteHeader(http.StatusNotFound)

	errBytes, _ := json.Marshal(struct {
		Status  string
		Code    string
		Message string
	}{
		Status:  "error",
		Code:    strconv.Itoa(http.StatusNotFound),
		Message: err.Error(),
	})

	w.Header().Set("Content-type", "application/json; charset=utf8")
	w.Write(errBytes)
}

func (a *API) ReportForbiddenError(w http.ResponseWriter, req *http.Request, err error) {
	dump, _ := httputil.DumpRequest(req, false)
	log.Printf("api: %s error: %s ", string(dump), err)
	w.WriteHeader(http.StatusForbidden)

	errBytes, _ := json.Marshal(struct {
		Status  string
		Code    string
		Message string
	}{
		Status:  "error",
		Code:    strconv.Itoa(http.StatusForbidden),
		Message: err.Error(),
	})

	w.Header().Set("Content-type", "application/json; charset=utf8")
	w.Write(errBytes)
}

func (a *API) ReportBadRequestError(w http.ResponseWriter, req *http.Request, err error) {
	dump, _ := httputil.DumpRequest(req, false)
	log.Printf("api: %s error: %s ", string(dump), err)
	w.WriteHeader(http.StatusBadRequest)

	errBytes, _ := json.Marshal(struct {
		Status  string
		Code    string
		Message string
	}{
		Status:  "error",
		Code:    strconv.Itoa(http.StatusBadRequest),
		Message: err.Error(),
	})

	w.Header().Set("Content-type", "application/json; charset=utf8")
	w.Write(errBytes)
}

func (a *API) ResponseSuccess(w http.ResponseWriter, req *http.Request, data interface{}) {
	resBytes, _ := json.Marshal(struct {
		Status string
		Data   interface{}
	}{
		Status: "success",
		Data:   data,
	})

	w.Header().Set("Content-type", "application/json; charset=utf8")
	w.Write(resBytes)
}

func (a *API) ResponseSuccessNoContent(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) ResponseFail(w http.ResponseWriter, req *http.Request, fails validator.ValidationErrors) {
	w.WriteHeader(http.StatusBadRequest)

	resBytes, _ := json.Marshal(struct {
		Status string
		Data   interface{}
	}{
		Status: "fail",
		Data:   fails,
	})

	w.Header().Set("Content-type", "application/json; charset=utf8")
	w.Write(resBytes)
}

func (a *API) UnzipResource(src *zip.File, destDir string) error {
	zipFileReader, err := src.Open()
	if err != nil {
		return err
	}
	defer zipFileReader.Close()

	destFile, err := os.OpenFile(filepath.Join(destDir, src.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, zipFileReader)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) TarDirectory(src string, dest *tar.Writer) error {
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(file, string(filepath.Separator))

		err = dest.WriteHeader(header)
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		_, err = io.Copy(dest, f)
		if err != nil {
			return err
		}

		f.Close()

		return nil
	})
}
