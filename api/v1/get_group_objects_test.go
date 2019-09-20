package v1

import (
	"database/sql"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"testing"
)

func TestGetGroupObjectsApi(t *testing.T) {
	Convey("Given error in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WillReturnError(errors.New("some database error"))

		Convey("When client GET /v1/group-objects/1", func() {
			res, err := http.Get(ts.URL + "/v1/group-objects/1")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNotFoundJsonResponse, `
{
	"Status":"error",
	"Code":"",
	"Message":"some database error"
}
`)
			})
		})
	})

	Convey("Given no groups in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		Convey("When GET /v1/group-objects/1", func() {
			res, err := http.Get(ts.URL + "/v1/group-objects/1")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNotFoundJsonResponse, `
{
	"Status":"error",
	"Code":"ErrGroupNotFound",
	"Message":"Group is not found."
}
`)
			})
		})
	})

	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "world_location_id", "name", "code", "editor_data"}).
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", 1, "Common", "Code", "{}")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WithArgs(1).WillReturnRows(rows)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+group_objects`).WithArgs(1).WillReturnError(sql.ErrNoRows)

		Convey("When GET /v1/group-objects/1", func() {
			res, err := http.Get(ts.URL + "/v1/group-objects/1")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"success",
	"Data":[]
}
`)
			})
		})
	})

	Convey("Given some rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "world_location_id", "name", "code", "editor_data"}).
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", 1, "Common", "Code", "{}")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WithArgs(1).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "group_id", "object_id", "instance_id", "name", "data", "locked", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 1, "Button 1", "{}", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00").
			AddRow(2, 1, 1, 2, "Button 2", "{}", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+group_objects`).WillReturnRows(rows)

		Convey("When GET /v1/group-objects/1", func() {
			res, err := http.Get(ts.URL + "/v1/group-objects/1")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
  "Status": "success",
  "Data": [
    {
      "Id": 1,
      "CreatedAt": "2018-01-01 00:00:00",
      "UpdatedAt": "2018-01-02 00:00:00",
      "Data": {},
      "GroupId": 1,
      "InstanceId": 1,
      "Name": "Button 1",
      "ObjectId": 1
    },
    {
      "Id": 2,
      "CreatedAt": "2018-01-01 00:00:00",
      "UpdatedAt": "2018-01-02 00:00:00",
      "Data": {},
      "GroupId": 1,
      "InstanceId": 2,
      "Name": "Button 2",
      "ObjectId": 1
    }
  ]
}
`)
			})
		})
	})
}
