package v1

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"strings"
	"testing"
)

func TestAddWorldLocationApi(t *testing.T) {
	Convey("Given invalid request", t, func() {
		ts, _ := GetTestHelpers(t)

		Convey("When client POST /v1/add-world-location", func() {
			res, err := http.Post(ts.URL+"/v1/add-world-location", "application/json", strings.NewReader(`{}`))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidBadRequestJsonResponse, `
{
  "Status": "fail",
  "Data": {
    "WorldLocation.LocationId": {
      "ActualTag": "required",
      "Field": "LocationId",
      "FieldNamespace": "WorldLocation.LocationId",
      "Kind": 7,
      "Name": "LocationId",
      "NameNamespace": "LocationId",
      "Param": "",
      "Tag": "required",
      "Type": {},
      "Value": 0
    },
    "WorldLocation.Name": {
      "ActualTag": "required",
      "Field": "Name",
      "FieldNamespace": "WorldLocation.Name",
      "Kind": 24,
      "Name": "Name",
      "NameNamespace": "Name",
      "Param": "",
      "Tag": "required",
      "Type": {},
      "Value": ""
    },
    "WorldLocation.WorldId": {
      "ActualTag": "required",
      "Field": "WorldId",
      "FieldNamespace": "WorldLocation.WorldId",
      "Kind": 7,
      "Name": "WorldId",
      "NameNamespace": "WorldId",
      "Param": "",
      "Tag": "required",
      "Type": {},
      "Value": 0
    }
  }
}
`)
			})
		})
	})

	Convey("Given valid request", t, func() {
		ts, sqlMocker := GetTestHelpers(t)

		Convey("When client POST /v1/add-world-location", func() {
			rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
			sqlMocker.ExpectQuery(`INSERT\s+INTO\s+world_locations`).WithArgs(
				"741bdb12-aa6b-40c4-8268-fe52a76c0b1e", // Sid
				1,                                      // WorldId
				1,                                      // LocationId
				"Test world location",                  // Name
			).WillReturnRows(rows)

			res, err := http.Post(ts.URL+"/v1/add-world-location", "application/json", strings.NewReader(`{
	"Sid": "741bdb12-aa6b-40c4-8268-fe52a76c0b1e",
	"LocationId" : 1,
	"WorldId" : 1,
	"Name" : "Test world location"
}`))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
  "Status": "success",
  "Data": {
    "Id": 1,
	"Sid": "741bdb12-aa6b-40c4-8268-fe52a76c0b1e",
    "CreatedAt": "2018-01-01 00:00:00",
    "UpdatedAt": "2018-01-02 00:00:00",
    "LocationId": 1,
    "WorldId": 1,
    "Name": "Test world location"
  }
}
`)
			})
		})
	})
}
