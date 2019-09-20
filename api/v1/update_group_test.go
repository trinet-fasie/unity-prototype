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

func TestUpdateGroupApi(t *testing.T) {
	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WillReturnError(sql.ErrNoRows)

		Convey("When client POST /v1/update-group/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/update-group/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"error",
	"Code":"ErrGroupNotFound",
	"Message":"Group is not found."
}
`)
			})
		})
	})

	Convey("Given some row in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "world_location_id", "name", "code", "editor_data"}).
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", 1, "Test group", "Code", "{}")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"updated_at"}).
			AddRow("2018-01-12 00:00:00")
		sqlMocker.ExpectQuery(`UPDATE\s+groups`).WithArgs("New group name", "Some new code", []byte(`{}`), 1).WillReturnRows(rows)

		Convey("When POST /v1/update-group/1", func() {
			res, err := http.Post(ts.URL+"/v1/update-group/1", "application/json", strings.NewReader(`{
	"Name" : "New group name",
	"Code" : "Some new code"
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
    "WorldLocationId": 1,
    "Name": "New group name",
	"Code": "Some new code",
	"EditorData": {}
  }
}
`)
			})
		})
	})
}
