// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Refuser refuses any outside attempt to connect to it.

func main() {
	var port serv.Port = 80

	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli.Var(&port, "port", "Port to listen to")
	_ = cli.Parse(os.Args[1:])

	tcp, err := net.Listen("tcp", port.Addr())
	errors.FatalOnErr(err)
	defer tcp.Close()
	log.Println("listening on", port.Addr())

	ctx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer stopFn()
		for {
			conn, err := tcp.Accept()
			if err != nil {
				if !strings.Contains(err.Error(), "closed network connection") {
					log.Println("err:", err)
				}
				continue
			}

			log.Println("refuse:", conn.RemoteAddr())
			_ = conn.Close()
		}
	}()
	<-ctx.Done()
}
