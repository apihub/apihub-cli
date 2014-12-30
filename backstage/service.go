package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

type Service struct {
	AllowKeylessUse bool   `json:"allow_keyless_use"`
	Description     string `json:"description"`
	Disabled        bool   `json:"disabled"`
	Documentation   string `json:"documentation"`
	Endpoint        string `json:"endpoint"`
	Owner           string `json:"owner"`
	Subdomain       string `json:"subdomain"`
	Team            string `json:"team"`
	Timeout         int    `json:"timeout"`
	client          *HTTPClient
}

func (s *Service) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "team-service-add",
			Usage:       "team-service-add --team <team> --subdomain <subdomain> --endpoint <api_endpoint>\n   Your new service has been created.",
			Description: "Create a new service.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "description, desc",
					Value: "",
					Usage: "Service description",
				},
				cli.StringFlag{
					Name:  "disabled, dis",
					Value: "",
					Usage: "Disable the service",
				},
				cli.StringFlag{
					Name:  "documentation, doc",
					Value: "",
					Usage: "Url with the documentation",
				},
				cli.StringFlag{
					Name:  "endpoint, e",
					Value: "",
					Usage: "Url where the service is running",
				},
				cli.StringFlag{
					Name:  "keyless, k",
					Value: "",
					Usage: "Allow keyless use",
				},
				cli.StringFlag{
					Name:  "subdomain, s",
					Value: "",
					Usage: "Desired subdomain",
				},
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Team responsible for the service",
				},
				cli.StringFlag{
					Name:  "timeout",
					Value: "",
					Usage: "Timeout",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-service-add")()
				keyless, err := strconv.ParseBool(c.String("keyless"))
				if err != nil {
					keyless = false
				}
				disabled, err := strconv.ParseBool(c.String("disabled"))
				if err != nil {
					disabled = false
				}
				timeout, err := strconv.ParseInt(c.String("timeout"), 10, 0)
				if err != nil {
					timeout = 0
				}

				service := &Service{
					Subdomain:       c.String("subdomain"),
					AllowKeylessUse: keyless,
					Description:     c.String("description"),
					Disabled:        disabled,
					Documentation:   c.String("documentation"),
					Endpoint:        c.String("endpoint"),
					Team:            c.String("team"),
					Timeout:         int(timeout),
					client:          NewHTTPClient(&http.Client{}),
				}
				result := service.save()
				fmt.Println(result)
			},
		},
		{
			Name:        "team-service-remove",
			Usage:       "team-service-remove --subdomain <subdomain>\n   The service `<subdomain>` has been deleted.",
			Description: "Remove an existing service.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "subdomain, s",
					Value: "",
					Usage: "Subdomain",
				},
			},
			Action: func(c *cli.Context) {
				context := &Context{Stdout: os.Stdout, Stdin: os.Stdin}
				if Confirm(context, "Are you sure you want to delete this service? This action cannot be undone.") != true {
					fmt.Println(ErrCommandCancelled)
				} else {
					defer RecoverStrategy("team-service-remove")()
					service := &Service{
						Subdomain: c.String("subdomain"),
						client:    NewHTTPClient(&http.Client{}),
					}
					result := service.remove()
					fmt.Println(result)
				}
			},
		},
	}
}

func (s *Service) save() string {
	path := fmt.Sprintf("/api/teams/%s/services", s.Team)
	service := &Service{}
	response, err := s.client.MakePost(path, s, service)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated {
		return "Your new service has been created."
	}
	return ErrBadRequest.Error()
}

func (s *Service) remove() string {
	path := fmt.Sprintf("/api/teams/%s/services/%s", s.Team, s.Subdomain)
	service := &Service{}
	response, err := s.client.MakeDelete(path, nil, service)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		return "The service `" + s.Subdomain + "` has been deleted."
	}
	return ErrBadRequest.Error()
}
