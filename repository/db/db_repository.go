package db

import (
	"errors"
	"github.com/RajatWithGolang/Books_OAuth_API/cassandra"
	"github.com/RajatWithGolang/Books_OAuth_API/domain/access_token"
	"github.com/RajatWithGolang/Books_Utils_API/restErrors"
	"github.com/gocql/gocql"
	
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, restErrors.RestErr)
	Create(access_token.AccessToken) restErrors.RestErr
	UpdateExpirationTime(access_token.AccessToken) restErrors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, restErrors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, restErrors.NewNotFoundError("no access token found with given id")
		}
		return nil, restErrors.NewInternalServerError("error when trying to get current id", errors.New("database error"))
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) restErrors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserId,
		at.ClientId,
		at.Expires,
	).Exec(); err != nil {
		return restErrors.NewInternalServerError("error when trying to save access token in database", err)
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) restErrors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return restErrors.NewInternalServerError("error when trying to update current resource", errors.New("database error"))
	}
	return nil
}
