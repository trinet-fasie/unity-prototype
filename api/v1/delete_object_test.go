package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/url"
	"testing"
)

func TestDeleteObjectApi(t *testing.T) {
	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+objects`).WillReturnError(sql.ErrNoRows)

		Convey("When client POST /v1/delete-object/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-object/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNotFoundJsonResponse, `
{
	"Status":"error",
	"Code":"ErrObjectNotFound",
	"Message":"Object not found."
}
`)
			})
		})
	})

	Convey("Given some used object in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "guid", "type_name", "name", "usages", "created_at", "updated_at"}).
			AddRow(1, "c949de70-7e44-42b7-b34e-b84efd1afbf1", "Trigger", "Button", 1, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+objects`).WillReturnRows(rows)

		Convey("When POST /v1/delete-object/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-object/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNotFoundJsonResponse, `
{
	"Status":"error",
	"Code":"ErrDeleteUsedObject",
	"Message":"Cannot delete object. Object is used."
}
`)
			})
		})
	})

	Convey("Given some unused object in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "guid", "type_name", "name", "usages", "created_at", "updated_at"}).
			AddRow(1, "c949de70-7e44-42b7-b34e-b84efd1afbf1", "Trigger", "Button", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+objects`).WillReturnRows(rows)
		sqlMocker.ExpectExec(`DELETE\s+FROM\s+objects`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		Convey("When POST /v1/delete-object/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-object/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"success",
	"Data":""
}
`)
			})
		})
	})
}
