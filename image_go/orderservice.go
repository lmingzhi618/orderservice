package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/6xiao/go/Common"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	orderUnAssign = 0
	orderTaken    = 1
	pageLen       = 5
	port          = ":8080"
	logDir        = "Log"
)

var (
	logErr  = Common.ErrorLog
	logInfo = Common.InfoLog
)
var db *sql.DB

func init() {
	Common.Init(os.Stdout)
	Common.SetLogDir(logDir)

}

type Order struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

type OrderInfo struct {
	Id       int64  `json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}

func GetDistance(order Order) (int, error) {
	if 2 != len(order.Origin) || 2 != len(order.Destination) {
		logErr("Request data invalid: ", order)
		return -1, fmt.Errorf("Request data invalid: %s", order)
	}

	url := "https://apis.map.qq.com/ws/distance/v1/?mode=driving&key=2LRBZ-QHNR3-GDN3N-YDF5B-DDLYJ-ZZBWW"
	url += "&from=" + strings.Join(order.Origin, ",") + "&to=" + strings.Join(order.Destination, ",")

	logInfo("url: ", url)
	ret, err := http.Get(url)
	if err != nil {
		logErr("http Get failed:", err)
		return -1, fmt.Errorf("http Get failed: %s", err)
	}
	defer ret.Body.Close()

	cnt, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		logErr("Read http response failed:", err)
		return -1, fmt.Errorf("Read http response failed, err: %v", err)
	}

	js, err := simplejson.NewJson(cnt)
	if err != nil {
		logErr("json unmarshal failed:", err)
		return -1, fmt.Errorf("json unmarshal failed: %s", err)
	}

	if status, err := js.Get("status").Int(); err == nil {
		if 0 == status {
			ele := js.Get("result").Get("elements").GetIndex(0)
			if distance, err := ele.Get("distance").Int(); err == nil {
				return distance, nil
			}
		}
	}
	logErr("http response data invalid:", string(cnt))
	return -1, fmt.Errorf("http response data invalid")
}

// Save order info to DB, and return order_id
func SaveOrder2DB(order Order, distance int) (int64, error) {
	stmt, err := db.Prepare(
		"INSERT INTO t_orders SET origin=?,destination=?,distance=?")
	if err != nil {
		logErr("db prepare failed:", err)
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(strings.Join(order.Origin, ","),
		strings.Join(order.Destination, ","), distance)
	if err != nil {
		logErr("db exec failed:", err)
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		logErr("db get LastInsertId failed:", err)
		return -1, err
	}
	return id, nil
}

func TakeOrder(orderid string) error {
	stmt, err := db.Prepare(`SELECT id,status FROM t_orders WHERE id=? AND status=?`)
	if err != nil {
		logErr("db prepare failed:", err)
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(orderid, orderUnAssign)
	if err != nil {
		logErr("db query failed:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		stmt, err := db.Prepare("UPDATE t_orders SET status=? WHERE id=? AND status=?")
		if err != nil {
			logErr("db prepare failed:", err)
			return err
		}

		res, err := stmt.Exec(orderTaken, orderid, orderUnAssign)
		if err != nil {
			logErr("db exec failed:", err)
			return err
		}
		if affected, err := res.RowsAffected(); err != nil || 0 == affected {
			logErr("rows affected invalid:", affected)
			return err
		}
		return nil
	}
	return fmt.Errorf("order requested not exists or has been taken")
}

func ListOrder(page, limit string) *[]OrderInfo {
	orderList := []OrderInfo{}
	for {
		stmt, err := db.Prepare(`SELECT id,distance,status FROM t_orders LIMIT ? OFFSET ?`)
		if err != nil {
			logErr("db prepare failed:", err)
			break
		}
		defer stmt.Close()

		n_page, err := strconv.Atoi(page)
		if err != nil {
			n_page = 0
		}
		n_limit, err := strconv.Atoi(limit)
		if err != nil {
			n_limit = pageLen
		}
		logInfo("page:", n_page, "limit:", n_limit)
		rows, err := stmt.Query(n_limit, n_page*pageLen)
		if err != nil {
			logErr("db query failed:", err)
			break
		}
		defer rows.Close()

		for rows.Next() {
			orderInfo := OrderInfo{}
			status := 0
			rows.Scan(&orderInfo.Id, &orderInfo.Distance, &status)
			if orderUnAssign == status {
				orderInfo.Status = "UNASSIGN"
			} else {
				orderInfo.Status = "TAKEN"
			}
			orderList = append(orderList, orderInfo)
		}
		break
	}
	return &orderList
}

func NewOrderHandler(w http.ResponseWriter, r *http.Request) {
	defer Common.CheckPanic()
	logInfo(fmt.Sprintf("client: %s, url: %s, method: %s\n", r.RemoteAddr, r.RequestURI, r.Method))
	for {
		if "POST" != r.Method {
			logErr("Method should be POST:", r.Method)
			break
		}

		cnt, err := ioutil.ReadAll(r.Body)
		if err != nil || len(cnt) != int(r.ContentLength) {
			logErr("content length invalid", "len:", len(cnt), "ContentLength:", r.ContentLength)
			break
		}

		order := Order{}
		err = json.Unmarshal(cnt, &order)
		if err != nil || order.Origin == nil || order.Destination == nil {
			logErr("Request data invalid:", string(cnt))
			break
		}

		distance, err := GetDistance(order)
		if err != nil {
			logErr("Get Distance failed")
			break
		}

		order_id, err := SaveOrder2DB(order, distance)
		if err != nil {
			logErr("Save order failed:", err)
			break
		}
		m_body := map[string]interface{}{}
		m_body["id"] = order_id
		m_body["distance"] = distance
		m_body["status"] = "UNASSIGN"
		data, err := json.MarshalIndent(m_body, "", "  ")
		if err != nil {
			logErr("Marshal data failed:", err)
			break
		}
		logInfo("Order Success, ret:", string(data))
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("{\n    \"error\":\"ERROR_DESCRIPTION\"\n}"))
	logErr("Order failed")
}

func TakeOrderHandler(w http.ResponseWriter, r *http.Request) {
	defer Common.CheckPanic()
	logInfo("client:", r.RemoteAddr, "url:", r.RequestURI, "method:", r.Method)
	for {
		if "PUT" != r.Method {
			logErr("Method should be PUT:", r.Method)
			break
		}

		cnt, err := ioutil.ReadAll(r.Body)
		if err != nil || len(cnt) != int(r.ContentLength) {
			logErr("content length invalid", "len:", len(cnt), "ContentLength:", r.ContentLength)
			break
		}

		r_body := map[string]string{}
		err = json.Unmarshal(cnt, &r_body)
		if "taken" != r_body["status"] || err != nil {
			logErr("Request data invalid:", string(cnt))
			break
		}

		orderId := path.Base(r.URL.Path)
		if err := TakeOrder(orderId); err != nil {
			logErr("TakeOrder failed:", err)
			break
		}
		w.Write([]byte("{\n    \"status\": \"SUCCESS\"\n}"))
		logInfo("Take order success")
		return
	}
	w.WriteHeader(http.StatusConflict)
	w.Write([]byte("{\n    \"error\":\"ORDER_ALREADY_BEEN_TAKEN\"\n}"))
	logErr("Take order failed")
}

func ListOrderHandler(w http.ResponseWriter, r *http.Request) {
	defer Common.CheckPanic()
	logInfo("client:", r.RemoteAddr, "url:", r.RequestURI, "method:", r.Method)

	for {
		if "GET" != r.Method {
			logErr("Method should be PUT:", r.Method)
			break
		}
		page, limit := r.FormValue("page"), r.FormValue("limit")
		if page == "" {
			page = "0"
		}
		if limit == "" {
			limit = "1"
		}
		logInfo("page:", page, "limit:", limit)
		res_data := ListOrder(page, limit)
		data, err := json.MarshalIndent(res_data, "", "  ")
		if err != nil {
			logErr("Json Marshal data failed:", err)
			break
		}
		logInfo("List Order Success, response:", string(data))
		w.Write(data)
		return
	}
	w.Write([]byte("[]"))
}

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	var err error
	db, err = sql.Open("mysql", "root:abc123456@tcp(mysql:3306)/orderserver?charset=utf8")
	if err != nil {
		logErr("sql open failed:", err)
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(100)

	logInfo("server listen:", port, " begin ...")

	http.HandleFunc("/order", NewOrderHandler)
	http.HandleFunc("/order/", TakeOrderHandler)
	http.HandleFunc("/orders", ListOrderHandler)
	logErr(http.ListenAndServe(port, nil))
}
