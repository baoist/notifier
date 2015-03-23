package settings

import (
	"fmt"
)

type Channels []string

var (
	prefixes = map[string]string{
		"#": "%23",
		"@": "%40",
	}
)

func (c Channels) PublicChannels() (public Channels) {
	for _, channel := range c {
		public = append(public, fmt.Sprint(prefixes["#"], channel))
	}

	return public
}

func (c Channels) UserChannels(user string) (assignees Channels) {
	for _, watcher := range c {
		if watcher == user {
			assignees = append(assignees, fmt.Sprint(prefixes["@"], watcher))
		}
	}

	return assignees
}
