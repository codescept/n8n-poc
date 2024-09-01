# n8n POC

## Setup

```bash

docker volume create n8n_data

docker run -d \
 --name n8n \
 -p 5678:5678 \
 -e GENERIC_TIMEZONE="Asia/Karachi" \
 -e TZ="Asia/Karachi" \
 -v n8n_data:/home/node/.n8n \
 docker.n8n.io/n8nio/n8n

```

## API Key

Create an API Key by going to - Settings > Create API Key

## Create Workflow

Create workflow by calling `/create-workflow` endpoint

Make sure to create webhook node so the workflow can be executed

## Execute Workflow

Execute workflow by calling `/execute-workflow` endpoint
