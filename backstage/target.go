package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"syscall"

	"github.com/tsuru/tsuru/fs"
	"gopkg.in/v1/yaml"
)

var (
	TargetFileName       = joinHomePath(".backstage_targets")
	ErrLabelExists       = errors.New("The label provided exists already.")
	ErrLabelNotFound     = errors.New("Label not found.")
	ErrBadFormattedFile  = errors.New("Bad formatted file. Please open an issue on or Github page: backstage/backstage")
	ErrCommandCancelled  = errors.New("Command Cancelled.")
	ErrFailedWritingFile = errors.New("Failed trying to write the target file.")
)

var fsystem fs.Fs

func filesystem() fs.Fs {
	if fsystem == nil {
		fsystem = fs.OsFs{}
	}
	return fsystem
}

type Target struct {
	Current string
	Options map[string]string
}

func (t *Target) add(label string, endpoint string) error {
	if _, ok := t.Options[label]; ok {
		return ErrLabelExists
	}
	t.Options[label] = endpoint
	return t.save()
}

func (t *Target) list() string {
	var targetList bytes.Buffer
	for label, endpoint := range t.Options {
		if t.Current == label {
			targetList.WriteString("* ")
		}
		targetList.WriteString(label + " - " + endpoint + "\n")
	}
	return targetList.String()
}

func (t *Target) remove(label string) error {
	if _, ok := t.Options[label]; !ok {
		return ErrLabelNotFound
	}
	if t.Current == label {
		t.Current = ""
	}
	delete(t.Options, label)
	return t.save()
}

func (t *Target) setDefault(label string) error {
	if _, ok := t.Options[label]; !ok {
		return ErrLabelNotFound
	}
	t.Current = label
	return t.save()
}

func (t *Target) save() error {
	d, err := yaml.Marshal(&t)
	if err != nil {
		return err
	}
	targetsFile, err := filesystem().OpenFile(TargetFileName, syscall.O_RDWR|syscall.O_CREAT, 0600)
	if err != nil {
		return err
	}
	n, err := targetsFile.WriteString(string(d))
	if n != len(string(d)) || err != nil {
		return ErrFailedWritingFile
	}
	return nil
}

func LoadTargets() (*Target, error) {
	targetsFile, err := filesystem().OpenFile(TargetFileName, syscall.O_RDWR|syscall.O_CREAT, 0600)
	if err != nil {
		return nil, err
	}
	defer targetsFile.Close()
	data, err := ioutil.ReadAll(targetsFile)
	if err == nil {
		var t Target
		err = yaml.Unmarshal([]byte(data), &t)
		if err != nil {
			return nil, ErrBadFormattedFile
		}
		if t.Options == nil {
			t.Options = map[string]string{}
		}
		return &t, nil
	}
	return nil, err
}
