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
			conn.SetDeadline(time.Now().Add(5 * time.Second))

			go func() {
				defer conn.Close()
				go func() {
					<-ctx.Done()
					// cancel
					conn.SetDeadline(time.Now())
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
				conn.Write([]byte(say))
			}()
		}
	}()

	<-ctx.Done()
}
