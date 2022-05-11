package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Path   string       `json:"path"`
	Filter SearchFilter `json:"filter"`
}

type SearchFilter struct {
	Postcode    string   `json:"postcode"`
	AM          uint     `json:"am"`
	PM          uint     `json:"pm"`
	RecipeNames []string `json:"recipe_names"`
}

func Load(filename string) (Config, error) {
	var conf Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(data, &conf)
	return conf, err
}
