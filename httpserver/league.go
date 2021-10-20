package httpserver

import (
	"encoding/json"
	"fmt"
	"os"
)

type League []Player

func(l League) Find(name string) (*Player, int) {
	for idx, player := range l {
		if player.Name ==name {
			return &l[idx], idx
		}
	}
	return nil, -2
}


func NewLeague(rdr *os.File) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)

	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}
