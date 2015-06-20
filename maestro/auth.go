package backstage

import (
	"io/ioutil"
	"syscall"
)

type TokenInfo struct {
	Expires   int    `json:"expires"`
	CreatedAt string `bson:"created_at" json:"created_at"`
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
}

var TokenFileName = JoinHomePath(".backstage_token")

func WriteToken(token string) error {
	tokenFile, err := filesystem().OpenFile(TokenFileName, syscall.O_RDWR|syscall.O_CREAT|syscall.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer tokenFile.Close()
	tokenFile.WriteString(token)
	return nil
}

func ReadToken() (string, error) {
	tokenFile, err := filesystem().OpenFile(TokenFileName, syscall.O_RDWR, 0600)
	if err != nil {
		return "", ErrLoginRequired
	}
	defer tokenFile.Close()
	data, _ := ioutil.ReadAll(tokenFile)
	return string(data), nil
}

func DeleteToken() error {
	err := filesystem().Remove(TokenFileName)
	return err
}
