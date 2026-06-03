# temp-backend-pipeline

Spin up a real, public HTTP backend for 1–60 minutes directly from a PR comment — no local environment, no staging server needed.

Useful for testing frontend integrations, mobile clients, Postman collections, or any scenario where you need a live API endpoint on demand.

## How it works

Comment on a PR with the duration:

```
/test-deploy sprites 30
```

Within ~3-4 min the workflow replies with a live URL:

```
🚀 Backend live at https://backend-12345.sprites.app — expires in 30 min

POST   /contacts
GET    /contacts
GET    /contacts/:id
PUT    /contacts/:id
DELETE /contacts/:id
```

The backend powers off automatically after the requested time. No runner kept alive waiting.

## Provider

[sprites.dev](https://sprites.dev) — Firecracker microVM with a public HTTPS URL. Boot ~3-4 min, cost ~$0.09/hr.

## API reference

In-memory contacts CRUD. Data resets when the session ends.

```
GET    /health               → 200 OK
POST   /contacts             → 201 Created
GET    /contacts             → 200 [ ]
GET    /contacts/:id         → 200 | 404
PUT    /contacts/:id         → 200 | 404
DELETE /contacts/:id         → 204 | 404
```

**Contact shape:**
```json
{ "name": "Alice", "email": "alice@example.com", "phone": "optional" }
```

## Setup

### Secrets

| Secret | Where to get it |
|--------|-----------------|
| `REPO_TOKEN` | GitHub fine-grained token — `Contents: R&W`, `Pull requests: R&W`, `Issues: R&W` |
| `SPRITES_TOKEN` | [sprites.dev](https://sprites.dev) dashboard → API tokens — format: `org-slug/org-id/token-id/token-value` |