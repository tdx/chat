package main

import (
	"log"
	"net"
)

func listen() {

	const tag = "[LISTENER]:"

	chat := NewChat()

	go func() {
		l, err := net.Listen("tcp4", ":20000")
		if err != nil {
			log.Println(tag, "listen failed:", err)
		}
		defer l.Close()

		log.Println(tag, "listening on 20000 port")

		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println(tag, "accept failed:", err)
			}

			user := NewUser(conn)
			remover := chat.Register(user)
			user.OnClose(remover)

			go user.WriteMessages(chat)
			go user.ReadMessages(chat)
		}

	}()
}
