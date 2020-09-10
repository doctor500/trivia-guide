# trivia-guide

## Step 1: Create the public repository

```shell
git clone https://github.com/your_name/your_repository.git
```

## Step 2: Create the api/api.go for backend

```go
package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

// NumbersHandler function
func NumbersHandler(ctx echo.Context) error {
	number := ctx.Param("number")

	url := "http://numbersapi.com/" + number

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("invalid number parameter: %s", number)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return ctx.JSON(http.StatusOK, string(body))
}

```

## Step 3: Create the unit test api/api_test.go

```go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

func TestNumbersHandler(t *testing.T) {
	type args struct {
		ctx echo.Context
	}
	type row struct {
		name    string
		args    args
		wantErr bool
	}
	tests := make([]row, 0, 2)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/numbersapi/:number")
	c.SetParamNames("number")
	c.SetParamValues("1")
	tests = append(tests, row{
		name:    "numberhandler",
		args:    args{c},
		wantErr: false,
	})

	d := e.NewContext(req, rec)
	d.SetPath("/numbersapi/:number")
	d.SetParamNames("number")
	d.SetParamValues("notanumber")
	tests = append(tests, row{
		name:    "numberhandler",
		args:    args{d},
		wantErr: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NumbersHandler(tt.args.ctx); (err != nil) != (tt.wantErr == true) {
				t.Errorf("NumbersHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

## Step 4: Commit the file

```shell
git add .
git commit -m "add api"
git push origin master
```

## Step 5: Get the api packages

```shell
go get github.com/your_name/your_repo/api
```

## Step 6: Create the numbers-be/main.go

```go
package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Set up basic auth with username and password
	e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: func(username, password string, c echo.Context) (bool, error) {
			if username == "dora" && password == "emon" {
				return true, nil
			}
			return false, nil
		},
	}))

	// Route => handler
	e.GET("/numbersapi/:number", api.NumbersHandler)

	// Start server
	//e.Logger.Fatal(e.Start(":80"))
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

```

## Step 7: Try run your app

```shell
cd numbers-be
PORT=5000 go run numbers-be/main.go
```

## Step 8: Try run the unit test

```shell
cd api
mkdir bin
go test . -coverprofile=bin/cov.out
```

## Step 9: Create .gitignore file

```shell
*.swp
*.out
**report.json**
```

## Step 10: Install and create .pre-commit-config.yml files

```shell
pip install pre-commit
or
pip3 install pre-commit
or
brew install pre-commit
```

Create .pre-commit-config.yaml

```shell
repos:
  - repo: 'https://github.com/pre-commit/pre-commit-hooks'
    rev: v3.1.0
    hooks:
      - id: trailing-whitespace
      - id: detect-private-key
      - id: detect-aws-credentials
        args: [--allow-missing-credentials]
      - id: check-merge-conflict
      - id: check-symlinks
      - id: end-of-file-fixer
```

```shell
pre-commit install
pre-commit run --all-files
```

## Step 11: Commit your changes

```shell
git add .
git commit -m "update"
git push origin master
```

## Step 12: Create new branch feature/workflow

```shell
git checkout -b feature/workflow
```

## Step 13: Create automatic unit test using workflow

```shell
mkdir -p .github/workflows/
touch .github/workflows/main.yml
```

```yml
on:
  push:
#     branches:
#       - master
#   pull_request:
#       types: [opened, synchronize, reopened]

name: Main Workflow
jobs:
  build:
    name: Compile and Test
    runs-on: ubuntu-latest
    steps:
      - name: Clone Repository
        uses: actions/checkout@v2
      - name: Setup go
        uses: actions/setup-go@v1
        with:
          go-version: '1.13'
      - name: Get dependencies
        run: go get -v -d ./...
      # Create directory bin to hold the unit test report
      - name: Create bin directory
        run: mkdir bin
      # Compile the binary
      - name: Build the binary
        run : go build -o numbers-be/bin/main numbers-be/main.go
      # 'Go test' automates testing the packages named by the import paths.
      - name: Run unit test
        run: go test ./... -coverprofile=bin/cov.out -json >> bin/report.json
      # Create artifact for the coverage results
      - name: Archive code coverage results
        uses: actions/upload-artifact@v1
        with:
          name: code-coverage-report
          path: bin
```

```shell
git add .; git commit -m "testing: workflow"; git push origin feature/workflow
```

## Step 14: Add sonarcloud workflow

- Open the sonarcloud web, access 

```shell
My Account -> Security -> Generate Tokens
```

Put the token to github secrets

```shell
Repo main page -> Settings -> Secrets
```

- And create the SONAR_TOKEN with the value of generated token.

Append to the main.yml

```yml
  sonarCloudTrigger:
    needs: build
    name: SonarCloud Trigger
    runs-on: ubuntu-latest
    steps:
      - name: Clone Repository
        uses: actions/checkout@v2
        with:
          # Disabling shallow clone is recommended for improving relevancy of reporting
          fetch-depth: 0
      - name: Download code coverage results
        uses: actions/download-artifact@v1
        with:
          name: code-coverage-report
          path: bin
      - name: Analyze with SonarCloud
        uses: sonarsource/sonarcloud-github-action@master
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

- Open your sonarcloud page, access the

```shell
My Projects -> Analyze new Project -> Create Project Manually -> Create Another organization -> Create an Organization manually -> Free Plan
```

- Also input the Project Key

- Then sonar-project.properties file

```yml
sonar.projectKey=PROJECT_KEY
sonar.organization=PROJECT_ORGANIZATION

# This is the name and version displayed in the SonarCloud UI.
sonar.projectName=trivia-numbers
sonar.projectVersion=1.0

# Path is relative to the sonar-project.properties file. Replace "\" by "/" on Windows.
sonar.sources=.
sonar.exclusions=**/*_test.go,**/vendor/**,**/testdata/*,numbers-fe/**,numbers-be/**

sonar.tests=.
sonar.test.inclusions=**/*_test.go
sonar.go.coverage.reportPaths=/github/workspace/bin/cov.out
sonar.go.tests.reportPaths=/github/workspace/bin/report.json
```

## Step 15: Test add another function to api.go

```go
// HoroscopeHandler function
func HoroscopeHandler(ctx echo.Context) error {
	star := ctx.Param("star")

	url := "http://horoscope-api.herokuapp.com/horoscope/today/" + star

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("invalid string parameter: %s", star)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return ctx.JSON(http.StatusOK, string(body))
}
```

- Edit the main.yml and set the condition to pull request only.

- Push the changes.

## Step 16: Try the PR
