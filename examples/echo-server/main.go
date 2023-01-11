package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	cowsay "github.com/Code-Hex/Neo-cowsay/v2"
)

const limit = 1000

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server addr =>", ln.Addr())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				log.Println(err)
				return
			}

			go func() {
				defer func(conn net.Conn) {
					err := conn.Close()
					if err != nil {
						log.Println(err)
					}
				}(conn)
				go func() {
					<-ctx.Done()
					// cancel
					err := conn.SetDeadline(time.Now())
					if err != nil {
						log.Println(err)
						return
					}
				}()

				var buf strings.Builder
				rd := io.LimitReader(conn, limit)
				if _, err := io.Copy(&buf, rd); err != nil {
					log.Println("error:", err)
					return
				}

				phrase := strings.TrimSpace(buf.String())
				log.Println(phrase)
				say, _ := cowsay.Say(phrase)
				_, err := conn.Write([]byte(say))
				if err != nil {
					log.Println(err)
					return
				}
			}()
		}
	}()

	<-ctx.Done()
}
