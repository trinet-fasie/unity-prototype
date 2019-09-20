package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/url"
	"testing"
)

func TestDeleteWorldLocationApi(t *testing.T) {
	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+world_locations`).WillReturnError(sql.ErrNoRows)

		Convey("When client POST /v1/delete-world-location/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-world-location/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNotFoundJsonResponse, `
{
	"Status":"error",
	"Code":"ErrWorldLocationNotFound",
	"Message":"World location is not found."
}
`)
			})
		})
	})

	Convey("Given some world location in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "world_id", "location_id", "name", "created_at", "updated_at"}).
			AddRow(1, 1, 1, "Test world location", "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+world_locations`).WillReturnRows(rows)
		sqlMocker.ExpectExec(`DELETE\s+FROM\s+world_locations`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		Convey("When POST /v1/delete-world-location/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-world-location/1", url.Values{})

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
