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
	client          *Client
}

func (s *Service) GetCommands() []cli.Command {
	return []cli.Command{
		{
			Name:        "service-add",
			Usage:       "service-add ...",
			Description: "Creates a new service.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyless, k",
					Value: "",
					Usage: "Allow keyless use",
				},
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
			},
			Action: func(c *cli.Context) {
				defer RecoverStrategy("service-add")()
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
					Owner:           c.String("owner"),
					Team:            c.String("team"),
					Timeout:         int(timeout),
					client:          NewClient(&http.Client{}),
				}
				result := service.save()
				fmt.Println(result)
			},
		},
		{
			Name:        "service-remove",
			Usage:       "service-remove --subdomain <subdomain>",
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
				if Confirm(context, "Are you sure you want to remove this service?") != true {
					fmt.Println(ErrCommandCancelled)
				} else {
					defer RecoverStrategy("service-remove")()
					service := &Service{
						Subdomain: c.String("subdomain"),
						client:    NewClient(&http.Client{}),
					}
					result := service.remove()
					fmt.Println(result)
				}
			},
		},
	}
}

func (s *Service) save() string {
	path := "/api/services"
	service := &Service{}
	response, err := s.client.MakePost(path, s, service)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusCreated {
		return "Service created successfully."
	}
	return ErrBadRequest.Error()
}

func (s *Service) remove() string {
	path := "/api/services/" + s.Subdomain
	service := &Service{}
	response, err := s.client.MakeDelete(path, nil, service)
	if err != nil {
		return err.Error()
	}

	if response.StatusCode == http.StatusOK {
		return "Service removed successfully."
	}
	return ErrBadRequest.Error()
}
