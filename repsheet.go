package main

import (
        "flag"
        "fmt"
        "github.com/repsheet/repsheet/Godeps/_workspace/src/github.com/fzzy/radix/redis"
        "time"
        "os"
        "syscall"
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
	Type StatusType
	reason string
	actor string
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
		fmt.Printf("%s is %s. Reason:  %s\n", status.actor, statusStr, status.reason)
	} else {
		fmt.Printf("%s is %s\n", status.actor, statusStr)
	}
}

func printList(actorType string, list *redis.Reply) {
        if len(list.Elems) > 0 {
                fmt.Printf("%s Actors\n", actorType)
        }
        for i := 0; i < len(list.Elems); i++ {
                fmt.Printf("  %s\n", list.Elems[i])
        }
}

func addToList(connection *redis.Client, list string, actor string, reason string) {
        actorString := fmt.Sprintf("%s:repsheet:ip:%sed", actor, list)

        if reason != "" {
                connection.Cmd("SET", actorString, reason)
        } else {
                connection.Cmd("SET", actorString, "true")
        }
}

func connect(host string, port int, timeout int) *redis.Client {
	connectionString := fmt.Sprintf("%s:%d", host, port)
        conn, err := redis.DialTimeout("tcp", connectionString, time.Duration(timeout)*time.Second)

        if err != nil {
                fmt.Println("Cannot connect to Redis, exiting.")
                os.Exit(int(syscall.ECONNREFUSED))
        }

	return conn;
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
        blacklistPtr := flag.String("blacklist", "", "Blacklist an actor")
        whitelistPtr := flag.String("whitelist", "", "Whitelist an actor")
        markPtr := flag.String("mark", "", "Mark an actor")
        reasonPtr := flag.String("reason", "", "Reason for the action")

        flag.Parse()

	connection := connect("localhost", 6379, 10)

        if *listPtr == true {
                whitelisted := connection.Cmd("KEYS", "*:repsheet:*:whitelisted")
                printList("Whitelisted", whitelisted)

                blacklisted := connection.Cmd("KEYS", "*:repsheet:*:blacklisted")
                printList("Blacklisted", blacklisted)

                marked := connection.Cmd("KEYS", "*:repsheet:*:marked")
                printList("Marked", marked)
        }

	if *statusPtr != "" {
		status := status(connection, *statusPtr)
		printStatus(status)
	}

        if *blacklistPtr != "" {
                addToList(connection, "blacklist", *blacklistPtr, *reasonPtr)
        }

        if *whitelistPtr != "" {
                addToList(connection, "whitelist", *whitelistPtr, *reasonPtr)
        }

        if *markPtr != "" {
                addToList(connection, "mark", *markPtr, *reasonPtr)
        }
}
