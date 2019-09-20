package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"strings"
	"testing"
)

func TestAddObjectTagApi(t *testing.T) {
	Convey("Given invalid request", t, func() {
		ts, _ := GetTestHelpers(t)

		Convey("When client POST /v1/add-object-tag", func() {
			res, err := http.Post(ts.URL+"/v1/add-object-tag", "application/json", strings.NewReader(`{}`))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidBadRequestJsonResponse, `
{
  "Status": "fail",
  "Data": {
    "AddObjectTagDto.Text": {
      "ActualTag": "required",
      "Field": "Text",
      "FieldNamespace": "AddObjectTagDto.Text",
      "Kind": 24,
      "Name": "Text",
      "NameNamespace": "Text",
      "Param": "",
      "Tag": "required",
      "Type": {},
      "Value": ""
    }
  }
}
`)
			})
		})

		Convey("When client POST /v1/add-object-tag with too long tag text", func() {
			res, err := http.Post(ts.URL+"/v1/add-object-tag", "application/json", strings.NewReader(`{
	"Text": "Too long tag text XXXXXXXX"
}`))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidBadRequestJsonResponse, `
{
  "Status": "fail",
  "Data": {
    "AddObjectTagDto.Text": {
      "ActualTag": "max",
      "Field": "Text",
      "FieldNamespace": "AddObjectTagDto.Text",
      "Kind": 24,
      "Name": "Text",
      "NameNamespace": "Text",
      "Param": "25",
      "Tag": "max",
      "Type": {},
      "Value": "too long tag text xxxxxxxx"
    }
  }
}
`)
			})
		})
	})

	Convey("Given valid request with non-existent text in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+object_tags`).WithArgs("some valid tag").WillReturnError(sql.ErrNoRows)

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`INSERT\s+INTO\s+object_tags`).WithArgs("some valid tag").WillReturnRows(rows)

		Convey("When client POST /v1/add-object-tag", func() {
			res, err := http.Post(ts.URL+"/v1/add-object-tag", "application/json", strings.NewReader(`{
	"Text": "Some valid tag"
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
    "UpdatedAt": "2018-01-02 00:00:00",
    "Text": "some valid tag"
  }
}
`)
			})
		})
	})

	Convey("Given valid request with existent text in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "text"}).
			AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00", "some valid tag")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+object_tags`).WithArgs("some valid tag").WillReturnRows(rows)

		Convey("When client POST /v1/add-object-tag", func() {
			res, err := http.Post(ts.URL+"/v1/add-object-tag", "application/json", strings.NewReader(`{
	"Text": "Some valid tag"
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
    "UpdatedAt": "2018-01-02 00:00:00",
    "Text": "some valid tag"
  }
}
`)
			})
		})
	})
}
