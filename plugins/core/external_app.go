package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) ExternalAppApi() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/external-app", []string{"context", "query:id"}, []string{"app"}, p.externalAppController.GetExternalApp),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/external-apps", []string{"context"}, []string{"apps"}, p.externalAppController.ListExternalApp),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/external-app", []string{"context", "body"}, []string{"app"}, p.externalAppController.CreateExternalApp),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/external-app", []string{"context", "query:id", "body"}, []string{"app"}, p.externalAppController.EditExternalApp),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/external-app", []string{"context", "query:id"}, nil, p.externalAppController.DeleteExternalApp),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/external-app/enable", []string{"context", "query:id"}, nil, p.externalAppController.EnableExternalApp),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/external-app/disable", []string{"context", "query:id"}, nil, p.externalAppController.DisableExternalApp),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/external-app/token", []string{"context", "query:id"}, []string{"token"}, p.externalAppController.UpdateExternalAppToken),
	}
}
