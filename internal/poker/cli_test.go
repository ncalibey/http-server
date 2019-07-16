package poker_test

import (
	"strings"
	"testing"

	"github.com/ncalibey/learn-go-with-tests/http-server/internal/poker"
)

func TestCLI(t *testing.T) {
	cases := []struct {
		desc   string
		winner string
		want   string
	}{
		{
			desc:   "record chris win from user input",
			winner: "Chris wins\n",
			want:   "Chris",
		},
		{
			desc:   "record cleo win from user input",
			winner: "Cleo wins\n",
			want:   "Cleo",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			in := strings.NewReader(c.winner)
			playerStore := &poker.StubPlayerStore{}

			cli := poker.NewCLI(playerStore, in)
			cli.PlayPoker()

			poker.AssertPlayerWin(t, playerStore, c.want)
		})
	}
}
