package main

import (
	"github.com/fzzy/radix/redis"
	"reflect"
	"testing"
)

func cleanup(conn *redis.Client) {
	conn.Cmd("FLUSHDB")
	conn.Close()
}

func TestWhitelist(t *testing.T) {
	connection := connect("localhost", 6379, 10)

	addToList(connection, "whitelist", "1.1.1.1", "whitelist test", 86400)
	status := status(connection, "1.1.1.1")
	if status.Type != Whitelisted {
		t.Error("1.1.1.1 should have been whitelisted, got:", status.Type)
	}

	cleanup(connection)
}

func TestBlacklist(t *testing.T) {
	connection := connect("localhost", 6379, 10)

	addToList(connection, "blacklist", "1.1.1.1", "blacklist test", 86400)
	status := status(connection, "1.1.1.1")
	if status.Type != Blacklisted {
		t.Error("1.1.1.1 should have been blacklisted, got:", status.Type)
	}

	cleanup(connection)
}

func TestMark(t *testing.T) {
	connection := connect("localhost", 6379, 10)

	addToList(connection, "mark", "1.1.1.1", "mark test", 86400)
	status := status(connection, "1.1.1.1")
	if status.Type != Marked {
		t.Error("1.1.1.1 should have been marked, got:", status.Type)
	}

	cleanup(connection)
}

func TestList(t *testing.T) {
	connection := connect("localhost", 6379, 10)

	addToList(connection, "blacklist", "1.1.1.1", "list test", 86400)
	addToList(connection, "whitelist", "2.2.2.2", "list test", 86400)
	addToList(connection, "mark", "3.3.3.3", "list test", 86400)

	var expected = map[string][]string{
		"blacklisted": []string{"1.1.1.1"},
		"whitelisted": []string{"2.2.2.2"},
		"marked":      []string{"3.3.3.3"},
	}

	l := list(connection)

	if !reflect.DeepEqual(l, expected) {
		t.Error("Did not get expected result, got:", l)
	}
}
