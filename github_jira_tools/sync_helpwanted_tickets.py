import requests
import click
import pprint
from requests.auth import HTTPBasicAuth

from utils import create_github_issues

@click.command()
@click.option('--jira-token', '-j', prompt='Your Jira access token', help='The token used to authenticate the user against Jira.')
@click.option('--jira-username', '-u', prompt='Your Jira username', help='Username of the user to get the ticket information.')
@click.option('--github-token', '-g', prompt='Your Github access token', help='The token used to authenticate the user against Github.')
@click.option('--webhook-url', '-w', help='Webhook URL to send the list of created issues', default='')
@click.option('--dry-run/--no-dry-run', help='Skip actually creating any tickets', default=False)
@click.option('--debug/--no-debug', help='Dump debugging information.', default=False)
def cli(jira_username, jira_token, github_token, webhook_url, dry_run, debug):
    data = {
            "jql":"project = MM AND status in (Open, Reopened) AND fixversion = \"Help Wanted\" AND \"GITHUB ISSUE\" IS EMPTY AND type != EPIC",
            "maxResults": 100,
            "fields": ["summary", "description"],
    }
    resp = requests.post(
        "https://mattermost.atlassian.net/rest/api/2/search",
        json=data,
        auth=HTTPBasicAuth(jira_username, jira_token)
    )
    issues = resp.json()['issues']
    if debug:
        pprint.pprint(resp.json())

    log = ""
    if len(issues) > 0:
        log = create_github_issues(jira_username, jira_token, github_token, 'sphotonchat/server', ['Help Wanted', 'Up For Grabs'], issues, dry_run)

    if log:
        if webhook_url:
            requests.post(webhook_url, json={"text": log})
        else:
            print(log)

if __name__ == "__main__":
    cli()
