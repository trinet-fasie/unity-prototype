package v1

import (
	"database/sql"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"testing"
)

func TestLocationApi(t *testing.T) {
	Convey("Given error in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WithArgs(1).WillReturnError(errors.New("some database error"))

		Convey("When client GET /v1/locations/1", func() {
			res, err := http.Get(ts.URL + "/v1/locations/1")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"error",
	"Code":"",
	"Message":"some database error"
}
`)
			})
		})
	})

	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		Convey("When GET /v1/locations/1", func() {
			res, err := http.Get(ts.URL + "/v1/locations/1")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"error",
	"Code":"ErrLocationNotFound",
	"Message":"Location is not found."
}
`)
			})
		})
	})

	Convey("Given some rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "guid", "name", "usages", "created_at", "updated_at"}).
			AddRow(1, "c949de70-7e44-42b7-b34e-b84efd1afbf1", "Location 1", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WithArgs(1).WillReturnRows(rows)

		Convey("When GET /v1/locations/1", func() {
			res, err := http.Get(ts.URL + "/v1/locations/1")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
  "Status": "success",
  "Data": {
    "Id": 1,
    "CreatedAt": "2018-01-01 00:00:00",
    "UpdatedAt": "2018-01-02 00:00:00",
    "Guid": "c949de70-7e44-42b7-b34e-b84efd1afbf1",
    "Name": "Location 1",
    "Resources": {
      "Bundle": "/v1/locations/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle",
      "Config": "/v1/locations/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle.json",
      "DllPath": "/v1/locations/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1",
      "Icon": "/v1/locations/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle.png",
      "Manifest": "/v1/locations/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle.manifest"
    },
    "Usages": 0
  }
}
`)
			})
		})
	})
}
