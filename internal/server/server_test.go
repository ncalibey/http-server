package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

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
//// Assertions //////////////////////////////////////////////////////////////////////////

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

func assertLeague(t *testing.T, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertPlayerScore(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didnt expect an error but got one: %v", err)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//// Helper Functions ////////////////////////////////////////////////////////////////////

func newPostWinRequest(player string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	return req
}

func getLeagueFromResponse(t *testing.T, body io.Reader) []Player {
	t.Helper()
	league, err := NewLeague(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return league
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
