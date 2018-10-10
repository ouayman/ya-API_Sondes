package helper

import "net/http"

// TODO Voir à gérer les défférentes autorisations basic, bearer, autres
type authentication interface {
	SetAuth(req *http.Request)
}

type basic struct {
	username string
	password string
}

func (obj basic) SetAuth(req *http.Request) {
	req.SetBasicAuth(obj.username, obj.password)
}

type bearer struct {
	jsonWebToken string
}

func (obj bearer) SetAuth(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+obj.jsonWebToken)
}

type headerKey struct {
	key   string
	value string
}

func (obj headerKey) SetAuth(req *http.Request) {
	req.Header.Set(obj.key, obj.value)
}
