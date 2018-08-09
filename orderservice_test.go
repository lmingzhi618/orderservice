package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	db, _ = sql.Open("mysql", "root:abc123456@tcp(127.0.0.1:3306)/orderserver?charset=utf8")
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
	orderList := []OrderInfo{}
	err = json.Unmarshal(body, &orderList)
	if err != nil {
		t.Errorf("response data should be %v, but got %v",
			`[
				{
					"id": <order_id>,
					"distance": <total_distance>,
					"status": <ORDER_STATUS>
				},
				...
			]`, string(body))
	}
}

func TestListOrderHandler1(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders?page=2&limit=5", nil)
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
	orderList := []OrderInfo{}
	err = json.Unmarshal(body, &orderList)
	if err != nil {
		t.Errorf("response data should be %v, but got %v",
			`[
				{
					"id": <order_id>,
					"distance": <total_distance>,
					"status": <ORDER_STATUS>
				},
				...
			]`, string(body))
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
	orderList := []OrderInfo{}
	err = json.Unmarshal(body, &orderList)
	if err != nil {
		t.Errorf("response data should be %v, but got %v",
			`[
				{
					"id": <order_id>,
					"distance": <total_distance>,
					"status": <ORDER_STATUS>
				},
				...
			]`, string(body))
	}
}

func Benchmark_ListOrderHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("GET", "/orders?page=2&limit=10", nil)
		if err != nil {
			b.Fatal(err)
		}

		rr := httptest.NewRecorder()

		ListOrderHandler(rr, req)
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		body, _ := ioutil.ReadAll(rr.Body)
		orderList := []OrderInfo{}
		err = json.Unmarshal(body, &orderList)
		if err != nil {
			b.Errorf("response data should be %v, but got %v",
				`[
				{
					"id": <order_id>,
					"distance": <total_distance>,
					"status": <ORDER_STATUS>
				},
				...
			]`, string(body))
		}
	}
}

func Test_NewOrderHandler(t *testing.T) {

	origin := []string{"39.983171", "116.308479"}
	destination := []string{"39.99606", "116.353455"}
	reqData := Order{Origin: origin, Destination: destination}

	reqBody, _ := json.Marshal(reqData)
	req := httptest.NewRequest(
		http.MethodPost,
		"/order",
		bytes.NewReader(reqBody),
	)

	rr := httptest.NewRecorder()
	NewOrderHandler(rr, req)

	ret := rr.Result()
	body, _ := ioutil.ReadAll(ret.Body)
	if ret.StatusCode == http.StatusOK {
		orderInfo := OrderInfo{}
		err := json.Unmarshal(body, &orderInfo)
		if err != nil {
			t.Errorf("response data invalid, got %v", string(body))
		}
	} else if ret.StatusCode == http.StatusInternalServerError {
		m := map[string]string{}
		err := json.Unmarshal(body, &m)
		if err != nil || m["error"] != "ERROR_DESCRIPTION" {
			t.Errorf("response data invalid, got %v", string(body))
		}
	} else {
		t.Errorf("expected status 200 or 500, but get %v", ret.StatusCode)
	}
}

func Benchmark_NewOrderHandler(b *testing.B) {

	origin := []string{"39.983171", "116.308479"}
	destination := []string{"39.99606", "116.353455"}
	for i := 0; i < b.N; i++ {
		reqData := Order{Origin: origin, Destination: destination}

		reqBody, _ := json.Marshal(reqData)
		req := httptest.NewRequest(
			http.MethodPost,
			"/order",
			bytes.NewReader(reqBody),
		)

		rr := httptest.NewRecorder()
		NewOrderHandler(rr, req)

		ret := rr.Result()
		body, _ := ioutil.ReadAll(ret.Body)
		if ret.StatusCode == http.StatusOK {
			orderInfo := OrderInfo{}
			err := json.Unmarshal(body, &orderInfo)
			if err != nil {
				b.Errorf("response data style should be %v, but got %v",
					`{
						"id": <order_id>,
						"distance": <total_distance>,
						"status": "UNASSIGN"
					}`,
					string(body))
			}
		} else if ret.StatusCode == http.StatusInternalServerError {
			m := map[string]string{}
			err := json.Unmarshal(body, &m)
			if err != nil || m["error"] != "ERROR_DESCRIPTION" {
				b.Errorf("response data should be %v, but got %v",
					`{"error":"ERROR_DESCRIPTION"}`, string(body))
			}
		} else {
			b.Errorf("expected status 200 or 500, but get %v", ret.StatusCode)
		}
	}
}

