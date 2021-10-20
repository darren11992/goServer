package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)



// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
	RecordNewPlayer(player Player)
	DeletePlayer(name string)
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	Store PlayerStore
	//http.Handler // Embedding - "PlayerServer" now has all the methods that http.handler has (ServeHTTP)
	http.Server
	// This is referenced with p.Handler in "NewPlayerServer"
}

// Player stores a name with a number of wins
type Player struct {
	Name string
	Wins int
}

type User struct {
	Username string
	Password string
}

const jsonContentType = "application/json"

var goodUsernames = map[string]string{"user_a": "passwordA", "user_b": "passwordB", "user_c": "passwordC", "admin": "Password1"}

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.Store = store
	router := http.NewServeMux()
	p.Addr = ":5000"
	router.Handle("/list", http.HandlerFunc(p.listHandler))
	router.Handle("/store/", http.HandlerFunc(p.playersHandler))
	router.Handle("/login", http.HandlerFunc(p.loginHandler))
	router.Handle("/ping", http.HandlerFunc(p.pingHandler))
	router.Handle("/shutdown", http.HandlerFunc(p.shutdownHandler))


	p.Handler = router // Can do this because NewServeMux has the method ServeHTTP
	return p
}

	func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL, r.RemoteAddr)
		player := strings.TrimPrefix(r.URL.Path, "/store/")

		switch r.Method {
		case http.MethodPost:
			p.processWin(w, player)
		case http.MethodGet:
			p.showScore(w, player)
		case http.MethodPut:
			p.processNewPlayer(w, r)
		case http.MethodDelete:
			p.processDelete(w, player)
		}
}

//func loggin(logger *log.Logger)

func (p *PlayerServer) listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL, r.RemoteAddr)
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.Store.GetLeague())

	//w.WriteHeader(http.StatusOK)

}

func (p*PlayerServer) loginHandler(w http.ResponseWriter, r *http.Request){
	//receivedUsername := r.Header.Get("Authorization")
	log.Println(r.Method, r.URL, r.RemoteAddr)
	//var receivedUsers []User
	//decoder := json.NewDecoder(r.Body)
	//err := decoder.Decode(&receivedUsers)
	u, pass, ok := r.BasicAuth()

	if !ok {
		fmt.Println("Error parsing basic aaauth")
		w.WriteHeader(401)
		return
	}

	validPass, userValid := goodUsernames[u]


	if !userValid{
		fmt.Printf("Username provided is incorrect: %s\n", u)
		w.WriteHeader(401)
		return
	}
	if pass != validPass {
		fmt.Printf("Password provided is incorrect: %s\n", pass)
		w.WriteHeader(401)
		return
	}
	//fmt.Printf("Username: %s\n", u)
	//fmt.Printf("Password: %s\n", p)
	w.WriteHeader(200)
	//r.Header.Set("Authorization", "Bearer gntkgmvjd")
	w.Header().Set("Authorization", "Bearer gntkgmvjd")
	w.Write([]byte("Bearer jfjklnbrepieb"))
	return
}

func (p *PlayerServer) pingHandler(w http.ResponseWriter, r *http.Request){
	log.Println(r.Method, r.URL, r.RemoteAddr)
	fmt.Fprint(w, "pong")
}

func (p *PlayerServer) shutdownHandler(w http.ResponseWriter, r *http.Request){
	log.Println(r.Method, r.URL, r.RemoteAddr)

	//Goroutine this
	err := p.Shutdown(context.Background())
	if err != nil{
		fmt.Printf("There was an error in shutdown")
	}
}

func (p *PlayerServer) getLeagueTable() []Player {
	return p.Store.GetLeague()
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string){
	score := p.Store.GetPlayerScore(player)

	if score == 0 {
		w.Write([]byte("404 Not Found"))
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.Store.RecordWin(player) // First value stored in the spy is the name of the player
	w.WriteHeader(http.StatusAccepted)
}

func (p* PlayerServer) processNewPlayer(w http.ResponseWriter, r *http.Request){
	var requestPlayer []Player
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestPlayer)
	
	if err != nil {
		fmt.Print("There was an error when decoding the Body of a PUT request", err)
	}
	
	for _, player := range requestPlayer {
		p.Store.RecordNewPlayer(player)
	}
		w.WriteHeader(http.StatusAccepted)
}

func (p* PlayerServer) processDelete(w http.ResponseWriter, player string){
	p.Store.DeletePlayer(player)
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusAccepted)


}