package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ncalibey/learn-go-with-tests/http-server/internal/poker"
)

func TestGETPlayers(t *testing.T) {
	store := poker.NewStubPlayerStore(map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	},
		nil,
		nil)
	server := poker.NewPlayerServer(store)

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
			request := poker.NewGetScoreRequest(c.name)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			poker.AssertStatus(t, response.Code, c.status)
			poker.AssertResponseBody(t, response.Body.String(), c.want)
		})
	}
}

func TestStoreWins(t *testing.T) {
	store := poker.NewStubPlayerStore(
		map[string]int{},
		nil,
		nil)
	server := poker.NewPlayerServer(store)

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
			request := poker.NewPostWinRequest(c.name)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			poker.AssertStatus(t, response.Code, http.StatusAccepted)

			if len(store.WinCalls) != 1 {
				t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
			}

			if store.WinCalls[0] != c.name {
				t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], c.name)
			}
		})
	}
}

func TestLeague(t *testing.T) {
	store := poker.NewStubPlayerStore(
		nil,
		nil,
		[]poker.Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		})
	server := poker.NewPlayerServer(store)

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
			request := poker.NewLeagueRequest()
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := poker.GetLeagueFromResponse(t, response.Body)
			poker.AssertStatus(t, response.Code, c.status)
			poker.AssertLeague(t, got, store.League)
			poker.AssertContentType(t, response, poker.JsonContentType)
		})
	}
}
