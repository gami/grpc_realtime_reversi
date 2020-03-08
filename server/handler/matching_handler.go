package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"reversi/build"
	"reversi/game"
	"reversi/gen/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MatchingHandler struct {
	sync.RWMutex
	Rooms       map[int32]*game.Room
	maxPlayerID int32
}

func NewMatchingHandler() *MatchingHandler {
	return &MatchingHandler{
		Rooms: make(map[int32]*game.Room),
	}
}

func (h *MatchingHandler) JoinRoom(req *pb.JoinRoomRequest, stream pb.MatchingService_JoinRoomServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), 2*time.Minute)
	defer cancel()

	h.Lock()

	// プレイヤーの新規作成
	h.maxPlayerID++
	me := &game.Player{
		ID: h.maxPlayerID,
	}

	// 空いている部屋を探す
	for _, room := range h.Rooms {
		if room.Guest == nil {
			me.Color = game.White
			room.Guest = me
			stream.Send(&pb.JoinRoomResponse{
				Status: pb.JoinRoomResponse_MATCHED,
				Room:   build.PBRoom(room),
				Me:     build.PBPlayer(room.Guest),
			})
			h.Unlock()
			fmt.Printf("matched room_id=%v\n", room.ID)
			return nil
		}
	}

	// 空いている部屋がなかったので部屋を作る
	me.Color = game.Black
	room := &game.Room{
		ID:   int32(len(h.Rooms)) + 1,
		Host: me,
	}
	h.Rooms[room.ID] = room
	h.Unlock()

	stream.Send(&pb.JoinRoomResponse{
		Status: pb.JoinRoomResponse_WAITTING,
		Room:   build.PBRoom(room),
	})

	ch := make(chan int)
	go func(ch chan<- int) {
		for {
			h.RLock()
			guest := room.Guest
			h.RUnlock()

			if guest != nil {
				stream.Send(&pb.JoinRoomResponse{
					Status: pb.JoinRoomResponse_MATCHED,
					Room:   build.PBRoom(room),
					Me:     build.PBPlayer(room.Host),
				})
				ch <- 0
				break
			}
			time.Sleep(1 * time.Second)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ch)

	select {
	case <-ch:
	case <-ctx.Done():
		return status.Errorf(codes.DeadlineExceeded, "マッチングできませんでした")
	}

	return nil
}
