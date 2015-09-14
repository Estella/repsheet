package main

import (
	"testing"
	"github.com/repsheet/repsheet/Godeps/_workspace/src/github.com/fzzy/radix/redis"
)


func cleanup(conn *redis.Client) {
	conn.Cmd("FLUSHDB")
	conn.Close()
}

func TestAddToListWhitelist(t *testing.T) {
	connection := connect("localhost", 6379, 10)

	addToList(connection, "whitelist", "1.1.1.1", "whitelist test")
	status := status(connection, "1.1.1.1")
	if status.Type != Whitelisted {
		t.Error("1.1.1.1 should have been whitelisted, got:", status.Type)
	}

	cleanup(connection)
}

func TestAddToListBlacklist(t *testing.T) {
	connection := connect("localhost", 6379, 10)

	addToList(connection, "blacklist", "1.1.1.1", "blacklist test")
	status := status(connection, "1.1.1.1")
	if status.Type != Blacklisted {
		t.Error("1.1.1.1 should have been blacklisted, got:", status.Type)
	}

	cleanup(connection)
}

func TestAddToListMark(t *testing.T) {
	connection := connect("localhost", 6379, 10)

	addToList(connection, "mark", "1.1.1.1", "mark test")
	status := status(connection, "1.1.1.1")
	if status.Type != Marked {
		t.Error("1.1.1.1 should have been marked, got:", status.Type)
	}

	cleanup(connection)
}

