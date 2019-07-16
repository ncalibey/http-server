package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const jsonContentType = "application/json"

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(store)

	cases := []struct {
		desc   string
		name   string
		want   string
		status int
	}{
		{
			desc:   "returns Pepper's score",
			name:   "Pepper",
			want:   "20",
			status: http.StatusOK,
		},
		{
			desc:   "return's Floyd's score",
			name:   "Floyd",
			want:   "10",
			status: http.StatusOK,
		},
		{
			desc:   "returns 404 on missing players",
			name:   "Apollo",
			want:   "0",
			status: http.StatusNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			request := newGetScoreRequest(c.name)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, c.status)
			assertResponseBody(t, response.Body.String(), c.want)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := &StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(store)

	cases := []struct {
		desc   string
		name   string
		status int
	}{
		{
			desc:   "it returns accepted on POST",
			name:   "Pepper",
			status: http.StatusAccepted,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			request := newPostWinRequest(c.name)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, http.StatusAccepted)

			if len(store.winCalls) != 1 {
				t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
			}

			if store.winCalls[0] != c.name {
				t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], c.name)
			}
		})
	}
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	cases := []struct {
		desc   string
		status int
	}{
		{
			desc:   "get score",
			status: http.StatusOK,
		},
		{
			desc:   "get league",
			status: http.StatusOK,
		},
	}
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	for i, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			response := httptest.NewRecorder()
			server.ServeHTTP(response, newGetScoreRequest(player))
			assertStatus(t, response.Code, http.StatusOK)

			if i > 0 {
				got := getLeagueFromResponse(t, response.Body)
				want := []Player{
					{"Pepper", 3},
				}
				assertLeague(t, got, want)
			} else {
				assertResponseBody(t, response.Body.String(), "3")
			}
		})
	}

}

func TestLeague(t *testing.T) {
	store := &StubPlayerStore{
		scores:   nil,
		winCalls: nil,
		league: []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		},
	}
	server := NewPlayerServer(store)

	cases := []struct {
		desc   string
		status int
	}{
		{
			desc:   "it returns 200 on /league",
			status: http.StatusOK,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			request := newLeagueRequest()
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := getLeagueFromResponse(t, response.Body)
			assertStatus(t, response.Code, c.status)
			assertLeague(t, got, store.league)
			assertContentType(t, response, jsonContentType)
		})
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//// Helper Functions ////////////////////////////////////////////////////////////////////

func newPostWinRequest(player string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func assertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
