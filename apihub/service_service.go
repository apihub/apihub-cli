package apihub

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ServiceService struct {
	client HTTPClient
}

func NewServiceService(client HTTPClient) *ServiceService {
	return &ServiceService{
		client: client,
	}
}

func (s ServiceService) Create(subdomain string, disabled bool, description string, documentation string, endpoint string, team string, timeout int64, transformers []string) (Service, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusCreated,
		Method:         "POST",
		Path:           "/api/services",
		Body: Service{
			Subdomain:     subdomain,
			Disabled:      disabled,
			Description:   description,
			Documentation: documentation,
			Endpoint:      endpoint,
			Team:          team,
			Timeout:       timeout,
			Transformers:  transformers,
		},
	})

	if err != nil {
		return Service{}, err
	}

	var service Service
	err = json.Unmarshal(body, &service)
	if err != nil {
		panic(err)
	}

	return service, nil
}

func (s ServiceService) Update(subdomain string, disabled bool, description string, documentation string, endpoint string, team string, timeout int64, transformers []string) (Service, error) {
	body, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "PUT",
		Path:           fmt.Sprintf("/api/services/%s", subdomain),
		Body: Service{
			Subdomain:     subdomain,
			Disabled:      disabled,
			Description:   description,
			Documentation: documentation,
			Endpoint:      endpoint,
			Team:          team,
			Timeout:       timeout,
			Transformers:  transformers,
		},
	})

	if err != nil {
		return Service{}, err
	}

	var service Service
	err = json.Unmarshal(body, &service)
	if err != nil {
		panic(err)
	}

	return service, nil
}

func (s ServiceService) Delete(subdomain string, team string) error {
	_, err := s.client.MakeRequest(RequestArgs{
		AcceptableCode: http.StatusOK,
		Method:         "DELETE",
		Path:           fmt.Sprintf("/api/services/%s", subdomain),
	})

	return err
}
