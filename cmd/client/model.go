package main

import (
	"errors"
	"math/rand"
)

type TableActions interface {
	GetId() int
	GetUID() string
}

type RowSchemeClient struct {
	ID    int    `json:"-"`
	UID   string `json:"uid"`
	Login string `json:"login"`
	Email string `json:"mail"`
}

type TableSchemeClient struct {
	nextId int
	rows   []RowSchemeClient
}

var _ TableActions = (*TableSchemeClient)(nil)

func NewTableClient() *TableSchemeClient {
	return &TableSchemeClient{
		nextId: 1,
		rows:   make([]RowSchemeClient, 0),
	}
}

func (s *TableSchemeClient) GetId() int {
	id := s.nextId
	s.nextId++
	return id
}

func (s *TableSchemeClient) GetUID() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	const length = 3

	var bytes = make([]byte, length)

	for i := 0; i < length; i++ {
		idx := rand.Intn(len(letterBytes))
		bytes[i] = letterBytes[idx]
	}

	return string(bytes)
}

func (s *TableSchemeClient) Save(data RowSchemeClient) (RowSchemeClient, error) {
	data.ID = s.GetId()
	data.UID = s.GetUID()

	s.rows = append(s.rows, data)

	s.nextId++

	return data, nil
}

func (s *TableSchemeClient) All() []RowSchemeClient {
	return s.rows
}

func (s *TableSchemeClient) Get(uid string) (RowSchemeClient, error) {
	for _, row := range s.rows {
		if row.UID == uid {
			return row, nil
		}
	}

	return RowSchemeClient{}, errors.New("not found")
}

func (s *TableSchemeClient) Update(uid string, data RowSchemeClient) (RowSchemeClient, error) {
	for i, row := range s.rows {
		if row.UID == uid {
			s.rows[i].Login = data.Login
			s.rows[i].Email = data.Email
			return s.rows[i], nil
		}
	}
	return data, errors.New("not found")
}

func (s *TableSchemeClient) Destroy(uid string) error {
	for i, row := range s.rows {
		if row.UID == uid {
			s.rows = append(s.rows[:i], s.rows[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
