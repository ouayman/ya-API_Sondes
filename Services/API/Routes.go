package apiService

import (
	"../../Helper/Http"
)

func routes() helperhttp.Routes {
	handler := newHandler()

	return helperhttp.Routes{
		helperhttp.Route{
			Name:    "GetInfo",
			Method:  "GET",
			Pattern: "/",
			Handler: helperhttp.ErrorFnHandler(handler.getInfo),
		},
		/*helperhttp.Route{
			Name:    "GetSPIS",
			Method:  "GET",
			Pattern: "/SPIS",
			Handler: helperhttp.ErrorFnHandler(handler.getSPIS),
		},*/
		helperhttp.Route{
			Name:    "GetEPlanning",
			Method:  "GET",
			Pattern: "/ePlanning",
			Handler: helperhttp.ErrorFnHandler(handler.getEPlanning),
		},
		helperhttp.Route{
			Name:    "GetAllOLT",
			Method:  "GET",
			Pattern: "/OLT",
			Handler: helperhttp.ErrorFnHandler(handler.getAllOLT),
		},
		helperhttp.Route{
			Name:    "GetSingleOLT",
			Method:  "GET",
			Pattern: "/OLT/{webservice}",
			Handler: helperhttp.ErrorFnHandler(handler.getSingleOLT),
		},
		/*		helperhttp.Route{
				Name:    "GetHistoriqueAcompte",
				Method:  "GET",
				Pattern: "/historiqueAcompte/magasin/{magasin}/commande/{numeroCommande}",
				Handler: helperhttp.ErrorFnHandler(handler.getHistoriqueAcompte),
			},*/
	}
}
