package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/url"
	"testing"
)

func TestDeleteLocationApi(t *testing.T) {
	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WillReturnError(sql.ErrNoRows)

		Convey("When client POST /v1/delete-location/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-location/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNotFoundJsonResponse, `
{
	"Status":"error",
	"Code":"ErrLocationNotFound",
	"Message":"Location is not found."
}
`)
			})
		})
	})

	Convey("Given some used location in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "guid", "name", "usages", "created_at", "updated_at"}).
			AddRow(1, "c949de70-7e44-42b7-b34e-b84efd1afbf1", "Example location", 1, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WillReturnRows(rows)

		Convey("When POST /v1/delete-location/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-location/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
  "Status": "error",
  "Code":"ErrDeleteUsedLocation",
  "Message":"Cannot delete location. Location is used."
}
`)
			})
		})
	})

	Convey("Given some unused location in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "guid", "name", "usages", "created_at", "updated_at"}).
			AddRow(1, "c949de70-7e44-42b7-b34e-b84efd1afbf1", "Example location", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WillReturnRows(rows)
		sqlMocker.ExpectExec(`DELETE\s+FROM\s+locations`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		Convey("When POST /v1/delete-location/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-location/1", url.Values{})

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
