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
	slackBaseUrl = "https://slack.com/api/chat.postMessage?token=%s&channel=%s&text=%s"
	statusColors = map[string]string{
		"error":   "#d60b09",
		"warning": "#dfa138",
		"success": "#36a755",
		"default": "#30d8e5",
	}
	author = Author{
		Name: "Software for Good",
		Link: "https://github.com/softwareforgood",
		Icon: "http://softwareforgood.com/wp-content/themes/sfg4/favicon.png?v=2",
	}
)

type Slack struct {
	settings.Webhook
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

type Attachments struct {
	Attachments []Attachment `json:"attachments"`
}

type Author struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Icon string `json:"icon"`
}

func NewAttachment(level string, title string, text string) {
	attachment := Attachment{
		Color:      statusColors[level],
		AuthorName: author.Name,
		AuthorLink: author.Link,
		AuthorIcon: author.Icon,
		Title:      title,
		Text:       text,
	}

	attachments := Attachments{[]Attachment{attachment}}

	z, err := json.Marshal(attachments)
	if err != nil {
		log.Fatal("fail marshal")
	}
	fmt.Printf("%s\n", z)

	//s.
}

func (s Slack) queryURL(message string) string {
	escaped := url.QueryEscape(message)
	return fmt.Sprintf(slackBaseUrl, s.Token, "%23hearthstone", escaped)
}

func (s Slack) Reassignment(user string) {
	if len(s.Watchers) == 0 {
		return
	}

	for _, watcher := range s.Watchers {
		fmt.Println(watcher)
	}
}

func (s Slack) Push(e *webhook.PushEvent) {
	message := fmt.Sprintf("%s pushed to <%s|%s>", e.Pusher.Email, e.Repository.URL, e.Repository.Name)

	_, err := http.Get(s.queryURL(message))
	if err != nil {
		log.Println(err)
	}
}

func (s Slack) PullRequest(e *webhook.PullRequestEvent) {
	var message string

	prefix := fmt.Sprintf("[%s]", e.PullRequest.Head.Repo.FullName)
	suffix := fmt.Sprintf("<%s|#%v %s> by <%s|%s>",
		e.PullRequest.HTMLURL,
		e.Number,
		e.PullRequest.Title,
		e.PullRequest.User.URL,
		e.PullRequest.User.Login)

	switch e.Action {
	case "opened":
		message = fmt.Sprintf("%s opened a new pull request %s", prefix, suffix)
	case "closed":
		message = fmt.Sprintf("%s deleted pull request %s", prefix, suffix)
	default:
		message = fmt.Sprintf("%s new action (%s) on pull request %s", prefix, e.Action, suffix)
	}

	_, err := http.Get(s.queryURL(message))
	if err != nil {
		log.Println(err)
	}
}

func (s *Slack) Listen() {
	NewAttachment("error", "Title", "body")
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), webhook.New(s.Secret, s))
}
