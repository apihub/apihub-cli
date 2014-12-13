package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
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

func sortMapKeys(m map[string]string) []string {
	mk := make([]string, len(m))
	i := 0
	for k, _ := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)
	return mk
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
