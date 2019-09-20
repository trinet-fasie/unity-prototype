package v1

import (
	"database/sql"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"testing"
)

func TestGetWorldsApi(t *testing.T) {
	Convey("Given error in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+worlds`).WillReturnError(errors.New("some database error"))

		Convey("When client GET /v1/worlds", func() {
			res, err := http.Get(ts.URL + "/v1/worlds")

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"error",
	"Code":"",
	"Message":"some database error"
}
`)
			})
		})
	})

	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+worlds`).WillReturnError(sql.ErrNoRows)

		Convey("When GET /v1/worlds", func() {
			res, err := http.Get(ts.URL + "/v1/worlds")

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
		rows := sqlmock.NewRows([]string{"id", "name", "world_configurations_count", "created_at", "updated_at"}).
			AddRow(1, "World 1", 0, "2018-01-01 00:00:00", "2018-01-02 00:00:00").
			AddRow(2, "World 2", 1, "2018-01-01 00:00:00", "2018-01-02 00:00:00").
			AddRow(3, "World 3", 2, "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+worlds`).WillReturnRows(rows)

		Convey("When GET /v1/worlds", func() {
			res, err := http.Get(ts.URL + "/v1/worlds")

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
      "Name": "World 1",
      "Configurations": 0,
      "CreatedAt": "2018-01-01 00:00:00",
      "UpdatedAt": "2018-01-02 00:00:00"
    },
    {
      "Id": 2,
      "Name": "World 2",
      "Configurations": 1,
      "CreatedAt": "2018-01-01 00:00:00",
      "UpdatedAt": "2018-01-02 00:00:00"
    },
    {
      "Id": 3,
      "Name": "World 3",
      "Configurations": 2,
      "CreatedAt": "2018-01-01 00:00:00",
      "UpdatedAt": "2018-01-02 00:00:00"
    }
  ]
}
`)
			})
		})
	})
}
