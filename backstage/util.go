package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func joinHomePath(p ...string) string {
	ps := []string{os.ExpandEnv("$HOME")}
	ps = append(ps, p...)
	return path.Join(ps...)
}

func parseBody(body io.ReadCloser, r interface{}) error {
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &r); err != nil {
		return err
	}
	return nil
}

func LoginRequired() error {
	_, err := ReadToken()
	if err != nil {
		err = ErrLoginRequired
		fmt.Println(err.Error())
		return err
	}
	return nil
}
