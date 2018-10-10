package helperhttp

import (
	"encoding/json"
	"log"
	"net/http"
)

type Error struct {
	Err  error
	Code int
}

func (obj *Error) Error() string {
	return obj.Err.Error()
}

type ErrorFn func(http.ResponseWriter, *http.Request) error

func ErrorFnHandler(functor ErrorFn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recover PANIC error ?
		// Mode TRACE timer function start
		err := functor(w, r)
		// Mode TRACE timer function end
		if err == nil {
			return
		}

		switch errType := err.(type) {
		case *Error:
			respondContent(w, errType.Code, map[string]string{"error": errType.Error()})
			log.Println(errType.Error())
		case error:
			respondContent(w, 500, map[string]string{"error": errType.Error()})
			log.Println(errType.Error())
		default:
			respondContent(w, 500, map[string]string{"error": http.StatusText(500)})
			log.Println("unknown error")
		}
	})
}

func ExtractObjectFromBody(r *http.Request, out interface{}) (err error) {
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(out)
	return
}

func RespondCheck(w http.ResponseWriter, payload interface{}, err error) error {
	if err != nil {
		return err
	}
	return respondContent(w, http.StatusOK, payload)
}

func Respond(w http.ResponseWriter, payload interface{}) error {
	return respondContent(w, http.StatusOK, payload)
}

func respondContent(w http.ResponseWriter, code int, payload interface{}) error {
	if payload == nil {
		w.WriteHeader(code)
		return nil
	}

	response, err := json.Marshal(payload)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(code)
		w.Write(response)
	}
	return err
}

/*
func buildJSONResponse(w http.ResponseWriter, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		buildJSONError(w, err)
	}
}

func buildJSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	if err != nil {
		println(err.Error())
		content := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(content)
	}
}
*/
