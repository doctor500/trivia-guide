package main

import (
	"os"

	"github.com/faldyfind/trivia-numbers/api"
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
