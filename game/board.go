package game

import "fmt"

type Board struct {
	Cells [][]Color
}

// 盤面を作成します
func NewBoard() *Board {
	// 8x8のマス+壁で、10x10の盤面を二次元配列で作ります
	b := &Board{
		Cells: make([][]Color, 10),
	}
	for i := 0; i < 10; i++ {
		b.Cells[i] = make([]Color, 10)
	}

	// 盤面の端に壁を設置します
	for i := 0; i < 10; i++ {
		b.Cells[0][i] = Wall
	}
	for i := 1; i < 9; i++ {
		b.Cells[i][0] = Wall
		b.Cells[i][9] = Wall
	}
	for i := 0; i < 9; i++ {
		b.Cells[9][i] = Wall
	}

	// 初期石を置きます
	b.Cells[4][4] = White
	b.Cells[5][5] = White
	b.Cells[5][4] = Black
	b.Cells[4][5] = Black

	return b
}

// 石を置きます
func (b *Board) PutStone(x int32, y int32, c Color) error {
	// そのマスに石を置けるかチェックします
	if !b.CanPutStone(x, y, c) {
		return fmt.Errorf("Can not put stone x=%v, y=%v color=%v", x, y, ColorToStr(c))
	}

	// マスに石を置きます
	b.Cells[x][y] = c

	// 置いた石の縦/横/斜めの各方向でひっくり返すことのできる石を全てひっくり返します
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			b.TurnStonesByDirection(x, y, c, int32(dx), int32(dy))
		}
	}

	return nil
}

// マスに石を置けるか判定します
func (b *Board) CanPutStone(x int32, y int32, c Color) bool {
	// すでに石がある場合は石を置けません
	if b.Cells[x][y] != Empty {
		return false
	}

	// 縦/横/斜めの各方向をチェックします
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			// ひっくり返すことのできる石がひとつでもあれば、石を置けます
			if b.CountTurnableStonesByDirection(x, y, c, int32(dx), int32(dy)) > 0 {
				return true
			}
		}
	}

	// ひとつもひっくり返すことできる石がなければ、石を置けません
	return false
}

// あるマスに石を置いた場合、ある方向にひっくり返すことのできる石がいくつあるか数えます
func (b *Board) CountTurnableStonesByDirection(x int32, y int32, c Color, dx int32, dy int32) int {
	cnt := 0

	nx := x + dx
	ny := y + dy
	for {
		nc := b.Cells[nx][ny]

		// 壁か自分の石であればループを終了する
		if nc != OpponentColor(c) {
			break
		}

		//相手の石なので数え上げ
		cnt++

		nx += dx
		ny += dy
	}

	// その方向にある相手の石の数がゼロより大きく、かつその先に自分の石がある場合は数を返す
	if cnt > 0 && b.Cells[nx][ny] == c {
		return cnt
	}

	// それ以外の場合はゼロ
	return 0
}

// ある方向の石をひっくり返す。ひっくり返しても良い場合だけ呼ぶ。
func (b *Board) TurnStonesByDirection(x int32, y int32, c Color, dx int32, dy int32) {
	nx := x + dx
	ny := y + dy

	for {
		nc := b.Cells[nx][ny]

		if nc != OpponentColor(c) {
			break
		}

		b.Cells[nx][ny] = c

		nx += dx
		ny += dy
	}
}

// 盤面内である色の石を置けるマスの数を数える
func (b *Board) AvailableCellCount(c Color) int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.CanPutStone(int32(i), int32(j), c) {
				cnt++
			}
		}
	}

	return cnt
}

// 盤面内に置かれている石の数を数える
func (b *Board) Score(c Color) int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.Cells[i][j] != c {
				continue
			}
			cnt++
		}
	}

	return cnt
}

// 盤面内で石が置かれてないマスの数を数える
func (b *Board) Rest() int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.Cells[i][j] == Empty {
				cnt++
			}
		}
	}

	return cnt
}
