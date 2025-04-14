package testdata

import (
	"net/http"
	"testing"
)

func _(tb testing.TB, client *http.Client, req *http.Request, statusCode int) {
	Errorf(tb, false, "expected %v status code, got %v %v", statusCode, req.URL, client.Transport)
}

func Errorf(tb testing.TB, cond bool, format string, args ...any) {}
