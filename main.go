package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/n8n-poc/client"
	"example.com/n8n-poc/config"

	"example.com/n8n-poc/errors"

	"github.com/gin-gonic/gin"
)

// Workflow represents the structure of the workflow JSON to be sent to n8n
type Workflow struct {
	Name        string                 `json:"name" binding:"required"`
	Nodes       []Node                 `json:"nodes" binding:"required"`
	Connections map[string]Connection  `json:"connections" binding:"required"`
	Settings    map[string]interface{} `json:"settings" binding:"required"`
	StaticData  map[string]interface{} `json:"staticData"`
}

// Node represents each node in the workflow
type Node struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	TypeVersion float64                `json:"typeVersion"`
	Position    []float64              `json:"position"`
	Parameters  map[string]interface{} `json:"parameters"`
	NotesInFlow bool                   `json:"notesInFlow"`
	Credentials map[string]interface{} `json:"credentials"`
	Notes       string                 `json:"notes"`
	WebhookID   string                 `json:"webhookId"`
}

// Connection represents the connection between nodes
type Connection struct {
	Main [][]ConnectionDetail `json:"main"`
}

// ConnectionDetail represents the details of a connection
type ConnectionDetail struct {
	Node  string `json:"node"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}

// Settings represents the settings for the workflow
type Settings struct {
	ExecutionOrder string `json:"executionOrder"`
}

type CreateWorkflowRequest struct {
	Workflow Workflow `json:"workflow"`
}

type CreateWorkflowResponse struct {
	WorkflowId string `json:"workflowId"`
	WebhookId  string `json:"webhookId"`
	Message    string `json:"message"`
}

type ExecuteWorkflowRequest struct {
	WorkflowId string `uri:"workflowId"`
	WebhookId  string `json:"webhookId" binding:"required"`
	Query      string `json:"query" binding:"required"`
}

func init() {
	config.Load()
}

func main() {
	conf := config.Get()
	router := gin.Default()

	router.POST("/create-workflow", createWorkflowHandler)
	router.POST("/execute-workflow/:workflowId", executeWorkflowHandler)

	router.Run(":" + conf.PORT)
}

func createWorkflowHandler(c *gin.Context) {
	var req CreateWorkflowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &errors.ErrResponse{
			Message: "Invalid request body",
		})
		return
	}

	payloadBytes, err := json.Marshal(req.Workflow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ErrResponse{
			Message: "Failed to marshal request body",
		})
		return
	}

	client := client.New()

	// create workflow
	resp, respErr := client.Do(c, http.MethodPost, "api/v1/workflows", payloadBytes)
	if respErr != nil {
		c.JSON(respErr.Status, respErr)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	var createdWorkflow struct {
		ID string `json:"id"`
	}
	// var createdWorkflow interface{}
	if err = json.Unmarshal(body, &createdWorkflow); err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	// activate workflow
	resp, respErr = client.Do(c, http.MethodPost, fmt.Sprintf("api/v1/workflows/%s/activate", createdWorkflow.ID), nil)
	if respErr != nil {
		c.JSON(respErr.Status, respErr)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	c.JSON(http.StatusCreated, &CreateWorkflowResponse{
		WorkflowId: createdWorkflow.ID,
		WebhookId:  req.Workflow.Nodes[0].WebhookID,
		Message:    "Workflow created successfully",
	})
}

func executeWorkflowHandler(c *gin.Context) {
	var req ExecuteWorkflowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &errors.ErrResponse{
			Message: "Invalid request body",
		})
		return
	}

	req.WorkflowId = c.Param("webhookId")

	payload := struct {
		Query string `json:"query"`
	}{
		Query: req.Query,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ErrResponse{
			Message: "Failed to marshal request body",
		})
		return
	}

	client := client.New()

	resp, respErr := client.Do(c, http.MethodPost, fmt.Sprintf("webhook/%s", req.WebhookId), payloadBytes)
	if respErr != nil {
		c.JSON(respErr.Status, respErr)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	var workflowResp interface{}
	if err = json.Unmarshal(body, &workflowResp); err != nil {
		c.JSON(http.StatusInternalServerError, &errors.ErrResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, workflowResp)
}
