package rest

import (
	"github.com/yossefaz/go-http-client/gohttp_mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gohttp_mock.MockupServer.Start()
	os.Exit(m.Run())
}

func TestLoginUserTimeOutFromApi(t *testing.T) {
	// Delete all mocks in every new test case to ensure a clean environment:
	gohttp_mock.MockupServer.DeleteMocks()

	// Configure a new mock:
	gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
		Method:             http.MethodPost,
		Url:                "https://api.micro-gis.com/users/login",
		RequestBody:        `{"email":"email@gmail.com","password":"password"}`,
		ResponseStatusCode: -1,
		ResponseBody:       "{}",
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restClient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	// Delete all mocks in every new test case to ensure a clean environment:
	gohttp_mock.MockupServer.DeleteMocks()

	// Configure a new mock:
	gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
		Method:             http.MethodPost,
		Url:                "https://api.micro-gis.com/users/login",
		RequestBody:        `{"email":"email@gmail.com","password":"password"}`,
		ResponseStatusCode: http.StatusNotFound,
		ResponseBody:       `{"message" : "invalid login credentials", "status" : "404", "error" : "not found"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)

}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	// Delete all mocks in every new test case to ensure a clean environment:
	gohttp_mock.MockupServer.DeleteMocks()

	// Configure a new mock:
	gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
		Method:             http.MethodPost,
		Url:                "https://api.micro-gis.com/users/login",
		RequestBody:        `{"email":"email@gmail.com","password":"password"}`,
		ResponseStatusCode: http.StatusNotFound,
		ResponseBody:       `{"message" : "invalid login credentials", "status" : 404, "error" : "not found"}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	// Delete all mocks in every new test case to ensure a clean environment:
	gohttp_mock.MockupServer.DeleteMocks()

	// Configure a new mock:
	gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
		Method:             http.MethodPost,
		Url:                "https://api.micro-gis.com/users/login",
		RequestBody:        `{"email":"email@gmail.com","password":"password"}`,
		ResponseStatusCode: http.StatusOK,
		ResponseBody: `{
			"id": "1",
			"first_name": "Regis",
			"last_name": "Azoulay",
			"email": "test122@test",
			"date_created": "2020-09-26 19:43:01",
			"status": "active"
		}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshall user response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	// Delete all mocks in every new test case to ensure a clean environment:
	gohttp_mock.MockupServer.DeleteMocks()

	// Configure a new mock:
	gohttp_mock.MockupServer.AddMock(gohttp_mock.Mock{
		Method:             http.MethodPost,
		Url:                "https://api.micro-gis.com/users/login",
		RequestBody:        `{"email":"email@gmail.com","password":"password"}`,
		ResponseStatusCode: http.StatusOK,
		ResponseBody: `{
			"id": 1,
			"first_name": "Regis",
			"last_name": "Azoulay",
			"email": "test122@test",
			"date_created": "2020-09-26 19:43:01",
			"status": "active"
		}`,
	})

	repository := userRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.EqualValues(t, "Regis", user.FirstName)
	assert.EqualValues(t, "Azoulay", user.LastName)
	assert.EqualValues(t, "test122@test", user.Email)
}
