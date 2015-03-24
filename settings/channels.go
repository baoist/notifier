package settings

import (
	"fmt"
)

type Channels []string

func (c Channels) PublicChannels() (public Channels) {
	for _, channel := range c {
		public = append(public, fmt.Sprint("#", channel))
	}

	return public
}

func (c Channels) UserChannels(user string) (assignees Channels) {
	for _, watcher := range c {
		if watcher == user {
			assignees = append(assignees, fmt.Sprint("@", watcher))
		}
	}

	return assignees
}
