package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var (
	env                = os.Getenv("env")
	slogger  *zap.SugaredLogger
)

//Response ...
type Response struct {
	ValidationErrors []string      `json:"validationErrors"`
	Error            string        `json:"error"`
	Data             *Default `json:"data"`
}

//Default ...
type Default struct {
	Text string `json:"text"`
}

func handleRequest(ctx context.Context, defaultRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var defaultBody Default
	err := json.Unmarshal([]byte(defaultRequest.Body), &defaultBody)
	if err != nil {
		slogger.Error(fmt.Sprintf("%s %s", "an error occurred trying to unmarshal request", err.Error()))
		return sendResponse(&Response{ValidationErrors: []string{"Invalid Request body, expected fields username and password"}}, http.StatusBadRequest), nil
	}

	slogger.Infof("default text from body bytes %s on env %s", defaultBody.Text, env)

	responseBodyBytes, err := json.Marshal(Response{Data: &defaultBody})
	if err != nil {
		slogger.Error(fmt.Sprintf("%s %s", "an error occurred trying to marshal response", err.Error()))
		return sendResponse(&Response{Error: "An unexpected error occurred"}, http.StatusInternalServerError), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBodyBytes),
	}, nil

}


func sendResponse(response *Response, statusCode int) events.APIGatewayProxyResponse {
	responseBody, err := json.Marshal(response)
	if err != nil {
		slogger.Error(fmt.Sprintf("%s %s", "an error occurred while trying to marshal response: ", err.Error()))
		return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       http.StatusText(http.StatusInternalServerError),
		}
	}
	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: statusCode,
	}
}

func init() {
    // Init the logger outside of the handler
	config := zap.NewProductionConfig()
	config.InitialFields = map[string]interface{}{
		"function" : "defaultFunction",
	}
    logger, _ := config.Build()
	slogger = logger.Sugar()
}

func main() {
	lambda.Start(handleRequest)
}
