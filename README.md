# temp-backend-pipeline

PoC that spins up a temporary Go/Gin contacts CRUD backend for 1–60 minutes, triggered by a slash command on a PR comment. Two providers are available — sprites.dev and AWS EC2 Spot — each in its own workflow file.

## Usage

Comment on any **PR** (not a plain issue):

```
/test-deploy 30          # sprites.dev, 30 minutes
/test-deploy-spot 15     # AWS EC2 Spot (t4g.nano), 15 minutes
```

The workflow replies with the live URL, waits the requested duration, tears everything down, and posts a confirmation comment.

## API

| Method | Path | Body |
|--------|------|------|
| `GET` | `/health` | — |
| `POST` | `/contacts` | `{"name":"Alice","email":"alice@example.com","phone":"optional"}` |
| `GET` | `/contacts` | — |
| `GET` | `/contacts/:id` | — |
| `PUT` | `/contacts/:id` | same as POST |
| `DELETE` | `/contacts/:id` | — |

Storage is in-memory. Data is lost when the session ends.

## Cost

| Provider | Instance | Cost / 60 min | Notes |
|----------|----------|--------------|-------|
| sprites.dev | 1 vCPU / 512 MB | ~$0.09 | No AWS account needed |
| AWS Spot (`t4g.nano`) | 0.5 vCPU / 512 MB | ~$0.0008 | ~100× cheaper |

## Required secrets

| Secret | Used by | How to obtain |
|--------|---------|---------------|
| `REPO_TOKEN` | `dispatch.yml` | GitHub PAT (classic) with `repo` scope |
| `SPRITES_TOKEN` | `deploy-sprites.yml` | sprites.dev dashboard → API tokens |
| `AWS_ACCESS_KEY_ID` | `deploy-spot.yml` | IAM user — see permissions below |
| `AWS_SECRET_ACCESS_KEY` | `deploy-spot.yml` | Same IAM user |

### Minimum IAM permissions for spot deployment

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:RunInstances",
        "ec2:DescribeInstances",
        "ec2:TerminateInstances",
        "ec2:CreateSecurityGroup",
        "ec2:DescribeSecurityGroups",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:DeleteSecurityGroup",
        "ec2:DescribeImages"
      ],
      "Resource": "*"
    }
  ]
}
```

## Local development

```bash
cd backend
go test ./... -v          # run all tests
go run .                  # start server on :8080

# or with Docker
docker build -t temp-backend .
docker run -p 8080:8080 temp-backend
curl localhost:8080/health
```

## Repository structure

```
.github/workflows/
  dispatch.yml          # PR slash-command listener
  deploy-sprites.yml    # sprites.dev provider
  deploy-spot.yml       # AWS EC2 Spot provider
backend/
  contacts_test.go      # tests (written first — TDD)
  main.go               # Go + Gin in-memory contacts CRUD
  go.mod / go.sum
  Dockerfile
deploy/
  sprites/entrypoint.sh # runs inside the sprite
  spot/user-data.sh     # EC2 bootstrap script
```
