package apiService

import (
	"net/http"

	"../../Configuration"
	"../../Helper/Http"
	"github.com/gorilla/mux"
)

func newHandler() *handler {
	return &handler{}
}

type handler struct {
	Brick brick
}

type brick struct {
	Name string    `json:"name"`
	Elts []element `json:"elements"`
}

type element struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	State bool   `json:"state"`
	Url   string `json:"url"`
}

func (obj handler) getInfo(w http.ResponseWriter, r *http.Request) error {
	return helperhttp.Respond(w, configuration.Get())
}

/*
func (obj *handler) getSPIS(w http.ResponseWriter, r *http.Request) error {
}
*/

func (obj *handler) getEPlanning(w http.ResponseWriter, r *http.Request) error {
	obj.Brick = brick{Name: "ePlanning"} // Initializing the brick

	obj.Brick.Elts = obj.Brick.Elts[:0] // Clearing brick's elements
	obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "ePlanning", Type: "webservice", Url: "http://yrfrlmdelivery.int.adeo.com/delivery-rest/api/#!/Delivery_Computing_API/computePriceServiceLevels"})

	_, err := http.Get(obj.Brick.Elts[0].Url)
	if err != nil {
		obj.Brick.Elts[0].State = false
	} else {
		obj.Brick.Elts[0].State = true
	}

	return helperhttp.Respond(w, obj.Brick)
}

func (obj *handler) getAllOLT(w http.ResponseWriter, r *http.Request) error {
	obj.Brick = brick{Name: "OLT"} // Initializing the brick

	obj.Brick.Elts = obj.Brick.Elts[:0] // Clearing brick's elements
	obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "deposit", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/deposits/v1?wsdl"})
	obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "transaction", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/processTransaction/v8?wsdl"})
	obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "invoice", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/processInvoice/v1?wsdl"})
	obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "identification", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/ticket/v1?wsdl"})

	for i, _ := range obj.Brick.Elts {
		_, err := http.Get(obj.Brick.Elts[i].Url)
		if err != nil {
			obj.Brick.Elts[i].State = false
		} else {
			obj.Brick.Elts[i].State = true
		}
	}

	return helperhttp.Respond(w, obj.Brick)
}

func (obj *handler) getSingleOLT(w http.ResponseWriter, r *http.Request) error {
	obj.Brick = brick{Name: "OLT"} // Initializing the brick

	obj.Brick.Elts = obj.Brick.Elts[:0] // Clearing brick's elements
	switch mux.Vars(r)["webservice"] {
	case "deposit":
		obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "deposit", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/deposits/v1?wsdl"})
	case "transaction":
		obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "transaction", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/processTransaction/v8?wsdl"})
	case "invoice":
		obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "invoice", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/processInvoice/v1?wsdl"})
	case "identification":
		obj.Brick.Elts = append(obj.Brick.Elts, element{Name: "identification", Type: "webservice", Url: "http://rlmfrasbdc01.corp.leroymerlin.com:24080/olt-web/elements/ticket/v1?wsdl"})
	default:
		return helperhttp.Respond(w, "Aucun webservice avec ce nom !")
	}

	_, err := http.Get(obj.Brick.Elts[0].Url)
	if err != nil {
		obj.Brick.Elts[0].State = false
	} else {
		obj.Brick.Elts[0].State = true
	}

	return helperhttp.Respond(w, obj.Brick)
}
