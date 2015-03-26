# Notifier

Notifies users when an issue or pull request is assigned to him/her. Notifies to specified public channels when a pull request or issues is opened or closed.

## Usage

Copy `[webhooks.yml.sample](webhooks.yml.example)` as `webhooks.yml`, customize it to your needs.

Note:
- You'll need set up a project/organization webhook in Github via https://github.com/<owner>/<project>/settings/hooks or https://github.com/organizations/<organization>/settings/hooks.
- A Slack API token will need to be [generated](https://api.slack.com/tokens).
- Author is optional.

## TODO:
- Map public channels to specific repositories.
- Add hipchat(?)
