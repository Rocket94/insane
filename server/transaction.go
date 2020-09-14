package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"io"
	"github.com/donnie4w/go-logger/logger"
	"time"
	"math/rand"
	"strconv"
	"insane/utils"
)

const (
	NEW_ORDER_rq    = 1
	PAYMENT_rq      = 2
	ORDER_STATUS_rq = 3
	DELIVERY_rq     = 4
	STOCK_LEVEL_rq  = 5
)

type neworder struct {
	Wid          int    `json:"wid"`
	Did          int    `json:"did"`
	Cid          int    `json:"cid"`
	Olcnt        int    `json:"olcnt"`
	Oentryd      string `json:"oentryd"`
	Oalllocal    int    `json:"oalllocal"`
	Olquantities []int  `json:"olquantities"`
	Oliids       []int  `json:"oliids"`
	Olsupplywids []int  `json:"olsupplywids"`
}
type payment struct {
	Wid     int     `json:"wid"`
	Did     int     `json:"did"`
	Cwid    int     `json:"cwid"`
	Cdid    int     `json:"cdid"`
	Cid     int     `json:"cid"` //如果根据clast检索，cid设置为-12345
	Clast   string  `json:"clast"`
	Hamount float64 `json:"hamount"`
	Hdate   string  `json:"hdate"` //格式："2020-09-10 14:25:47"
}
type orderstatus struct {
	Cwid   int    `json:"cwid"`
	Cdid   int    `json:"cdid"`
	Cid    int    `json:"cid"`
	Clast  string `json:"clast"`
	Byname bool   `json:"byname"` //如果byname为true，cid可以设置为任意值，事务中会被重新赋值
}
type delivery struct {
	Wid        int `json:"wid"`
	Ocarrierid int `json:"ocarrierid"`
}
type stocklevel struct {
	Wid       int `json:"wid"`
	Did       int `json:"did"`
	Threshold int `json:"threshold"`
}

func CreateBody(TxnType int, r *rand.Rand) io.Reader {
	var body interface{}
	switch TxnType {
	case NEW_ORDER_rq:
		wid := utils.GetRandomIntRange(1, 10, r)
		cnt := utils.GetRandomIntRange(5, 15, r)
		supply, alllocal := utils.GetOlsupplywidsRandom(r, cnt, wid)
		body = neworder{
			Wid:          wid,
			Did:          utils.GetRandomIntRange(1, 10, r),
			Cid:          utils.GetCIDRandom(r),
			Olcnt:        cnt,
			Oentryd:      time.Now().Format("2020-09-10 14:25:47"),
			Oalllocal:    alllocal,
			Olquantities: utils.GetOlquantitiesRandom(r, cnt),
			Oliids:       utils.GetOliidsRandom(r, cnt),
			Olsupplywids: supply,
		}
		break
	case PAYMENT_rq:
		wid := utils.GetRandomIntRange(1, 10, r)
		did := utils.GetRandomIntRange(1, 10, r)
		cwid, cdid := utils.GetCwdidRandom(r, wid, did)
		body = payment{
			Wid:     wid,
			Did:     did,
			Cwid:    0,
			Cdid:    0,
			Cid:     0,
			Clast:   "",
			Hamount: float64(utils.GetRandomIntRange(100, 500000, r) / 100),
			Hdate:   time.Now().Format("2020-09-10 14:25:47"),
		}
		break
	case ORDER_STATUS_rq:
		body = orderstatus{
			Cwid:   0,
			Cdid:   0,
			Cid:    0,
			Clast:  "",
			Byname: false,
		}
		break
	case DELIVERY_rq:
		body = delivery{
			Wid:        0,
			Ocarrierid: 0,
		}
		break
	case STOCK_LEVEL_rq:
		body = stocklevel{
			Wid:       0,
			Did:       0,
			Threshold: 0,
		}
		break
	default:
		body = nil
	}

	s, err := json.Marshal(body)
	if err != nil || s == nil {
		logger.Debug("request body is nil ", err)
		s = nil
	}
	return strings.NewReader(string(s))
}

func SetRequest(request *Request) (req *http.Request, transactionType int, err error) {
	r := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	tt := rand.Intn(100)
	if tt >= 0 && tt <= 44 {
		tt = 1
	} else if tt >= 45 && tt <= 87 {
		tt = 2
	} else if tt >= 88 && tt <= 91 {
		tt = 3
	} else if tt >= 92 && tt <= 95 {
		tt = 4
	} else {
		tt = 5
	}
	body := CreateBody(tt, r)
	//tommorow look why
	req, err = http.NewRequest("POST", request.Url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transaction-Type", strconv.FormatInt(int64(tt), 10))
	if err != nil {
		return nil, tt, err
	}
	return req, tt, nil
}
