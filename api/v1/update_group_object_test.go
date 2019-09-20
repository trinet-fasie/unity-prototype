package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestUpdateGroupObjectApi(t *testing.T) {
	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+group_objects`).WillReturnError(sql.ErrNoRows)

		Convey("When client POST /v1/update-group-object/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/update-group-object/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"error",
	"Code":"ErrGroupObjectNotFound",
	"Message":"Group object is not found."
}
`)
			})
		})
	})

	Convey("Given some row in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "group_id", "object_id", "instance_id", "name", "data", "locked", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 1, "Button", "{}", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+group_objects`).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"updated_at"}).
			AddRow("2018-01-12 00:00:00")
		sqlMocker.ExpectQuery(`UPDATE\s+group_objects`).WithArgs(2, "Red button", `{"Some":"value"}`, 1).WillReturnRows(rows)

		Convey("When POST /v1/update-group-object/1", func() {
			res, err := http.Post(ts.URL+"/v1/update-group-object/1", "application/json", strings.NewReader(`{
	"Name" : "Red button",
	"InstanceId" : 2,
	"Data" : {"Some" : "value"}
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
    "CreatedAt": "2018-01-01 00:00:00",
    "UpdatedAt": "2018-01-12 00:00:00",
    "Data": {
      "Some": "value"
    },
    "GroupId": 1,
    "InstanceId": 2,
    "Name": "Red button",
    "ObjectId": 1
  }
}
`)
			})
		})
	})
}
