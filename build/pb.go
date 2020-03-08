package build

import (
	"reversi/game"
	"reversi/gen/pb"
)

func PBRoom(r *game.Room) *pb.Room {
	return &pb.Room{
		Id:    r.ID,
		Host:  PBPlayer(r.Host),
		Guest: PBPlayer(r.Guest),
	}
}

func PBPlayer(p *game.Player) *pb.Player {
	if p == nil {
		return nil
	}
	return &pb.Player{
		Id:    p.ID,
		Color: PBColor(p.Color),
	}
}

func PBColor(c game.Color) pb.Color {
	switch c {
	case game.Black:
		return pb.Color_BLACK
	case game.White:
		return pb.Color_WHITE
	case game.Empty:
		return pb.Color_EMPTY
	case game.Wall:
		return pb.Color_WALL
	}

	return pb.Color_UNKNOWN
}

func PBBoard(b *game.Board) *pb.Board {
	pbCols := make([]*pb.Board_Col, 0, 10)

	for _, col := range b.Cells {
		pbCells := make([]pb.Color, 0, 10)
		for _, c := range col {
			pbCells = append(pbCells, PBColor(c))
		}
		pbCols = append(pbCols, &pb.Board_Col{
			Cells: pbCells,
		})
	}

	return &pb.Board{
		Cols: pbCols,
	}
}
