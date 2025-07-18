package hash

import (
	"slices"
	"testing"

	"github.com/nalgeon/be"
	"github.com/nalgeon/redka/redsrv/internal/redis"
)

func TestHKeysParse(t *testing.T) {
	tests := []struct {
		cmd string
		key string
		err error
	}{
		{
			cmd: "hkeys",
			key: "",
			err: redis.ErrInvalidArgNum,
		},
		{
			cmd: "hkeys person",
			key: "person",
			err: nil,
		},
		{
			cmd: "hkeys person name",
			key: "",
			err: redis.ErrInvalidArgNum,
		},
	}

	for _, test := range tests {
		t.Run(test.cmd, func(t *testing.T) {
			cmd, err := redis.Parse(ParseHKeys, test.cmd)
			be.Equal(t, err, test.err)
			if err == nil {
				be.Equal(t, cmd.key, test.key)
			} else {
				be.Equal(t, cmd, HKeys{})
			}
		})
	}
}

func TestHKeysExec(t *testing.T) {
	t.Run("key found", func(t *testing.T) {
		red := getRedka(t)

		_, _ = red.Hash().Set("person", "name", "alice")
		_, _ = red.Hash().Set("person", "age", 25)

		cmd := redis.MustParse(ParseHKeys, "hkeys person")
		conn := redis.NewFakeConn()
		res, err := cmd.Run(conn, red)

		be.Err(t, err, nil)
		got := res.([]string)
		slices.Sort(got)
		be.Equal(t, got, []string{"age", "name"})
		be.True(t, conn.Out() == "2,age,name" || conn.Out() == "2,name,age")
	})
	t.Run("key not found", func(t *testing.T) {
		red := getRedka(t)

		cmd := redis.MustParse(ParseHKeys, "hkeys person")
		conn := redis.NewFakeConn()
		res, err := cmd.Run(conn, red)

		be.Err(t, err, nil)
		be.Equal(t, res.([]string), []string{})
		be.Equal(t, conn.Out(), "0")
	})
}
