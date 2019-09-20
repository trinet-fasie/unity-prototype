package v1

import (
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"strings"
	"testing"
)

func TestUpdateLocationTagsApi(t *testing.T) {
	Convey("Given invalid request", t, func() {
		ts, sqlMocker := GetTestHelpers(t)

		Convey("When client POST /v1/update-location-tags/1 with empty body", func() {
			sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WillReturnError(sql.ErrNoRows)

			res, err := http.Post(ts.URL+"/v1/update-location-tags/1", "application/json", strings.NewReader(``))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidBadRequestJsonResponse, `
{
  "Code": "400",
  "Message": "unexpected end of JSON input",
  "Status": "error"
}
`)
			})
		})

		Convey("When client POST /v1/update-location-tags/1 without location in database", func() {
			sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WillReturnError(sql.ErrNoRows)

			res, err := http.Post(ts.URL+"/v1/update-location-tags/1", "application/json", strings.NewReader(`[]`))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNotFoundJsonResponse, `
{
  "Code": "404",
  "Message": "Location is not found.",
  "Status": "error"
}
`)
			})
		})
	})

	Convey("Given valid request", t, func() {
		ts, sqlMocker := GetTestHelpers(t)

		Convey("When client POST /v1/update-location-tags/1 with empty body", func() {
			locationRows := sqlmock.NewRows([]string{"id", "guid", "name", "usages", "created_at", "updated_at"}).
				AddRow(
					1,                                      // id
					"c949de70-7e44-42b7-b34e-b84efd1afbf1", // guid
					"1",                                    // name
					0,                                      // usages
					"2018-01-01 00:00:00",                  // created_at
					"2018-01-02 00:00:00",                  // updated_at
				)
			sqlMocker.ExpectQuery(`SELECT.*FROM\s+locations`).WithArgs(1).WillReturnRows(locationRows)

			tagRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "text"}).
				AddRow(
					1,                     // id
					"2018-01-01 00:00:00", // created_at
					"2018-01-02 00:00:00", // updated_at
					"tag1",                // text
				).
				AddRow(
					2,                     // id
					"2018-01-01 00:00:00", // created_at
					"2018-01-02 00:00:00", // updated_at
					"tag2",                // text
				)
			sqlMocker.ExpectQuery(`SELECT.*FROM\s+location_tags`).WithArgs(1).WillReturnRows(tagRows)

			tagByIdRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "text"}).
				AddRow(
					3,                     // id
					"2018-01-01 00:00:00", // created_at
					"2018-01-02 00:00:00", // updated_at
					"tag3",                // text
				)
			sqlMocker.ExpectQuery(`SELECT.*FROM\s+location_tags`).WithArgs(3).WillReturnRows(tagByIdRows)
			// insert tag
			sqlMocker.ExpectExec(`INSERT\s+INTO\s+location_tag_locations`).WithArgs().WillReturnResult(sqlmock.NewResult(0, 1))
			// delete tag
			sqlMocker.ExpectExec(`DELETE\s+FROM\s+location_tag_locations`).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

			res, err := http.Post(ts.URL+"/v1/update-location-tags/1", "application/json", strings.NewReader(`
	[
		{"Id": 1, "Text": "tag1"},
		{"Id": 3, "Text": "tag3"}
	]
`))

			Convey("Request error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Response should be valid JSON", func() {
				So(res, ShouldBeValidNoContentResponse)
			})
		})

	})
}