func Test_TakeOrderHandler(t *testing.T) {
	reqBody := []byte(`{"status":"taken"}`)
	req := httptest.NewRequest(
		http.MethodPut,
		"/order/1",
		bytes.NewReader(reqBody),
	)

	rr := httptest.NewRecorder()
	TakeOrderHandler(rr, req)

	ret := rr.Result()
	body, _ := ioutil.ReadAll(ret.Body)
	if ret.StatusCode == http.StatusOK {
		m := map[string]string{}
		err := json.Unmarshal(body, &m)
		if err != nil || m["status"] != "SUCCESS" {
			t.Errorf("response data should be %v, but got %v",
				`"{"status":"SUCCESS"}"`, string(body))
		}
	} else if ret.StatusCode == http.StatusConflict {
		m := map[string]string{}
		err := json.Unmarshal(body, &m)
		if err != nil || m["error"] != "ORDER_ALREADY_BEEN_TAKEN" {
			t.Errorf("response data should be %v, but got %v",
				`{"error":"ORDER_ALREADY_BEEN_TAKEN"}`, string(body))
		}
	} else {
		t.Errorf("expected status 200 or 409, but get %v", ret.StatusCode)
	}
}

func Test_TakeOrderHandler1(t *testing.T) {
	reqBody := []byte(`{"status":"taken"}`)
	req := httptest.NewRequest(
		http.MethodPut,
		"/order/99999999",
		bytes.NewReader(reqBody),
	)

	rr := httptest.NewRecorder()
	TakeOrderHandler(rr, req)

	ret := rr.Result()
	body, _ := ioutil.ReadAll(ret.Body)
	if ret.StatusCode == http.StatusOK {
		m := map[string]string{}
		err := json.Unmarshal(body, &m)
		if err != nil || m["status"] != "SUCCESS" {
			t.Errorf("response data should be %v, but got %v",
				`"{"status":"SUCCESS"}"`, string(body))
		}
	} else if ret.StatusCode == http.StatusConflict {
		m := map[string]string{}
		err := json.Unmarshal(body, &m)
		if err != nil || m["error"] != "ORDER_ALREADY_BEEN_TAKEN" {
			t.Errorf("response data should be %v, but got %v",
				`{"error":"ORDER_ALREADY_BEEN_TAKEN"}`, string(body))
		}
	} else {
		t.Errorf("expected status 200 or 409, but get %v", ret.StatusCode)
	}
}

func Benchmark_TakeOrderHandler(b *testing.B) {
	reqBody := []byte(`{"status":"taken"}`)
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(
			http.MethodPut,
			"/order/1000",
			bytes.NewReader(reqBody),
		)

		rr := httptest.NewRecorder()
		TakeOrderHandler(rr, req)

		ret := rr.Result()
		body, _ := ioutil.ReadAll(ret.Body)
		if ret.StatusCode == http.StatusOK {
			m := map[string]string{}
			err := json.Unmarshal(body, &m)
			if err != nil || m["status"] != "SUCCESS" {
				b.Errorf("response data should be %v, but got %v",
					`"{"status":"SUCCESS"}"`, string(body))
			}
		} else if ret.StatusCode == http.StatusConflict {
			m := map[string]string{}
			err := json.Unmarshal(body, &m)
			if err != nil || m["error"] != "ORDER_ALREADY_BEEN_TAKEN" {
				b.Errorf("response data should be %v, but got %v",
					`{"error":"ORDER_ALREADY_BEEN_TAKEN"}`, string(body))
			}
		} else {
			b.Errorf("expected status 200 or 409, but get %v", ret.StatusCode)
		}
	}
}
