package gobot_test

import (
	"testing"
)

func TestPing(t *testing.T) {
	AssertGobotResponseContains(t, "!ping", "pong")
}
