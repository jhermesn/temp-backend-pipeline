# temp-backend-pipeline

Spin up a real, public HTTP backend for 1–60 minutes directly from a PR comment without need for a local environment, staging server, or paying extra costs.

Useful for testing frontend integrations, mobile clients, Postman collections, or any scenario where you need a live API endpoint on demand.

## How it works

Comment on a PR with the provider and duration:

```
/test-deploy sprites 30
/test-deploy aws 15
```

Within ~2 minutes the workflow replies with a live URL:
```
🚀 Sprites: Backend live at https://backend-12345.sprites.dev — expires in 30 min
🚀 AWS: Backend live at http://[IP_ADDRESS] — expires in 30 min

POST   /contacts
GET    /contacts
GET    /contacts/:id
PUT    /contacts/:id
DELETE /contacts/:id
```

The backend runs for the requested time, then the workflow tears it down automatically and confirms in the same thread.

## Providers

Two isolated implementations, pick the one that fits your setup:

| Provider arg | Provider | Avg boot | Cost / 60 min |
|--------------|----------|----------|---------------|
| `sprites` | sprites.dev microVM | ~3-4 min | ~$0.09 |
| `aws` | AWS EC2 Spot `t4g.nano` | ~2-3 min | ~$0.002 |

Both containerize the same Go/Gin image via Docker and expose the same API.

## API reference

The backend is a contacts CRUD with in-memory storage. Data resets when the session ends.

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
{
  "name": "Alice",
  "email": "alice@example.com",
  "phone": "optional"
}
```

## Setup

### Secrets

| Secret | Required for | Where to get it |
|--------|--------------|-----------------|
| `REPO_TOKEN` | Both providers | GitHub fine-grained token — `Contents: R&W`, `Pull requests: R&W` |
| `SPRITES_TOKEN` | sprites.dev | sprites.dev dashboard → API tokens |
| `AWS_ROLE_ARN` | AWS Spot | IAM role ARN |