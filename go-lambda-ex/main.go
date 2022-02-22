package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"go_src/go-lambda-ex/database"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

// global state로 선언해 lambda의 메모리 낭비를 줄인다.
var logger *zap.Logger
var db *sql.DB

func init() {
	l, _ := zap.NewProduction()
	logger = l
	defer logger.Sync()

	dbConnection, err := database.GetConnection()
	if err != nil {
		logger.Error("error connection to database", zap.Error(err))
		panic(err)
	}

	// Ping 메소드는 DB가 제대로 연결되어있는지를 확인하고 필요한 경우 연결을 설정한다.
	dbConnection.Ping()
	if err != nil {
		logger.Error("error pinging database", zap.Error(err))
		panic(err)
	}
	db = dbConnection
}

type DefaultResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GETEmployeesResponse struct {
	Employees []*database.Employee `json:"employees"`
}

func MyHandler(ctx context.Context, event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse
	// lambda에 대해 여러개의 endpoint를 가지고 있다면 event.Path를 사용한다.
	// event.Path를 사용해 endpoint에 접근할 수 있다.

	// 어떤 request를 받던간에 request에 Body가 있다면 log를 띄울 것임
	logger.Info("received event", zap.Any("method", event.HTTPMethod), zap.Any("path", event.Path), zap.Any("body", event.Body))

	if event.Path == "/migrate" {
		err := database.CreateEmployeesTable(ctx, db)

		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusInternalServerError),
				Message: "could not create employees table",
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&DefaultResponse{
			Status:  string(http.StatusOK),
			Message: "success",
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil

	} else if event.Path == "/employees" && event.HTTPMethod == http.MethodPost {
		// create a new employee
		employee := &database.Employee{}
		err := json.Unmarshal([]byte(event.Body), employee)
		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusBadRequest),
				Message: err.Error(),
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       string(body),
			}, nil
		}

		err = database.CreateEmployee(ctx, db, employee.Email, employee.FirstName, employee.LastName)
		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusInternalServerError),
				Message: err.Error(),
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       string(body),
			}, nil
		}
		body, _ := json.Marshal(&DefaultResponse{
			Status:  string(http.StatusOK),
			Message: "success",
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil
	} else if event.Path == "/employees" && event.HTTPMethod == http.MethodGet {
		// get all employees
		employees, err := database.GetEmployees(ctx, db)
		if err != nil {
			body, _ := json.Marshal(&DefaultResponse{
				Status:  string(http.StatusInternalServerError),
				Message: "error fetching employees",
			})

			return &events.APIGatewayProxyResponse{
				StatusCode: http.StatusOK,
				Body:       string(body),
			}, nil
		}

		body, _ := json.Marshal(&GETEmployeesResponse{
			Employees: employees,
		})

		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}, nil
	} else {
		body, _ := json.Marshal(&DefaultResponse{
			Status:  string(http.StatusOK),
			Message: "default path",
		})

		res = &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
	}

	return res, nil
}

func main() {
	lambda.Start(MyHandler)
}
