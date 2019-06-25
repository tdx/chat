package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"sync"
)

// User struct
type User struct {
	mu      sync.Mutex
	onClose func()

	once sync.Once
	conn net.Conn
	out  chan []byte
}

// NewUser creates user
func NewUser(conn net.Conn) *User {
	return &User{
		conn: conn,
		out:  make(chan []byte, 16),
	}
}

var (
	errUserOutError = errors.New("user output is full")
)

func (u *User) OnClose(f func()) {
	u.onClose = f
}

// Close ...
func (u *User) Close() {
	u.once.Do(func() {
		close(u.out)
		u.conn.Close()

		u.mu.Lock()
		fc := u.onClose
		u.onClose = nil
		u.mu.Unlock()

		fc()
	})
}

// ReadMessages reads messages from user connection and broadcast it to chat
func (u *User) ReadMessages(chat *Chat) {

	buf := bufio.NewReader(u.conn)

	for {
		msg, err := buf.ReadBytes('\n')
		if err != nil {
			log.Println("client read failed:", err)
			u.Close()
			return
		}
		chat.Broadcast(u, msg)
	}
}

// Equal compares user's remote addresses
func (u *User) Equal(u2 *User) bool {
	return u.conn.RemoteAddr() == u2.conn.RemoteAddr()
}

// WriteMessages from broadcast to user connection
func (u *User) WriteMessages(chat *Chat) {

	// Write chat history to new user
	for _, msg := range chat.History() {
		if _, err := u.conn.Write(msg); err != nil {
			log.Println("client write failed:", err)
			u.Close()
			return
		}
	}

	var msg []byte
	for ok := true; ok; msg, ok = <-u.out {
		if _, err := u.conn.Write(msg); err != nil {
			log.Println("client write failed:", err)
			u.Close()
			return
		}
	}
}

// WriteMessage writes message to user channel
func (u *User) WriteMessage(msg []byte) error {
	select {
	case u.out <- msg:
		return nil
	default:
		return errUserOutError
	}
}

// func (u *User) readMessage(r io.Reader) ([]byte, error) {

// 	return r.ReadBytes("\n")
// }
