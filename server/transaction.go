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
)

const (
	NEW_ORDER_rq    = 1
	PAYMENT_rq      = 2
	DELIVERY_rq     = 3
	ORDER_STATUS_rq = 4
	STOCK_LEVEL_rq  = 5
)

type NewOrder struct {
	WID         int
	DID         int
	CID         int
	OLcnt       int
	OQuantity   int
	OEntryD     string
	OLIID       int
	OLSupplyWID int
	IsHome      bool
}
type Payment struct {
	WID     int
	DID     byte
	CDID    byte
	IsHome  bool
	CLast   int
	CID     int
	HAmount float32
	HDate   string
}
type OrderStatus struct {
	WID          int
	DID          int
	OLDeliveryID int
}
type Delivery struct {
	WID         int
	OCarrierID  int
	OLDeliveryD string
}
type StockLevel struct {
	WID       int
	DID       int
	Threshold int
}

func CreateBody(TxnType int) io.Reader {
	var body interface{}
	switch TxnType {
	case NEW_ORDER_rq:
		body = NewOrder{
			WID:         0,
			DID:         0,
			CID:         0,
			OLcnt:       0,
			OQuantity:   0,
			OEntryD:     "",
			OLIID:       0,
			OLSupplyWID: 0,
			IsHome:      false,
		}
		break
	case PAYMENT_rq:
		body = Payment{
			WID:     0,
			DID:     0,
			CDID:    0,
			IsHome:  false,
			CLast:   0,
			CID:     0,
			HAmount: 0,
			HDate:   "",
		}
		break
	case DELIVERY_rq:
		body = Delivery{
			WID:         0,
			OCarrierID:  0,
			OLDeliveryD: "",
		}
		break
	case ORDER_STATUS_rq:
		body = OrderStatus{
			WID:          0,
			DID:          0,
			OLDeliveryID: 0,
		}
		break
	case STOCK_LEVEL_rq:
		body = StockLevel{
			WID:       0,
			DID:       0,
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

func SetRequest(request *Request) (req *http.Request, err error) {
	rand.Seed(time.Now().UnixNano())
	tt := rand.Intn(5)
	body := CreateBody(tt + 1)
	//tommorow look why
	req, err = http.NewRequest(request.Method, request.Url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("transaction-Type", strconv.FormatInt(int64(tt),10))
	if err != nil {
		return nil, err
	}
	return req, nil
}
