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
🚀 Backend live at https://backend-12345.sprites.dev — expires in 30 min

POST   /contacts
GET    /contacts
GET    /contacts/:id
PUT    /contacts/:id
DELETE /contacts/:id
```

The backend runs for the requested time, then the workflow tears it down automatically and confirms in the same thread.

## Providers

Two isolated implementations — pick the one that fits your setup:

| Provider arg | Provider | Avg boot | Cost / 60 min |
|--------------|----------|----------|--------------|
| `sprites` | [sprites.dev](https://sprites.dev) microVM | ~90 s | ~$0.09 |
| `aws` | AWS EC2 Spot `t4g.nano` | ~2 min | ~$0.001 |

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
{ "name": "Alice", "email": "alice@example.com", "phone": "optional" }
```

## Setup

### 1. Secrets

| Secret | Required for | Where to get it |
|--------|-------------|-----------------|
| `REPO_TOKEN` | Both providers | GitHub → Settings → Developer settings → PAT (classic), `repo` scope |
| `SPRITES_TOKEN` | sprites.dev | [sprites.dev](https://sprites.dev) dashboard → API tokens |
| `AWS_ROLE_ARN` | AWS Spot | IAM role ARN — see below |

### 2. AWS OIDC (spot provider only)

The spot workflow uses GitHub OIDC — no long-term keys.

**Add GitHub as an identity provider in IAM** (once per AWS account):
- Provider URL: `https://token.actions.githubusercontent.com`
- Audience: `sts.amazonaws.com`

**Create an IAM role** with this trust policy:
```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Federated": "arn:aws:iam::<account-id>:oidc-provider/token.actions.githubusercontent.com"
    },
    "Action": "sts:AssumeRoleWithWebIdentity",
    "Condition": {
      "StringEquals": { "token.actions.githubusercontent.com:aud": "sts.amazonaws.com" },
      "StringLike":  { "token.actions.githubusercontent.com:sub": "repo:jhermesn/temp-backend-pipeline:*" }
    }
  }]
}
```

**Attach this permission policy** to the role:
```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:RunInstances", "ec2:DescribeInstances", "ec2:TerminateInstances",
      "ec2:CreateSecurityGroup", "ec2:DescribeSecurityGroups",
      "ec2:AuthorizeSecurityGroupIngress", "ec2:DeleteSecurityGroup",
      "ec2:DescribeImages"
    ],
    "Resource": "*"
  }]
}
```

Set the role ARN as the `AWS_ROLE_ARN` repo secret.

## Local development

```bash
cd backend
go test ./... -v        # run tests
go run .                # server on :8080

docker build -t temp-backend .
docker run -p 8080:8080 temp-backend
curl localhost:8080/health
```
