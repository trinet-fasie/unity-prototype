package v1

import (
	"database/sql"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"testing"
)

func TestGetObjectTagsApi(t *testing.T) {
	Convey("Given error in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+object_tags`).WillReturnError(errors.New("some database error"))

		Convey("When client GET /v1/object-tags", func() {
			res, err := http.Get(ts.URL + "/v1/object-tags")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidInternalErrorJsonResponse, `
{
	"Status":"error",
	"Code":"500",
	"Message":"some database error"
}
`)
			})
		})
	})

	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+object_tags`).WillReturnError(sql.ErrNoRows)

		Convey("When GET /v1/object-tags", func() {
			res, err := http.Get(ts.URL + "/v1/object-tags")

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
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "text"}).
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", "tag1").
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", "tag2").
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", "tag3")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+object_tags`).WithArgs("sometext", ObjectTagsLimit).WillReturnRows(rows)

		Convey("When GET /v1/object-tags?search=SomeText", func() {
			res, err := http.Get(ts.URL + "/v1/object-tags?search=SomeText")

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
      "Text": "tag1"
    },
    {
      "Id": 1,
      "Text": "tag2"
    },
    {
      "Id": 1,
      "Text": "tag3"
    }
  ]
}
`)
			})
		})
	})
}
