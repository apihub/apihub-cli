package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

type Service struct {
	Description   string   `json:"description,omitempty"`
	Disabled      bool     `json:"disabled,omitempty"`
	Documentation string   `json:"documentation,omitempty"`
	Endpoint      string   `json:"endpoint,omitempty"`
	Owner         string   `json:"owner,omitempty"`
	Subdomain     string   `json:"subdomain,omitempty"`
	Team          string   `json:"team,omitempty"`
	Timeout       int      `json:"timeout,omitempty"`
	Transformers  []string `json:"transformers,omitempty"`
	client        *HTTPClient
}

func (s *Service) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "team-service-add",
			Usage:       "team-service-add --team <team> --subdomain <subdomain> --endpoint <api_endpoint>",
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
				cli.StringFlag{
					Name:  "transformers, tf",
					Value: "",
					Usage: "Transformers",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-service-add")()
				disabled, err := strconv.ParseBool(c.String("disabled"))
				if err != nil {
					disabled = false
				}
				timeout, err := strconv.ParseInt(c.String("timeout"), 10, 0)
				if err != nil {
					timeout = 0
				}

				service := &Service{
					Subdomain:     c.String("subdomain"),
					Description:   c.String("description"),
					Disabled:      disabled,
					Documentation: c.String("documentation"),
					Endpoint:      c.String("endpoint"),
					Team:          c.String("team"),
					Timeout:       int(timeout),
					client:        NewHTTPClient(&http.Client{}),
				}
				result := service.create()
				fmt.Println(result)
			},
		},
		{
			Name:        "team-service-update",
			Usage:       "team-service-update --team <team> --subdomain <subdomain> --endpoint <api_endpoint>",
			Description: "Update an existing service.",
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
				cli.StringFlag{
					Name:  "transformers, tf",
					Value: "",
					Usage: "Transformers",
				},
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("team-service-update")()
				disabled, err := strconv.ParseBool(c.String("disabled"))
				if err != nil {
					disabled = false
				}
				timeout, err := strconv.ParseInt(c.String("timeout"), 10, 0)
				if err != nil {
					timeout = 0
				}

				service := &Service{
					Subdomain:     c.String("subdomain"),
					Disabled:      disabled,
					Description:   c.String("description"),
					Documentation: c.String("documentation"),
					Endpoint:      c.String("endpoint"),
					Team:          c.String("team"),
					Timeout:       int(timeout),
					client:        NewHTTPClient(&http.Client{}),
				}
				result := service.update()
				fmt.Println(result)
			},
		},
		{
			Name:        "team-service-remove",
			Usage:       "team-service-remove --subdomain <subdomain>",
			Description: "Remove an existing service.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "subdomain, s",
					Value: "",
					Usage: "Subdomain",
				},
				cli.StringFlag{
					Name:  "team, t",
					Value: "",
					Usage: "Team",
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
						Team:      c.String("team"),
						client:    NewHTTPClient(&http.Client{}),
					}
					result := service.remove()
					fmt.Println(result)
				}
			},
		},
	}
}

func (s *Service) create() string {
	path := fmt.Sprintf("/api/teams/%s/services", s.Team)
	service := &Service{}
	response, err := s.client.MakePost(path, s, service)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated {
		return "Your new service has been created."
	}
	panic("The team was not found.")
}

func (s *Service) update() string {
	path := fmt.Sprintf("/api/teams/%s/services/%s", s.Team, s.Subdomain)
	service := &Service{}
	response, err := s.client.MakePut(path, s, service)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		return "Your service has been updated."
	}
	panic("The team was not found.")
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
	panic("The service was not found for the team provided.")
}
