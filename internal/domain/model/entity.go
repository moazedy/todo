package model

type Entity interface {
	WithIDSet(string) Entity
	GetID() string
}
