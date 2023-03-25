package bybit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order-book/pkg/types"
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

	l2 := NewL2OrderBook()
	s.l2BySymbol[snapshot.Params.Symbol] = l2

	for _, item := range snapshot.Payload.Snapshot {
		side := types.SideFromString(item.Side)
		tm := time.Unix(0, item.Timestamp)
		l2.Apply(item.Price, side, item.Volume, tm)
	}
	return nil
}

func (s *Source) onUpdate(data []byte) error { return nil }
