package main

import (
        "flag"
        "fmt"
        "github.com/repsheet/repsheet/Godeps/_workspace/src/github.com/fzzy/radix/redis"
        "time"
        "os"
        "syscall"
)

func printList(actorType string, list *redis.Reply) {
        if len(list.Elems) > 0 {
                fmt.Printf("%s Actors\n", actorType)
        }
        for i := 0; i < len(list.Elems); i++ {
                fmt.Printf("  %s\n", list.Elems[i])
        }
}

func main() {
        listPtr := flag.Bool("list", false, "Show the contents of Repsheet")
        blacklistPtr := flag.String("blacklist", "", "Blacklist an actor")
        whitelistPtr := flag.String("whitelist", "", "Whitelist an actor")
        markPtr := flag.String("mark", "", "Mark an actor")
        reasonPtr := flag.String("reason", "", "Reason for the action")

        flag.Parse()

        conn, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
        if err != nil {
                fmt.Println("Cannot connect to Redis, exiting.")
                os.Exit(int(syscall.ECONNREFUSED))
        }

        if *listPtr == true {
                whitelisted := conn.Cmd("KEYS", "*:repsheet:*:whitelisted")
                printList("Whitelisted", whitelisted)

                blacklisted := conn.Cmd("KEYS", "*:repsheet:*:blacklisted")
                printList("Blacklisted", blacklisted)

                marked := conn.Cmd("KEYS", "*:repsheet:*:marked")
                printList("Marked", marked)
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
