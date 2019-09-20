package v1

import (
	"database/sql"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"testing"
)

func TestGetObjectsApi(t *testing.T) {
	Convey("Given error in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+objects`).WillReturnError(errors.New("some database error"))

		Convey("When client GET /v1/objects", func() {
			res, err := http.Get(ts.URL + "/v1/objects")

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
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+objects`).WillReturnError(sql.ErrNoRows)

		Convey("When GET /v1/objects", func() {
			res, err := http.Get(ts.URL + "/v1/objects")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
  "Status": "success",
  "Data": []
}
`)
			})
		})
	})

	Convey("Given some rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "guid", "type_name", "name", "usages", "created_at", "updated_at"}).
			AddRow(1, "c949de70-7e44-42b7-b34e-b84efd1afbf1", "Trigger", "Button", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00").
			AddRow(2, "c949de70-7e44-42b7-b34e-b84efd1afbf2", "Display", "Display", 1, "2018-01-01 00:00:00", "2018-01-02 00:00:00").
			AddRow(3, "c949de70-7e44-42b7-b34e-b84efd1afbf3", "Valve", "Valve", 2, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+objects`).WillReturnRows(rows)

		Convey("When GET /v1/objects", func() {
			res, err := http.Get(ts.URL + "/v1/objects")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
  "Status": "success",
  "Data": [
    {
      "CreatedAt": "2018-01-01 00:00:00",
      "Guid": "c949de70-7e44-42b7-b34e-b84efd1afbf1",
      "Id": 1,
      "Name": "Button",
      "Resources": {
        "Bundle": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle",
        "Config": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle.json",
        "Icon": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle.png",
        "Manifest": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf1/bundle.manifest"
      },
      "Type": "Trigger",
      "UpdatedAt": "2018-01-02 00:00:00",
      "Usages": 0
    },
    {
      "CreatedAt": "2018-01-01 00:00:00",
      "Guid": "c949de70-7e44-42b7-b34e-b84efd1afbf2",
      "Id": 2,
      "Name": "Display",
      "Resources": {
        "Bundle": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf2/bundle",
        "Config": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf2/bundle.json",
        "Icon": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf2/bundle.png",
        "Manifest": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf2/bundle.manifest"
      },
      "Type": "Display",
      "UpdatedAt": "2018-01-02 00:00:00",
      "Usages": 1
    },
    {
      "CreatedAt": "2018-01-01 00:00:00",
      "Guid": "c949de70-7e44-42b7-b34e-b84efd1afbf3",
      "Id": 3,
      "Name": "Valve",
      "Resources": {
        "Bundle": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf3/bundle",
        "Config": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf3/bundle.json",
        "Icon": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf3/bundle.png",
        "Manifest": "/v1/objects/resources/c9/c949de70-7e44-42b7-b34e-b84efd1afbf3/bundle.manifest"
      },
      "Type": "Valve",
      "UpdatedAt": "2018-01-02 00:00:00",
      "Usages": 2
    }
  ]
}
`)
			})
		})
	})
}
