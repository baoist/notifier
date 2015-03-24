package settings

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Webhook struct {
	Service  string
	Secret   string
	Token    string
	Port     int
	Watchers Channels
	Public   Channels
	Author   Author
}

type Webhooks []Webhook

type Author struct {
	Name string
	Link string
	Icon string
}

func Import(path string) Webhooks {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var webhooks Webhooks

	err = yaml.Unmarshal(file, &webhooks)
	if err != nil {
		panic(err)
	}

	return webhooks
}
