package v1

import (
	"fmt"
	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/gorilla/mux"
	"github.com/smartystreets/assertions"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetTestHelpers(t *testing.T) (ts *httptest.Server, sqlMocker sqlmock.Sqlmock) {
	db, sqlMocker, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlMocker.MatchExpectationsInOrder(false)

	fakeServer := server.NewServer("amqp://localhost:5672/%2f")
	fakeServer.Start()

	conn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		t.Fatalf(err.Error())
	}

	r := mux.NewRouter()
	New(r, conn, db)

	ts = httptest.NewServer(r)

	return ts, sqlMocker
}

func ShouldBeValidSuccessJsonResponse(actual interface{}, expected ...interface{}) string {
	return ShouldBeValidJsonResponse(http.StatusOK, actual, expected...)
}

func ShouldBeValidInternalErrorJsonResponse(actual interface{}, expected ...interface{}) string {
	return ShouldBeValidJsonResponse(http.StatusInternalServerError, actual, expected...)
}

func ShouldBeValidNotFoundJsonResponse(actual interface{}, expected ...interface{}) string {
	return ShouldBeValidJsonResponse(http.StatusNotFound, actual, expected...)
}

func ShouldBeValidBadRequestJsonResponse(actual interface{}, expected ...interface{}) string {
	return ShouldBeValidJsonResponse(http.StatusBadRequest, actual, expected...)
}

func ShouldBeValidJsonResponse(expectedStatus int, actual interface{}, expected ...interface{}) string {
	response := actual.(*http.Response)
	if response == nil {
		return "Actual value should be http.Response."
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != expectedStatus {
		return fmt.Sprintf("Expected status code: %d. Actual: %d. Body: %s", expectedStatus, response.StatusCode, string(bodyBytes))
	}

	if err != nil {
		return err.Error()
	}

	return assertions.ShouldEqualJSON(string(bodyBytes), expected[0])
}

func ShouldBeValidNoContentResponse(actual interface{}, expected ...interface{}) string {
	response := actual.(*http.Response)
	if response == nil {
		return "Actual value should be http.Response."
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Sprintf("Expected status code: %d. Actual: %d. Body: %s", http.StatusNoContent, response.StatusCode, string(bodyBytes))
	}

	if err != nil {
		return err.Error()
	}

	return assertions.ShouldBeEmpty(string(bodyBytes))
}
