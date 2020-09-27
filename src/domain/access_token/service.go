package access_token

type Service interface {
	GetById(string)(*AccessToken, error)

}