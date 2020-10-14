package db

import (
	"github.com/gocql/gocql"
	"github.com/micro-gis/oauth-api/src/clients/cassandra"
	"github.com/micro-gis/oauth-api/src/domain/access_token"
	errors "github.com/micro-gis/utils/rest_errors"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires from access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens (access_token, user_id, client_id, expires) VALUES (?,?,?,?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
	querySelectALlUsersAccessTokens = "SELECT access_token from access_tokens WHERE user_id=?;"
	queryDeleteUserToken = "DELETE FROM access_tokens WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, errors.RestErr)
	Create(access_token.AccessToken) errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) errors.RestErr
	DeleteUserTokens(access_token.AccessToken) errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("no access token found with given id")
		}
		return nil, errors.NewInternalServerError(err.Error(), err)
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) errors.RestErr {
	session := cassandra.GetSession()
	if err := session.Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserId,
		at.ClientId,
		at.Expires,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error(), err)
	}
	return nil
}
func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) errors.RestErr {
	session := cassandra.GetSession()
	if err := session.Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error(), err)
	}
	return nil
}

func (r *dbRepository) DeleteUserTokens(at access_token.AccessToken) errors.RestErr {
	session := cassandra.GetSession()
	// list all tokens
	iter := session.Query(querySelectALlUsersAccessTokens, at.UserId).Iter()
	var token string
	for iter.Scan(&token) {
		// delete each retireved token
		if err := session.Query(queryDeleteUserToken, token).Exec(); err != nil {
			return errors.NewInternalServerError(err.Error(), err)
		}
	}
	return nil
}


