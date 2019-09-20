package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/url"
	"testing"
)

func TestDeleteGroupApi(t *testing.T) {
	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WillReturnError(sql.ErrNoRows)

		Convey("When client POST /v1/delete-group/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-group/1", url.Values{})

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

	Convey("Given some group in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "world_location_id", "name", "code", "editor_data"}).
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", 1, "Test group", "Code", "{}")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+groups`).WillReturnRows(rows)
		sqlMocker.ExpectExec(`DELETE\s+FROM\s+groups`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		Convey("When POST /v1/delete-group/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-group/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
  "Status": "success",
  "Data": ""
}
`)
			})
		})
	})
}
