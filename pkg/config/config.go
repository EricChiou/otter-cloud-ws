package config

import (
	"io/ioutil"
	"strings"
)

// LoadConfigFile load config file
func LoadConfigFile(configFilePath string) (map[string]string, error) {
	config := make(map[string]string)

	lines, err := fileHandler(configFilePath)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if len(line) > 0 && line[0:1] != "#" {
			val := strings.Split(line, "=")
			if len(val) == 2 {
				config[strings.TrimSpace(val[0])] = strings.TrimSpace(val[1])
			}
		}
	}

	return config, nil
}

func fileHandler(filePath string) ([]string, error) {
	var lines []string
	binary, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	str := string(binary)
	str = strings.Replace(str, "\r", "", -1)
	lines = strings.Split(str, "\n")
	return lines, nil
}
