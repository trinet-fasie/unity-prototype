package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/url"
	"testing"
)

func TestDeleteWorldApi(t *testing.T) {
	Convey("Given no rows in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+worlds`).WillReturnError(sql.ErrNoRows)

		Convey("When client POST /v1/delete-world/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-world/1", url.Values{})

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidSuccessJsonResponse, `
{
	"Status":"error",
	"Code":"ErrWorldNotFound",
	"Message":"World is not found."
}
`)
			})
		})
	})

	Convey("Given some world in database", t, func() {
		ts, sqlMocker := GetTestHelpers(t)
		rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
			AddRow(1, "Example world", "2018-01-01 00:00:00", "2018-01-02 00:00:00")
		sqlMocker.ExpectQuery(`SELECT.*FROM\s+worlds`).WillReturnRows(rows)
		sqlMocker.ExpectExec(`DELETE\s+FROM\s+worlds`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		Convey("When POST /v1/delete-world/1", func() {
			res, err := http.PostForm(ts.URL+"/v1/delete-world/1", url.Values{})

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
