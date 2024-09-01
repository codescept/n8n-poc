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

If you want to create your own workflow. Make sure it must have a webhook node so the workflow can be executed.

### Sample Workflow

Use the above request to create sample workflow. Remember to add credentials manually in the n8n dashbaord
if you are using the above workflow.

```bash
curl --location 'http://localhost:6000/create-workflow' \
--header 'Content-Type: application/json' \
--data '{
    "workflow": {
        "name": "RAG Automated Workflow",
        "nodes": [
            {
                "parameters": {
                    "httpMethod": "POST",
                    "path": "1",
                    "responseMode": "lastNode",
                    "options": {}
                },
                "id": "0b4e243c-1fee-48b3-9c22-5ccc29a072cc",
                "name": "Webhook",
                "type": "n8n-nodes-base.webhook",
                "typeVersion": 2,
                "position": [
                    820,
                    620
                ],
                "webhookId": "1",
                "credentials": {}
            },
            {
                "parameters": {
                    "method": "POST",
                    "url": "https://gateway-dev.on-demand.io/chat/v1/sessions",
                    "authentication": "genericCredentialType",
                    "genericAuthType": "httpHeaderAuth",
                    "sendBody": true,
                    "specifyBody": "json",
                    "jsonBody": "={\n  \"pluginIds\": [],\n  \"externalUserId\": \"1\"\n}",
                    "options": {}
                },
                "id": "fd5d691a-f2b3-4ed7-97ea-fada80d4d214",
                "name": "Create Chat Session",
                "type": "n8n-nodes-base.httpRequest",
                "typeVersion": 4.2,
                "position": [
                    1060,
                    620
                ],
                "notesInFlow": false,
                "credentials": {
                    "httpHeaderAuth": {
                        "id": "mcZvcRyVJI2aZ7Bl",
                        "name": "Header Auth account"
                    }
                }
            },
            {
                "parameters": {
                    "method": "POST",
                    "url": "=https://gateway-dev.on-demand.io/chat/v1/sessions/{{ $json.data.id }}/query",
                    "authentication": "genericCredentialType",
                    "genericAuthType": "httpHeaderAuth",
                    "sendBody": true,
                    "specifyBody": "json",
                    "jsonBody": "={\n    \"endpointId\": \"predefined-openai-gpt4o-mini\",\n    \"query\": \"{{ $('\''Webhook'\'').item.json.body.query.replace(/\"/g, '\''\\\\\"'\'') }}\",\n    \"pluginIds\": [\n        \"plugin-1714419354\",\n        \"plugin-1713924030\"\n    ],\n    \"responseMode\": \"sync\"\n}",
                    "options": {}
                },
                "id": "43ebed18-0667-4277-bdc4-fe86ebb7b184",
                "name": "Submit Query",
                "type": "n8n-nodes-base.httpRequest",
                "typeVersion": 4.2,
                "position": [
                    1280,
                    620
                ],
                "notesInFlow": true,
                "credentials": {
                    "httpHeaderAuth": {
                        "id": "mcZvcRyVJI2aZ7Bl",
                        "name": "Header Auth account"
                    }
                },
                "notes": "Submit Query"
            }
        ],
        "connections": {
            "Webhook": {
                "main": [
                    [
                        {
                            "node": "Create Chat Session",
                            "type": "main",
                            "index": 0
                        }
                    ]
                ]
            },
            "Create Chat Session": {
                "main": [
                    [
                        {
                            "node": "Submit Query",
                            "type": "main",
                            "index": 0
                        }
                    ]
                ]
            }
        },
        "settings": {
            "executionOrder": "v1"
        }
    }
}'
```

## Execute Workflow

Execute workflow by calling `/execute-workflow` endpoint
