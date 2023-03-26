package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order-book/pkg/types"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	bybitWsUrl = "wss://stream.bybit.com/v5/public/spot"

	defaultRetryInterval = time.Second
	defaultReadTimeout   = time.Second * 5
	defaultWriteTimeout  = time.Second * 5
)

type Source struct {
	sync.RWMutex
	l2BySymbol map[string]*L2OrderBook
}

func NewSource() *Source {
	return &Source{
		l2BySymbol: make(map[string]*L2OrderBook),
	}
}

func (s *Source) Start(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			log.Printf("source stopped: context cancelled")
			return

		default:
			if err := s.receiveData(ctx); err != nil {
				log.Printf("error recieveing data: %v", err)
			}

			// TODO : refactor retry connection code
			log.Printf("sleep for %v ", defaultRetryInterval)
			time.Sleep(defaultRetryInterval)
		}
	}
}

func (s *Source) receiveData(ctx context.Context) error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(ctx, bybitWsUrl, nil)
	if err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout)); err != nil {
		return err
	}

	if err := conn.WriteJSON(GetWsOrderBookSubscribeRequest()); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default: // read from WS connection
			mt, data, err := conn.ReadMessage()
			if err != nil {
				return err
			}
			if mt != websocket.TextMessage {
				return fmt.Errorf("wrong message type")
			}
			header := struct{ Type string }{}
			if err := json.Unmarshal(data, &header); err != nil {
				return err
			}

			switch header.Type {
			case "snapshot":
				s.onSnapshot(data)
			case "update":
				s.onUpdate(data)
			default:
			}
		}
	}
}

func (s *Source) onSnapshot(data []byte) error {
	var snapshot WsL2Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	l2 := NewL2OrderBook()
	s.l2BySymbol[snapshot.Params.Symbol] = l2

	for _, items := range [][]*WsL2Item{
		snapshot.Payload.Snapshot,
		snapshot.Payload.Updates,
	} {
		for _, item := range items {
			side := types.SideFromString(item.Side)
			tm := time.Unix(0, item.Timestamp)
			l2.Apply(item.Price, side, item.Volume, tm)
		}
	}

	log.Printf("snapshot applied: symbol=%s, bid=%d, ask=%d",
		snapshot.Params.Symbol, l2.bid.Len(), l2.ask.Len())

	return nil
}

func (s *Source) onUpdate(data []byte) error {
	var update WsL2Update
	if err := json.Unmarshal(data, &update); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	l2, ok := s.l2BySymbol[update.Params.Symbol]
	if !ok {
		log.Printf("inconsistent update for symbol %s", update.Params.Symbol)
		return nil
	}

	side := types.SideFromString(update.Payload.Side)
	tm := time.Unix(0, update.Payload.Timestamp)
	l2.Apply(update.Payload.Price, side, update.Payload.Volume, tm)

	log.Printf("update applied: symbol=%s, bid=%d, ask=%d",
		update.Params.Symbol, l2.bid.Len(), l2.ask.Len())

	return nil
}

func (s *Source) GetL2OrderBook(symbol string, size int) (types.L2OrderBook, error) {
	s.RLock()
	defer s.RUnlock()

	l2, ok := s.l2BySymbol[symbol]
	if !ok {
		return types.L2OrderBook{}, fmt.Errorf("no data for symbol %s", symbol)
	}

	return types.L2OrderBook{
		Bid: l2.GetBid(size),
		Ask: l2.GetAsk(size),
	}, nil
}
