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

func (s Slack) NewAttachment(level string, pretext string, title string, text string) (attachment Attachment) {
	attachment = Attachment{
		Color:      statusColors[level],
		Pretext:    pretext,
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

func (s Slack) pretext(action string, repository webhook.Repository, user webhook.User) string {
	prefix := fmt.Sprintf("[<%s|%s>]", repository.HTMLURL, repository.FullName)

	return fmt.Sprintf("%s %s by <%s|%s>",
		prefix,
		action,
		user.HTMLURL,
		user.Login)
}

func (s Slack) Issue(e *webhook.IssuesEvent) {
	title := fmt.Sprintf("<%s|%s>", e.Issue.HTMLURL, e.Issue.Title)

	switch e.Action {
	case "opened", "closed":
		attachment := s.NewAttachment("success",
			s.pretext(fmt.Sprintf("Issue %s", e.Action), e.Repository, e.Issue.User),
			title,
			e.Issue.Body)

		s.Notify(s.Public.PublicChannels(), attachment)
	case "assigned":
		attachment := s.NewAttachment("default",
			s.pretext("Issue assigned to you", e.Repository, e.Issue.User),
			title,
			e.Issue.Body)

		s.Notify(s.Watchers.UserChannels(e.Issue.Assignee.Login), attachment)
	}
}

func (s Slack) PullRequest(e *webhook.PullRequestEvent) {
	title := fmt.Sprintf("<%s|%s>", e.PullRequest.HTMLURL, e.PullRequest.Title)

	switch e.Action {
	case "opened", "closed":
		attachment := s.NewAttachment("success",
			s.pretext(fmt.Sprintf("Pull Request %s", e.Action), e.Repository, e.PullRequest.User),
			title,
			e.PullRequest.Body)

		s.Notify(s.Public.PublicChannels(), attachment)
	case "assigned":
		attachment := s.NewAttachment("success",
			s.pretext("Pull Request assigned to you", e.Repository, e.PullRequest.User),
			title,
			e.PullRequest.Body)

		s.Notify(s.Watchers.UserChannels(e.PullRequest.Assignee.Login), attachment)
	}
}

func (s *Slack) Listen() {
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), webhook.New(s.Secret, s))
}
