# ginProject

Go API with JWT auth and Google OAuth.

## Google OAuth setup

1. Create OAuth credentials in [Google Cloud Console](https://console.cloud.google.com/apis/credentials).
2. Application type: **Web application**.
3. Authorized redirect URI:
   ```
   http://localhost:8080/api/v1/oauth/google/callback
   ```
4. Copy Client ID and Client Secret into `.env`:

```env
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
OAUTH_REDIRECT_URL=http://localhost:8080/api/v1/oauth/google/callback
FRONTEND_URL=http://localhost:3000
```

5. Restart the Go server. On the React app, use **Continue with Google** on login or signup.

## OAuth flow

- `GET /api/v1/oauth/google` — redirects to Google
- `GET /api/v1/oauth/google/callback` — creates or links the user, then redirects to `{FRONTEND_URL}/oauth/callback` with `token`, `userId`, `email`, and `role`
