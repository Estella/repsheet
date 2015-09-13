package main

import (
	"flag"
	"fmt"
	"github.com/repsheet/repsheet/Godeps/_workspace/src/github.com/fzzy/radix/redis"
	"time"
)

func main() {
	listPtr := flag.Bool("list", false, "Show the contents of Repsheet")
	blacklistPtr := flag.String("blacklist", "", "Blacklist an actor")
	whitelistPtr := flag.String("whitelist", "", "Whitelist an actor")
	markPtr := flag.String("mark", "", "Mark an actor")
	reasonPtr := flag.String("reason", "", "Reason for the action")

	flag.Parse()

	conn, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
	if err != nil {
		panic("Error connecting to Redis")
	}

	if *listPtr == true {
		whitelisted := conn.Cmd("KEYS", "*:repsheet:*:whitelisted")
		if len(whitelisted.Elems) > 0 {
			fmt.Println("Whitelisted Actors")
		}
		for i := 0; i < len(whitelisted.Elems); i++ {
			fmt.Printf("  %s\n", whitelisted.Elems[i])
		}

		blacklisted := conn.Cmd("KEYS", "*:repsheet:*:blacklisted")
		if len(blacklisted.Elems) > 0 {
			fmt.Println("Blacklisted Actors")
		}
		for i := 0; i < len(blacklisted.Elems); i++ {
			fmt.Printf("  %s\n", blacklisted.Elems[i])
		}

		marked := conn.Cmd("KEYS", "*:repsheet:*:marked")
		if len(marked.Elems) > 0 {
			fmt.Println("Marked Actors")
		}
		for i := 0; i < len(marked.Elems); i++ {
			fmt.Printf("  %s\n", marked.Elems[i])
		}
	}

	if *blacklistPtr != "" {
		fmt.Println("Blacklisting", *blacklistPtr)
		command := fmt.Sprintf("%s:repsheet:ip:blacklisted", *blacklistPtr)

		if *reasonPtr != "" {
			conn.Cmd("SET", command, *reasonPtr)
		} else {
			conn.Cmd("SET", command, "true")
		}

	}

	if *whitelistPtr != "" {
		fmt.Println("Whitelisting", *whitelistPtr)
		command := fmt.Sprintf("%s:repsheet:ip:whitelisted", *whitelistPtr)

		if *reasonPtr != "" {
			conn.Cmd("SET", command, *reasonPtr)
		} else {
			conn.Cmd("SET", command, "true")
		}
	}

	if *markPtr != "" {
		fmt.Println("Marking", *markPtr)
		command := fmt.Sprintf("%s:repsheet:ip:marked", *markPtr)

		if *reasonPtr != "" {
			conn.Cmd("SET", command, *reasonPtr)
		} else {
			conn.Cmd("SET", command, "true")
		}
	}
}
