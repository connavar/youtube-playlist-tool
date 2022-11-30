# youtube-playlist-tool
[![Go Report Card](https://goreportcard.com/badge/github.com/connavar/youtube-playlist-tool)](https://goreportcard.com/report/github.com/connavar/youtube-playlist-tool)
[![codebeat badge](https://codebeat.co/badges/2fe38773-f5dd-4f78-ad97-2a3fee550975)](https://codebeat.co/projects/github-com-connavar-youtube-playlist-tool-main)
[![Maintainability](https://api.codeclimate.com/v1/badges/3cb5f260cb902b420fd3/maintainability)](https://codeclimate.com/github/connavar/youtube-playlist-tool/maintainability)

## Build

1. Register for a GCP YouTube API key. For instructions, see: https://developers.google.com/youtube/v3/getting-started

2. Create youtube-secret.json in the secrets directory, using youtube-secret-template.json.

The file should be of the format:
```json
{
  "installed": {
    "client_id": "<client_id>",
    "project_id": "<project_id>",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_secret": "<client_secret>",
    "redirect_uris": [
      "urn:ietf:wg:oauth:2.0:oob",
      "http://localhost"
    ]
  }
}
```
**DO NOT** commit your credentials or the youtube-secret.json file.
3. Build or run the project

```bash
go run .

# OR

go install .
```   