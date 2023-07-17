package file

import (
	"net/url"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Proxy   string            `yaml:"proxy"`
	Timeout int               `yaml:"timeout"`
}

type File struct {
	Url      string `yaml:"url"`
	Filename string `yaml:"filename"`
	Config   Config `yaml:"config"`
}

type Yaml struct {
	Config Config `yaml:"config"`
	Files  []File `yaml:"files"`
}

func Unmarshal(in []byte) (Yaml, error) {
	var fileYaml Yaml
	err := yaml.Unmarshal(in, &fileYaml)
	if err == nil {
		for index := range fileYaml.Files {
			if fileYaml.Files[index].Filename == "" {
				u, _ := url.Parse(fileYaml.Files[index].Url)
				fileYaml.Files[index].Filename = path.Base(u.Path)
			}
		}
	}
	return fileYaml, err
}
