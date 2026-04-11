# OAuth2 Client Demo

This is a small demo OAuth 2.0 client for the companion auth server in `../oauth2`.

The project is intentionally minimal and meant for learning the authorization code flow on `localhost`.

## What It Does

- builds the `/oauth/authorize` URL
- generates and stores `state`
- handles the `/callback`
- exchanges the authorization code at `/oauth/token`

## Demo-Only Defaults

- `client_id`: `web_client`
- `client_secret`: `demo-client-secret`
- `redirect_uri`: `http://localhost:8081/callback`

These are fake local demo credentials only.
Do not copy this setup into production.

## Run

Start the auth server first from `../oauth2`, then run the client:

```bash
go run .
```

Open `http://localhost:8081`.

## Not Production Ready

This demo still simplifies or omits several things you would need in a real OAuth client:

- secure secret management
- HTTPS-only deployment
- hardened cookies
- proper token response handling
- PKCE
- refresh token handling

## Goal

Use this repo as a local playground for understanding the browser redirect and backend token exchange parts of OAuth.
