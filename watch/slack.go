package watch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/baoist/notifier/settings"
	"github.com/rjeczalik/gh/webhook"
)

var (
	slackBaseUrl = "https://slack.com/api/chat.postMessage?token=%s&channel=%s&attachments=[%s]"
	statusColors = map[string]string{
		"error":   "#d60b09",
		"warning": "#dfa138",
		"success": "#36a755",
		"default": "#30d8e5",
	}
)

type Slack struct {
	settings.Webhook
	settings.Channels
}

type Attachment struct {
	Fallback   string `json:"fallback"`
	Color      string `json:"color"`
	Pretext    string `json:"pretext"`
	AuthorName string `json:"author_name"`
	AuthorLink string `json:"author_link"`
	AuthorIcon string `json:"author_icon"`
	Title      string `json:"title"`
	TitleLink  string `json:"title_link"`
	Text       string `json:"text"`
	ImageURL   string `json:"image_url"`
	Fields     []struct {
		Title string `json:"title"`
		Value string `json:"value"`
		Short bool   `json:"short"`
	} `json:"fields"`
}

type Author struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Icon string `json:"icon"`
}

func (s Slack) NewAttachment(level string, title string, text string) (attachment Attachment) {
	attachment = Attachment{
		Color:      statusColors[level],
		AuthorName: s.Author.Name,
		AuthorLink: s.Author.Link,
		AuthorIcon: s.Author.Icon,
		Title:      title,
		Text:       text,
		Fallback:   text,
	}

	return attachment
}

func (s Slack) queryURL(channel string, message string) string {
	encoded := struct {
		channel string
		message string
	}{
		url.QueryEscape(channel),
		url.QueryEscape(message),
	}
	return fmt.Sprintf(slackBaseUrl, s.Token, encoded.channel, encoded.message)
}

func (s *Slack) Notify(channels settings.Channels, attachment Attachment) {
	for _, channel := range channels {
		attachmentJSON, err := json.Marshal(attachment)
		if err != nil {
			log.Fatal("Unable to parse JSON.")
		}

		_, errz := http.Get(s.queryURL(channel, string(attachmentJSON)))
		if errz != nil {
			log.Println(errz)
		}
	}
}

func (s Slack) Push(e *webhook.PushEvent) {
	message := fmt.Sprintf("%s pushed to <%s|%s>", e.Pusher.Email, e.Repository.URL, e.Repository.Name)
	attachment := s.NewAttachment("success", "Pushed", message)

	s.Notify(s.Public.PublicChannels(), attachment)
}

func (s Slack) PullRequest(e *webhook.PullRequestEvent) {
	var message string

	prefix := fmt.Sprintf("[<%s|%s>]", e.PullRequest.Head.Repo.HTMLURL, e.PullRequest.Head.Repo.FullName)
	suffix := fmt.Sprintf("<%s|#%v %s> by <%s|%s>.\n%s",
		e.PullRequest.HTMLURL,
		e.Number,
		e.PullRequest.Title,
		e.PullRequest.User.HTMLURL,
		e.PullRequest.User.Login,
		e.PullRequest.Body)

	switch e.Action {
	case "opened":
		message = fmt.Sprintf("%s opened a new pull request %s", prefix, suffix)
		attachment := s.NewAttachment("success", "Opened Pull Request", message)

		s.Notify(s.Public.PublicChannels(), attachment)
	case "closed":
		message = fmt.Sprintf("%s deleted pull request %s", prefix, suffix)
		attachment := s.NewAttachment("default", "Closed Pull Request", message)

		s.Notify(s.Public.PublicChannels(), attachment)
	case "assigned":
		message = fmt.Sprintf("%s assigned pull request <%s|%s> to you.\n%s",
			prefix,
			e.PullRequest.HTMLURL,
			e.PullRequest.Title,
			e.PullRequest.Body)
		attachment := s.NewAttachment("default", "Assigned Pull Request To You", message)

		s.Notify(s.Watchers.UserChannels(e.PullRequest.Assignee.Login), attachment)
	}
}

func (s *Slack) Listen() {
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), webhook.New(s.Secret, s))
}
