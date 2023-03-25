package bybit

import "github.com/shopspring/decimal"

//  {"req_id":"test", "op":"subscribe", "args":["orderbook.1.BTCUSDT"]}

func GetWsOrderBookSubscribeRequest() interface{} {
	subsribeReq := struct {
		ReqId     string
		Operation string
		Args      []string
	}{}

	subsribeReq.ReqId = "test"
	subsribeReq.Operation = "subscribe"
	subsribeReq.Args = append(subsribeReq.Args, "orderbook.1.BTCUSDT")

	return subsribeReq
}

type WsL2Item struct {
	Price     decimal.Decimal `json:"price"`
	Volume    uint64          `json:"volume,string"`
	Side      string          `json:"side"`
	Timestamp int64           `json:"timestamp"`
}

type WsL2Update struct {
	Params struct {
		Symbol string `json:"symbol"`
	} `json:"params"`
	Payload WsL2Item `json:"payload"`
}

type WsL2Snapshot struct {
	Params struct {
		Symbol string `json:"symbol"`
	} `json:"params"`
	Payload struct {
		Snapshot []*WsL2Item `json:"snapshot"`
		Updates  []*WsL2Item `json:"updates"`
	} `json:"payload"`
}
