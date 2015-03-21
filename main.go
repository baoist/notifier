package main

import (
	"path/filepath"

	"github.com/baoist/notifier/settings"
	"github.com/baoist/notifier/watch"
)

func main() {
	filepath, err := filepath.Abs("./webhooks.yml")
	if err != nil {
		panic(err)
	}

	webhooks := settings.Import(filepath)
	for _, webhook := range webhooks {
		watch.Connect(webhook.Service, webhook)
	}
}
