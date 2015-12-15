package main

import (
	"flag"
	"fmt"
	"github.com/fzzy/radix/redis"
	"os"
	"strings"
	"syscall"
	"time"
)

type StatusType uint8

const (
	UnitType StatusType = iota
	Blacklisted
	Whitelisted
	Marked
	OK
)

type Status struct {
	Type   StatusType
	reason string
	actor  string
}

func printStatus(status *Status) {
	var statusStr string
	switch status.Type {
	case Blacklisted:
		statusStr = "blacklisted"
	case Whitelisted:
		statusStr = "whitelisted"
	case Marked:
		statusStr = "marked"
	default:
		statusStr = "OK"
	}

	if status.reason != "" {
		fmt.Printf("%s is %s. Reason: %s\n", status.actor, statusStr, status.reason)
	} else {
		fmt.Printf("%s is %s\n", status.actor, statusStr)
	}
}

func printList(lists map[string][]string) {
	for k, v := range lists {
		fmt.Printf("%s actors\n", k)
		for i := 0; i < len(v); i++ {
			fmt.Printf("  %s\n", v[i])
		}
	}
}

func list(connection *redis.Client) map[string][]string {
	var blacklisted []string
	var whitelisted []string
	var marked []string

	repsheet := connection.Cmd("KEYS", "*:repsheet:*:*")
	for i := 0; i < len(repsheet.Elems); i++ {
		value, _ := repsheet.Elems[i].Str()
		parts := strings.Split(value, ":")
		repsheetType := parts[len(parts)-1]

		switch repsheetType {
		case "blacklisted":
			blacklisted = append(blacklisted, parts[0])
		case "whitelisted":
			whitelisted = append(whitelisted, parts[0])
		case "marked":
			marked = append(marked, parts[0])
		}
	}

	var list = make(map[string][]string)
	list["blacklisted"] = blacklisted
	list["whitelisted"] = whitelisted
	list["marked"] = marked

	return list
}

func addToList(connection *redis.Client, list string, actor string, reason string, ttl int) {
	actorString := fmt.Sprintf("%s:repsheet:ip:%sed", actor, list)

	if reason != "" {
		connection.Cmd("SETEX", actorString, ttl, reason)
	} else {
		connection.Cmd("SETEX", actorString, ttl, "true")
	}
}

func removeFromLists(connection *redis.Client, actor string) {
	removeFromList(connection, "blacklist", actor)
	removeFromList(connection, "whitelist", actor)
	removeFromList(connection, "mark", actor)
}

func removeFromList(connection *redis.Client, list string, actor string) {
	actorString := fmt.Sprintf("%s:repsheet:ip:%sed", actor, list)
	connection.Cmd("DEL", actorString)
}

func connect(host string, port int, timeout int) *redis.Client {
	connectionString := fmt.Sprintf("%s:%d", host, port)
	conn, err := redis.DialTimeout("tcp", connectionString, time.Duration(timeout)*time.Second)

	if err != nil {
		fmt.Println("Cannot connect to Redis, exiting.")
		os.Exit(int(syscall.ECONNREFUSED))
	}

	return conn
}

func status(connection *redis.Client, actor string) *Status {
	blacklistedString := fmt.Sprintf("%s:repsheet:ip:blacklisted", actor)
	whitelistedString := fmt.Sprintf("%s:repsheet:ip:whitelisted", actor)
	markedString := fmt.Sprintf("%s:repsheet:ip:marked", actor)
	connection.Cmd("MULTI")
	connection.Cmd("GET", whitelistedString)
	connection.Cmd("GET", blacklistedString)
	connection.Cmd("GET", markedString)
	reply := connection.Cmd("EXEC")

	if reply.Elems[0].Type != redis.NilReply {
		str, _ := reply.Elems[0].Str()
		return &Status{Type: Whitelisted, reason: str, actor: actor}
	} else if reply.Elems[1].Type != redis.NilReply {
		str, _ := reply.Elems[1].Str()
		return &Status{Type: Blacklisted, reason: str, actor: actor}
	} else if reply.Elems[2].Type != redis.NilReply {
		str, _ := reply.Elems[1].Str()
		return &Status{Type: Marked, reason: str, actor: actor}
	} else {
		return &Status{Type: OK, reason: "", actor: actor}
	}
}

func main() {
	listPtr := flag.Bool("list", false, "Show the contents of Repsheet")
	statusPtr := flag.String("status", "", "Get the status of an actor")
	removePtr := flag.String("remove", "", "Remove an actor from all lists")
	blacklistPtr := flag.String("blacklist", "", "Blacklist an actor")
	whitelistPtr := flag.String("whitelist", "", "Whitelist an actor")
	markPtr := flag.String("mark", "", "Mark an actor")
	reasonPtr := flag.String("reason", "", "Reason for the action")
	ttlPtr := flag.Int("ttl", 86400, "Set expiry in Redis")
	hostPtr := flag.String("host", "localhost", "Redis host")
	portPtr := flag.Int("port", 6379, "Redis port")
	timeoutPtr := flag.Int("timeout", 10, "Redis connection timeout")

	flag.Parse()

	connection := connect(*hostPtr, *portPtr, *timeoutPtr)

	if *listPtr == true {
		l := list(connection)
		printList(l)
	}

	if *statusPtr != "" {
		status := status(connection, *statusPtr)
		printStatus(status)
	}

	if *blacklistPtr != "" {
		addToList(connection, "blacklist", *blacklistPtr, *reasonPtr, *ttlPtr)
	}

	if *whitelistPtr != "" {
		addToList(connection, "whitelist", *whitelistPtr, *reasonPtr, *ttlPtr)
	}

	if *markPtr != "" {
		addToList(connection, "mark", *markPtr, *reasonPtr, *ttlPtr)
	}

	if *removePtr != "" {
		removeFromLists(connection, *removePtr)
	}
}
