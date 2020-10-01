package rest

import (
	"encoding/json"
	"github.com/micro-gis/oauth-api/src/domain/users"
	errors "github.com/micro-gis/oauth-api/src/utils/errors_util"
	"github.com/yossefaz/go-http-client/gohttp"
	"time"
)

var (
	Timeout         = 100 * time.Millisecond
	BaseURL         = "http://127.0.0.1:8086"
	usersRestClient gohttp.Client
)

func init() {
	usersRestClient = gohttp.NewBuilder().Build()
}

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type userRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &userRepository{}
}

func (r *userRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response, _ := usersRestClient.Post(BaseURL+"/users/login", request)
	if response == nil || response.StatusCode < 100 {
		return nil, errors.NewInternalServerError("invalid restClient response when trying to login user")
	}

	if response.StatusCode > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshall user response")
	}
	return &user, nil
}
