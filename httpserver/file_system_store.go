package httpserver

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type FileSystemPlayerStore struct {
	Database *json.Encoder
	league League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialisePlayerDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem with initilising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		json.NewEncoder(&tape{file}),
		league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i int, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})


	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(playerName string) int {

	player, _ := f.league.Find(playerName)

	if player != nil {
		return player.Wins
	}
	return 0

}

func (f *FileSystemPlayerStore) RecordWin(playerName string) {
	player, _ := f.league.Find(playerName)

	if player != nil {
		player.Wins++
	}else{
		f.league = append(f.league, Player{playerName, 1})
	}
	f.Database.Encode(f.league)
}

func (f *FileSystemPlayerStore) RecordNewPlayer(player Player){

	playerFound, idx := f.league.Find(player.Name)
	if playerFound != nil{
		//player already exists- write over
		f.league = append(f.league[:idx], f.league[idx+1:]...)
	}
	f.league = append(f.league, Player{player.Name, player.Wins})
	f.Database.Encode(f.league)
}

func (f *FileSystemPlayerStore) DeletePlayer(name string){
	playerFound, idx := f.league.Find(name)
	if playerFound != nil{
		//player found- delete
		f.league = append(f.league[:idx], f.league[idx+1:]...)
	}else{
		fmt.Errorf("player cound not be found, and deleted: %s", name)
	}
	f.Database.Encode(f.league)

}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err )
	}

	// Path for an empty file (make an empty JSON for the rest of my code)
	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}
	return nil
}

