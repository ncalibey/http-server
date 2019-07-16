package poker

import "testing"

func TestFileSystemStore(t *testing.T) {
	cases := []struct {
		desc string
		want []Player
	}{
		{
			desc: "league sorted",
			want: []Player{
				{"Chris", 33},
				{"Cleo", 10},
			},
		},
	}
	database, cleanDatabase := CreateTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			store, err := NewFileSystemPlayerStore(database)

			AssertNoError(t, err)

			got := store.GetLeague()
			AssertLeague(t, got, c.want)

			// read again
			got = store.GetLeague()
			AssertLeague(t, got, c.want)
		})
	}
}

func TestFileSystemGetPlayerStore(t *testing.T) {
	cases := []struct {
		desc  string
		want1 int
		want2 int
	}{
		{
			desc:  "get player score",
			want1: 33,
			want2: 34,
		},
	}
	database, cleanDatabase := CreateTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}]`)
	defer cleanDatabase()

	store, err := NewFileSystemPlayerStore(database)
	AssertNoError(t, err)

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			got := store.GetPlayerScore("Chris")
			AssertPlayerScore(t, got, c.want1)
		})

		t.Run(c.desc, func(t *testing.T) {
			store.RecordWin("Chris")
			got := store.GetPlayerScore("Chris")
			AssertPlayerScore(t, got, c.want2)
		})
	}
}
