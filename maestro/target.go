package apihub

import (
	"io/ioutil"
	"syscall"

	"github.com/tsuru/tsuru/fs"
	"gopkg.in/v1/yaml"
)

var (
	Fsystem        fs.Fs
	TargetFileName = JoinHomePath(".apihub_targets")
)

func filesystem() fs.Fs {
	if Fsystem == nil {
		Fsystem = fs.OsFs{}
	}
	return Fsystem
}

type Target struct {
	Current string
	Options map[string]string
}

func LoadTargets() (*Target, error) {
	targetsFile, err := filesystem().OpenFile(TargetFileName, syscall.O_RDWR|syscall.O_CREAT, 0600)
	defer targetsFile.Close()
	if err != nil {
		return nil, err
	}

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

func (t *Target) Add(label string, endpoint string) error {
	if _, ok := t.Options[label]; ok {
		return ErrLabelExists
	}

	t.Options[label] = endpoint
	return t.Save()
}

func (t *Target) Remove(label string) error {
	if _, ok := t.Options[label]; !ok {
		return ErrTargetNotFound
	}

	if t.Current == label {
		t.Current = ""
	}

	delete(t.Options, label)
	return t.Save()
}

func (t *Target) SetDefault(label string) error {
	if _, ok := t.Options[label]; !ok {
		return ErrTargetNotFound
	}

	t.Current = label
	return t.Save()
}

func (t *Target) Save() error {
	d, err := yaml.Marshal(&t)
	if err != nil {
		return err
	}

	targetsFile, err := filesystem().OpenFile(TargetFileName, syscall.O_RDWR|syscall.O_CREAT|syscall.O_TRUNC, 0600)
	defer targetsFile.Close()
	if err != nil {
		return err
	}

	n, err := targetsFile.WriteString(string(d))
	if n != len(string(d)) || err != nil {
		return ErrFailedWritingTargetFile
	}
	return nil
}

func (t *Target) GetOptions() (string, []string, map[string]string) {
	current := t.Current
	sortedMapKeys := SortMapKeys(t.Options)

	return current, sortedMapKeys, t.Options
}

func GetCurrentTarget() (string, error) {
	t, err := LoadTargets()
	if err != nil {
		return "", err
	}

	current := t.Options[t.Current]
	if current == "" {
		return "", ErrEndpointNotFound
	}

	return current, nil
}
