package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	db, _ = sql.Open("mysql", "orderserver:abc123456@tcp(localhost:3306)/orderserver?charset=utf8")
}

func TestListOrderHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ListOrderHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body, _ := ioutil.ReadAll(rr.Body)
	_, err = simplejson.NewJson(body)
	if err != nil {
		t.Errorf("response data should be json format, but got %v", rr.Body.String())
	}
}

func TestListOrderHandler1(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders?page=2&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ListOrderHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, _ := ioutil.ReadAll(rr.Body)
	_, err = simplejson.NewJson(body)
	if err != nil {
		t.Errorf("response data should be json format, but got %v", rr.Body.String())
	}
}

func TestListOrderHandler2(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders?page=2000&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ListOrderHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	body, _ := ioutil.ReadAll(rr.Body)
	_, err = simplejson.NewJson(body)
	if err != nil {
		t.Errorf("response data should be json format, but got %v", rr.Body.String())
	}
}

func Benchmark_ListOrderHandler(b *testing.B) {
	req, err := http.NewRequest("GET", "/orders?page=2&limit=10", nil)
	if err != nil {
		b.Fatal(err)
	}

	rr := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		ListOrderHandler(rr, req)
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		body, _ := ioutil.ReadAll(rr.Body)
		_, err = simplejson.NewJson(body)
		if err != nil {
			b.Errorf("response data should be json format, but got %v", rr.Body.String())
		}
	}
}

func Test_NewOrderHandler(t *testing.T) {

	origin := []string{"39.983171", "116.308479"}
	destination := []string{"39.99606", "116.353455"}
	reqData := Order{Origin: origin, Destination: destination}

	reqBody, _ := json.Marshal(reqData)
	fmt.Println("input:", string(reqBody))
	req := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewReader(reqBody),
	)

	rr := httptest.NewRecorder()
	NewOrderHandler(rr, req)

	ret := rr.Result()
	if ret.StatusCode != http.StatusOK {
		t.Errorf("expected status 200,", ret.StatusCode)
	}

	body, _ := ioutil.ReadAll(ret.Body)
	orderInfo := OrderInfo{}
	err := json.Unmarshal(body, &orderInfo)
	if err != nil {
		t.Errorf("response data should be json format, but got %v", rr.Body.String())
	}
}

func Benchmark_NewOrderHandler(b *testing.B) {

	origin := []string{"39.983171", "116.308479"}
	destination := []string{"39.99606", "116.353455"}
	reqData := Order{Origin: origin, Destination: destination}

	reqBody, _ := json.Marshal(reqData)
	req := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewReader(reqBody),
	)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		NewOrderHandler(rr, req)

		ret := rr.Result()
		if ret.StatusCode != http.StatusOK {
			b.Errorf("expected status 200, but got %v", ret.StatusCode)
		}

		body, _ := ioutil.ReadAll(ret.Body)
		orderInfo := OrderInfo{}
		err := json.Unmarshal(body, &orderInfo)
		if err != nil {
			b.Errorf("response data should be json format, but got %v", rr.Body.String())
		}
	}
}
