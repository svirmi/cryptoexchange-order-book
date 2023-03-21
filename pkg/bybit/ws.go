package bybit

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
