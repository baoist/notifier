package watch

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/baoist/notifier/settings"
	"github.com/rjeczalik/gh/webhook"
)

var (
	slackBaseUrl = "https://slack.com/api/chat.postMessage?token=%s&channel=%s&text=%s"
)

type Slack struct {
	settings.Webhook
}

type Attachment struct {
	Attachments []struct {
		Fallback    string
		Color       string
		Pretext     string
		Author_name string
		Author_link string
		Author_icon string
		Title       string
		Title_link  string
		Text        string
		Image_url   string
		Fields      []struct {
			Title string
			Value string
			Short bool
		}
	}
}

func (s Slack) queryURL(message string) string {
	escaped := url.QueryEscape(message)
	return fmt.Sprintf(slackBaseUrl, s.Token, "%20hearthstone", escaped)
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

func (s *Slack) Ping() {
	fmt.Println("Pong")
}

func (s *Slack) Listen() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), webhook.New(s.Secret, s)))
}
