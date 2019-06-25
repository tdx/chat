package main

import (
	"container/list"
	"sync"
)

// Chat holds all active users
type Chat struct {
	mu    sync.Mutex
	users list.List
	hist  [][]byte
}

// NewChat creates chat object
func NewChat() *Chat {
	return &Chat{}
}

// Register adds user to chat and returns 'remover' func
func (c *Chat) Register(user *User) (remove func()) {

	c.mu.Lock()
	defer c.mu.Unlock()

	el := c.users.PushBack(user)

	return func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.users.Remove(el)
	}
}

// Broadcast sends message to all users in the caht
func (c *Chat) Broadcast(from *User, msg []byte) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.hist) > 9 {
		c.hist = c.hist[1:]
	}
	c.hist = append(c.hist, msg)

	for el := c.users.Front(); el != nil; el = el.Next() {
		user := el.Value.(*User)
		// don't send to self
		if from != nil && from.Equal(user) {
			continue
		}
		if err := user.WriteMessage(msg); err != nil {
			// remove user from chat ?
		}
	}
}

// History ...
func (c *Chat) History() [][]byte {
	return c.hist
}
