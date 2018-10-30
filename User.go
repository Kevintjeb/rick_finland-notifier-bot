package main

import (
	"fmt"
	"gopkg.in/tucnak/telebot.v2"
)

type User struct {
	Id   int
	Name string
}

func New(id int, name string) User {
	return User{id, name}
}

func (user *User) toRecipient() telebot.Recipient {
	return &telebot.User{ID: user.Id, FirstName: user.Name}
}

func (user *User) SendMessage(message string, b *telebot.Bot) {
	send, e := b.Send(user.toRecipient(), message)
	if e != nil {
		panic(e)
	}
	fmt.Println(send)
}

func (user *User) LogUser() {
	fmt.Println(user)
}
