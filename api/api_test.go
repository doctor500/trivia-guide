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
