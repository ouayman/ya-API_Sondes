package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RestClient struct {
	authentication authentication
}

func (obj *RestClient) SetBasicAuth(username string, password string) {
	obj.authentication = &basic{username: username, password: password}
}

func (obj *RestClient) SetBearerAuth(jsonWebToken string) {
	obj.authentication = &bearer{jsonWebToken: jsonWebToken}
}

func (obj *RestClient) SetHeaderKey(key string, value string) {
	obj.authentication = &headerKey{key: key, value: value}
}

// Retourne le contenu, le code d'erreur ainsi qu'une erreur majeure éventuellement (Voir à mettre le status code dans l'erreur? définition d'une classe d'erreur )
func (obj RestClient) Request(url string) ([]byte, int, error) {
	// Création de la requête
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	// Ajout de l'authentification à la requête
	if obj.authentication != nil {
		obj.authentication.SetAuth(request)
	}

	// Exécution de la requête
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	// return binary de la réponse
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return body, response.StatusCode, err
}

func (obj RestClient) AdvRequest(method string, url string, object interface{}) ([]byte, int, error) {
	valueJSON, err := json.Marshal(object)
	if err != nil {
		return nil, 0, err
	}

	// Création de la requête
	request, err := http.NewRequest(method, url, bytes.NewBuffer(valueJSON))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-Type", "application/json")

	// Ajout de l'authentification à la requête
	if obj.authentication != nil {
		obj.authentication.SetAuth(request)
	}

	// Exécution de la requête
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	// return binary de la réponse
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return body, response.StatusCode, err
}

func (obj RestClient) PostRequest(url string, value interface{}, out interface{}) (int, error) {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		// DEBUG println(url, " error Marshal")
		return 0, err
	}

	// Exécution de la requête
	response, err := http.Post(url, "application/json", bytes.NewBuffer(valueJSON))
	if err != nil {
		// DEBUG println(url, " error Post")
		if response != nil {
			return response.StatusCode, err
		}
		return -1, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		var errorMsg = fmt.Sprintln("Http error status code: ", response.StatusCode)
		if body, err := ioutil.ReadAll(response.Body); err == nil {
			errorMsg = string(body)
		}
		// DEBUG println(url, " error ", errorMsg)

		return response.StatusCode, errors.New(errorMsg)
	}

	// return binary de la réponse
	/*
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			println("ReadAll")
			return response.StatusCode, err
		}

		err = json.Unmarshal(body, &out)
		if err != nil {
			println("Unmarshal")
			println(string(body))
			return response.StatusCode, err
		}

		return response.StatusCode, err
	*/
	return response.StatusCode, json.NewDecoder(response.Body).Decode(out)
}
