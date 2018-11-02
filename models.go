package api

import (
	"time"
)

type userModel struct {
	tableName struct{} `sql:"users"`

	ID       int64
	Name     string
	Email    string   `sql:",unique"`
	Roles    []string `sql:",array"`
	Salt     []byte
	Password []byte
}

type sessionData struct {
	Auth              bool
	Name              string
	Email             string
	Roles             []string
	MarkedForDeletion bool
}

type sessionModel struct {
	tableName struct{} `sql:"sessions"`

	ID       int64
	Key      []byte `sql:",unique"`
	Data     sessionData
	LastUsed time.Time
}
