package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/WeTrustPlatform/charity-management-serv/seed"
)

func TestMain(m *testing.M) {
	dbInstance := DB(false)
	defer dbInstance.Close()

	seed.Populate(dbInstance, "data_test.txt", false)
	code := m.Run()
	os.Exit(code)
}

func execute(req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	Router().ServeHTTP(res, req)
	return res
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetCharities(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/charities", nil)
	res := execute(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestSearchFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/charities?search=Ig", nil)
	res := execute(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestSearchNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/charities?search=000000", nil)
	res := execute(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestEINFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/charities?ein=000587764", nil)
	res := execute(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestEINNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/charities?ein=000000", nil)
	res := execute(req)
	checkResponseCode(t, http.StatusNotFound, res.Code)
}

func TestGetCharity(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/charities/1", nil)
	res := execute(req)
	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestGetCharityNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/charities/4", nil)
	res := execute(req)
	checkResponseCode(t, http.StatusNotFound, res.Code)
}
