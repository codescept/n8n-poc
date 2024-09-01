# n8n Proof of Concept (POC)

## Setup Instructions

To set up the n8n instance using Docker, follow the steps below:

1. **Create a Docker Volume for Persistent Data Storage:**

    ```bash
    docker volume create n8n_data
    ```

2. **Run the n8n Container:**

    Execute the following command to run the n8n container:

    ```bash
    docker run -d \
     --name n8n \
     -p 5678:5678 \
     -e GENERIC_TIMEZONE="Asia/Karachi" \
     -e TZ="Asia/Karachi" \
     -v n8n_data:/home/node/.n8n \
     docker.n8n.io/n8nio/n8n
    ```

    This will start n8n, making it accessible on port `5678` with the timezone set to `Asia/Karachi`.

## Generating an API Key

To interact with the n8n API, you will need an API Key. Follow these steps to generate one:

1. Navigate to `Settings` within the n8n dashboard.
2. Select `Create API Key`.
3. Save the generated API Key for use in API requests.

## Creating a Workflow

Workflows can be created programmatically by calling the `/create-workflow` endpoint.

### Requirements for Workflows

- **Webhook Node:** Ensure that your workflow includes at least one webhook node. This node is essential for triggering and executing the workflow.

**Important:** If your workflow includes external services or APIs, you must manually add the corresponding credentials through the n8n dashboard.

### Sample Workflow

To create a sample workflow, use the provided API request template to call the `/create-workflow` endpoint. You also need to create credentials to run this workflow.  

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
