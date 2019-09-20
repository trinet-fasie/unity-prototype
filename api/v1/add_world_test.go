package v1

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"strings"
	"testing"
)

func TestAddWorldApi(t *testing.T) {
	Convey("Given invalid request", t, func() {
		ts, _ := GetTestHelpers(t)

		Convey("When client POST /v1/add-world", func() {
			res, err := http.Post(ts.URL+"/v1/add-world", "application/json", strings.NewReader(`{}`))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidBadRequestJsonResponse, `
{
  "Status": "fail",
  "Data": {
    "World.Name": {
      "ActualTag": "required",
      "Field": "Name",
      "FieldNamespace": "World.Name",
      "Kind": 24,
      "Name": "Name",
      "NameNamespace": "Name",
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
	})

	Convey("Given valid request", t, func() {
		ts, sqlMocker := GetTestHelpers(t)

		Convey("When client POST /v1/add-world", func() {
			rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(1, "2018-01-01 00:00:00", "2018-01-02 00:00:00")

			sqlMocker.ExpectQuery(`INSERT\s+INTO\s+worlds`).WithArgs("Example world").WillReturnRows(rows)

			res, err := http.Post(ts.URL+"/v1/add-world", "application/json", strings.NewReader(`{
	"Name" : "Example world"
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
    "Name": "Example world",
	"Configurations":0
  }
}
`)
			})
		})
	})
}
