package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type Client struct {
	ID    int64  `json:"-"`
	UUID  string `json:"uuid"`
	Login string `json:"login"`
	Email string `json:"mail"`
}

type Clients []Client

func NewTableClient() *Clients {
	return &Clients{}
}

func (s *Clients) All(db *sql.DB) ([]Client, error) {
	rows, err := db.Query("SELECT id, uuid, login, email FROM clients WHERE date_delete IS NULL ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	items := make([]Client, 0)

	for rows.Next() {
		var item Client

		err = rows.Scan(
			&item.ID,
			&item.UUID,
			&item.Login,
			&item.Email,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Clients) Save(item Client, db *sql.DB) (Client, error) {
	err := db.QueryRow(
		"INSERT INTO clients(login, email) VALUES ($1, $2) RETURNING id, uuid",
		item.Login,
		item.Email,
	).Scan(
		&item.ID,
		&item.UUID,
	)
	if err != nil {
		return Client{}, err
	}

	return item, nil
}

func (s *Clients) Get(uuid string, db *sql.DB) (Client, error) {
	item := Client{}

	err := db.QueryRow(
		"SELECT id, uuid, login, email FROM clients WHERE uuid = $1 AND date_delete IS NULL",
		uuid,
	).Scan(
		&item.ID,
		&item.UUID,
		&item.Login,
		&item.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return item, fmt.Errorf("client not found")
		}
		return item, err
	}

	return item, nil
}

func (s *Clients) Update(uuid string, item Client, db *sql.DB) (Client, error) {
	if isset := clients.Isset(uuid, db); isset == false {
		return item, fmt.Errorf("client not found")
	}

	err := db.QueryRow(
		"UPDATE clients SET login = $1, email = $2 WHERE uuid = $3 RETURNING id, uuid",
		item.Login,
		item.Email,
		uuid,
	).Scan(
		&item.ID,
		&item.UUID,
	)
	if err != nil {
		return Client{}, err
	}

	return item, nil
}

func (s *Clients) Destroy(uuid string, db *sql.DB) error {
	_, err := db.Query("UPDATE clients SET date_delete = NOW() WHERE uuid = $1 AND date_delete IS NULL", uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Clients) Isset(uuid string, db *sql.DB) bool {
	item := Client{}
	err := db.QueryRow(
		"SELECT id FROM clients WHERE uuid = $1 AND date_delete IS NULL",
		uuid,
	).Scan(
		&item.ID,
	)
	if err != nil {
		return false
	}
	if item.ID == 0 {
		return false
	}

	return true
}
