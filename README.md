# cryptoexchange-order-book
#### Inspired by https://youtu.be/OsXjDZ52dyQ video
##### Code [https://github.com/forward32/quotes](https://github.com/forward32/quotes)

##### Documentation:
- [https://www.bybit.com](https://www.bybit.com)
    - [WS](https://bybit-exchange.github.io/docs/v5/websocket/public/orderbook)
    - [How to Subscribe to Topics](https://bybit-exchange.github.io/docs/v5/ws/connect#how-to-subscribe-to-topics)


    
- [blocktrade.com (not implemented)](https://trade.blocktrade.com/api_documentation#blocktrade-api)
    - [trading assets](https://trade.blocktrade.com/api_documentation#trading-assets)
    - [trading pairs](https://trade.blocktrade.com/api_documentation#trading-pairs)
    - [WS](https://trade.blocktrade.com/api_documentation#order-book--websockets-)

##### Some useful code examples
- [CCXW : WebSocket client for 38 cryptocurrency exchanges](https://github.com/altangent/ccxws)
- [Transpiled version of the CCXT library to Golang; support 100+ cryptocurrency exchanges with a unified API](https://github.com/prompt-cash/ccxt-go)
- [https://github.com/trever-io/go-blocktrade](https://github.com/trever-io/go-blocktrade)
- [https://github.com/johnwashburne/Crypto-Price-Aggregator](https://github.com/johnwashburne/Crypto-Price-Aggregator)

##### Reading list:
- [CCXT Crypto Exchange Order Book (Python)](https://blog.shrimpy.io/blog/ccxt-crypto-exchange-order-book-snapshot)

##### To run test
```bash
/cryptoexchange-order-book/pkg/bybit$ go test -v -race
```

