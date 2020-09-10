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
