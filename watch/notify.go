package watch

import (
	"errors"
	"fmt"

	"github.com/baoist/notifier/settings"
)

type Service interface {
	Listen()
}

func Connect(serviceName string, webhook settings.Webhook) {
	service, err := NewService(serviceName, webhook)
	if err != nil {
		panic(err)
	}

	service.Listen()
}

func NewService(service string, data settings.Webhook) (serviceInterface Service, err error) {
	switch service {
	case "slack":
		serviceInterface = &Slack{Webhook: data}
	default:
		err = errors.New(fmt.Sprintf("Unrecognized service '%v'", service))
	}

	return serviceInterface, err
}
