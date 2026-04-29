# Backend Deployment

This backend is ready to deploy on a free service that supports Go web apps and PostgreSQL.

## Recommended setup

- Backend host: Render web service
- Production database: Neon Postgres or Render Postgres

This codebase expects these environment variables in production:

- `DATABASE_URL`
- `JWT_SECRET`
- `PORT`
- `ALLOWED_ORIGINS`

## Required environment variables

### `DATABASE_URL`

Use the full Postgres connection string from your database provider.

Example:

```text
postgres://user:password@host:5432/dbname?sslmode=require
```

### `JWT_SECRET`

Set this to a long random secret in production.

### `PORT`

Most hosts inject this automatically. The backend now honors it.

### `ALLOWED_ORIGINS`

Comma-separated frontend origins that may call the API.

Example:

```text
https://your-frontend.vercel.app,https://courseshare.example.com
```

## Render deployment steps

1. Push the repo to GitHub.
2. In Render, create a new Web Service from the repo.
3. Set the root directory to `backend` if Render asks for it.
4. Use this build command:

```text
go build -o app .
```

5. Use this start command:

```text
./app
```

6. Add environment variables:

```text
DATABASE_URL=...
JWT_SECRET=...
ALLOWED_ORIGINS=https://your-frontend-url
```

7. Deploy the service.

## Database setup

The app runs GORM auto-migrations on startup for:

- `users`
- `courses`
- `notes`
- `helpful_votes`
- `saved_notes`

That means the production schema is created automatically once `DATABASE_URL` is valid.

## Backend checklist for demo readiness

After deploy, verify:

1. `GET /` returns `Server is running`
2. registration works
3. login returns a JWT
4. course creation and listing work
5. note create/edit/delete works
6. My Notes works
7. Saved Notes works
8. helpful votes work
9. tag filtering works

## Local production-style run

You can smoke test locally with:

```bash
DATABASE_URL="postgres://..." \
JWT_SECRET="change-me" \
ALLOWED_ORIGINS="http://localhost:4200" \
go run .
```
