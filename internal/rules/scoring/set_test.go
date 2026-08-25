package scoring

import (
	"testing"

	"github.com/gabriel-d0/rummy_backend/internal/rules/meld"
	"github.com/gabriel-d0/rummy_backend/internal/rules/tile"
)

func TestScoreSetSimple(t *testing.T) {
	// 7 red, 7 yellow, 7 blue => 5+5+5=15
	m, _ := meld.New("m1", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 7),
		tile.MustTile("t2", tile.Yellow, 7),
		tile.MustTile("t3", tile.Blue, 7),
	}, nil)
	if v, _ := ScoreSet(m); v != 15 {
		t.Fatalf("7 set 15 got %d", v)
	}
	// 10 red, yellow, blue, black => 10*4=40
	m2, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 10),
		tile.MustTile("t2", tile.Yellow, 10),
		tile.MustTile("t3", tile.Blue, 10),
		tile.MustTile("t4", tile.Black, 10),
	}, nil)
	if v, _ := ScoreSet(m2); v != 40 {
		t.Fatalf("10 set 40 got %d", v)
	}
}

func TestScoreSetAce(t *testing.T) {
	// Ace set of 3 => 25*3=75
	m, _ := meld.New("m1", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Yellow, 1),
		tile.MustTile("t3", tile.Blue, 1),
	}, nil)
	if v, _ := ScoreSet(m); v != 75 {
		t.Fatalf("Ace set 75 got %d", v)
	}
	// Ace set of 4 => 25*4=100
	m2, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Yellow, 1),
		tile.MustTile("t3", tile.Blue, 1),
		tile.MustTile("t4", tile.Black, 1),
	}, nil)
	if v, _ := ScoreSet(m2); v != 100 {
		t.Fatalf("Ace set 4 100 got %d", v)
	}
}

func TestScoreSetJoker(t *testing.T) {
	// 5 red, 5 yellow + joker as 5 blue => 5+5+5=15
	j := tile.MustJoker("j1")
	m, _ := meld.New("m1", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 5),
		tile.MustTile("t2", tile.Yellow, 5),
		j,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Blue, 5)})
	if v, _ := ScoreSet(m); v != 15 {
		t.Fatalf("joker 5 set 15 got %d", v)
	}
	// Joker as Ace in Ace set => 25
	j2 := tile.MustJoker("j1")
	m2, _ := meld.New("m2", meld.KindSet, []tile.TileInstance{
		tile.MustTile("t1", tile.Red, 1),
		tile.MustTile("t2", tile.Yellow, 1),
		j2,
	}, map[tile.TileInstanceId]tile.TileInstance{"j1": tile.MustTile("rep", tile.Blue, 1)})
	if v, _ := ScoreSet(m2); v != 75 {
		t.Fatalf("joker Ace set 75 got %d", v)
	}
}
