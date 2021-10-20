package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
	NewUserCalls []string
	league League
	DeleteCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int{
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string){
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) RecordNewPlayer(player Player){
	s.NewUserCalls = append(s.NewUserCalls, player.Name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func(s *StubPlayerStore) DeletePlayer(name string){
	s.DeleteCalls = append(s.DeleteCalls, name)
}

func newStore(scores map[string]int) StubPlayerStore{
	store := StubPlayerStore{scores, nil, nil, nil, nil}
	return store
}

func TestGETPlayers(t *testing.T){
	store := newStore(map[string]int{"Pepper": 20, "Floyd": 10})
	server := NewPlayerServer(&store)

	t.Run("returns Pepper's score", func(t *testing.T){
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request) // Implements handler interface

		responseCodeGot := response.Code
		responseCodeWant := http.StatusOK
		assertStatus(t, responseCodeGot, responseCodeWant)

		got := response.Body.String()
		want:= "20"
		assertResponseBody(t, got, want)

	})
	t.Run("Returns Floyd's Score", func(t *testing.T){
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		responseCodeGot := response.Code
		responseCodeWant := http.StatusOK
		assertStatus(t, responseCodeGot, responseCodeWant)

		got := response.Body.String()
		want:= "10"
		assertResponseBody(t, got, want)
	})
	t.Run("Return response for missing player", func(t *testing.T){
		request := newGetScoreRequest("Potato")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Code
		want:= http.StatusNotFound


		assertStatus(t, got, want)
	})
}

func TestScoreWins(t *testing.T){

	t.Run("We get a good status code from a POST", func(t *testing.T){
		store := newStore(map[string]int{})
		//Store has to be local to the individual test otherwise each affects the subsequent winCalls check...
		server := NewPlayerServer(&store)
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1{
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
	})

	t.Run("it records wins on POST", func(t *testing.T){
		store := newStore(map[string]int{})
		server := NewPlayerServer(&store)
		player := "Pepper"

		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})

	t.Run("We get a good status from a PUT", func(t *testing.T) {

		store := newStore(map[string]int{})
		server := NewPlayerServer(&store)
		newPlayer := Player{"Potato", 10}
		jsonPlayer, err := json.Marshal(newPlayer)
		if err != nil {
			t.Errorf("Error when converting PUT data to Json: %s", err)
		}

		request := newPutPlayerRequest(newPlayer.Name, jsonPlayer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("We record new players with a set win count from a PUT", func(t *testing.T){
		store := newStore(map[string]int{})
		server := NewPlayerServer(&store)
		newPlayer := Player{"Potato", 10}
		jsonPlayer, err := json.Marshal(newPlayer)
		if err != nil {
			t.Errorf("Error when converting PUT data to Json: %s", err)
		}

		//request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/store/%s", newPlayer.Name), bytes.NewBuffer(jsonPlayer))
		request := newPutPlayerRequest(newPlayer.Name, jsonPlayer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.NewUserCalls) != 1 {
			t.Fatalf("got %d calls to RecordNewUser want %d", len(store.NewUserCalls), 1)
		}

		if store.NewUserCalls[0] != newPlayer.Name {
			t.Errorf("did not store correct winner got %q want %q", store.NewUserCalls[0], newPlayer.Name)
		}
	})
	t.Run("We get a good response from a DELETE", func(t *testing.T){
		store := newStore(map[string]int{"Pepper": 20, "Floyd": 10})
		server := NewPlayerServer(&store)

		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/store/%s", "Pepper"), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
	t.Run("We remove players in URL of a DELETE request", func(t *testing.T){
		store := newStore(map[string]int{"Pepper": 20, "Floyd": 10})
		server := NewPlayerServer(&store)
		deletePerson := "Pepper"

		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/store/%s", deletePerson), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.DeleteCalls) != 1 {
			t.Fatalf("got %d calls to DeleteUser want %d", len(store.NewUserCalls), 1)
		}

		if store.DeleteCalls[0] != deletePerson {
			t.Errorf("did not store correct winner got %q want %q", store.NewUserCalls[0], deletePerson)
		}


	})
}

func TestLogin(t *testing.T){
	t.Run("")
}

func TestLeague(t *testing.T){
	//store := StubPlayerStore{}
	//server := NewPlayerServer(&store)

	t.Run("it returns the league table as JSON", func(t *testing.T){
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, nil ,wantedLeague, nil}
		server := NewPlayerServer(&store)

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)

		assertContentType(t, response, jsonContentType)

	})
}


func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/store/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request{
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/store/%s", name),nil)
	return req

}

func newPutPlayerRequest(name string, reqBody []byte) *http.Request{
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/store/%s", name), bytes.NewBuffer(reqBody))
	return req
}

func assertResponseBody(t testing.TB, got string, want string){
	t.Helper()
	if want != got {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got int, want int){
	t.Helper()
	if got != want{
		t.Errorf("Did not get the expected status. got %d, want %d", got, want)
	}
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player){
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v", body, err)
	}

	return //league
}

func assertLeague(t testing.TB, got []Player, want []Player){
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string){
	t.Helper()
	contentTypeHeaderGot := response.Result().Header.Get("content-type")
	if contentTypeHeaderGot != want {
		t.Errorf("response did not have content-type of %s, got %v ", want, contentTypeHeaderGot)
	}
}

func assertScoreEquals(t testing.TB, got int, want int){
	t.Helper()
	if got != want{
		t.Errorf("different Score than expected. got %d want %d", got, want)
	}

}




