package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Clients struct{}

type ClientItem struct {
	ID    int64  `json:"-"`
	UUID  string `json:"uuid"`
	Login string `json:"login"`
	Email string `json:"mail"`
}

type ClientItems []ClientItem

func NewClients() *Clients {
	return &Clients{}
}

var clients = NewClients()

var (
	ErrNotFound = errors.New("client not found")
)

func (s *Clients) All(ctx context.Context, db *sql.DB) (ClientItems, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctxTimeout, "SELECT id, uuid, login, email FROM clients WHERE date_delete IS NULL ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	items := make([]ClientItem, 0)

	for rows.Next() {
		var item ClientItem

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

func (s *Clients) Save(ctx context.Context, db *sql.DB, item ClientItem) (ClientItem, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := db.QueryRowContext(
		ctxTimeout,
		"INSERT INTO clients(login, email) VALUES ($1, $2) RETURNING id, uuid",
		item.Login,
		item.Email,
	).Scan(
		&item.ID,
		&item.UUID,
	)
	if err != nil {
		return ClientItem{}, err
	}

	return item, nil
}

func (s *Clients) Get(ctx context.Context, db *sql.DB, uuid string) (ClientItem, error) {
	ctxTimeOut, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	item := ClientItem{}

	err := db.QueryRowContext(
		ctxTimeOut,
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
			return item, ErrNotFound
		}
		return item, err
	}

	return item, nil
}

func (s *Clients) Update(ctx context.Context, db *sql.DB, uuid string, item ClientItem) (ClientItem, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := db.QueryRowContext(
		ctxTimeout,
		"UPDATE clients SET login = $1, email = $2 WHERE uuid = $3 AND date_delete IS NULL RETURNING id, uuid",
		item.Login,
		item.Email,
		uuid,
	).Scan(
		&item.ID,
		&item.UUID,
	)
	if err != nil {
		return ClientItem{}, err
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

func (s *Clients) Isset(ctx context.Context, db *sql.DB, uuid string) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	item := ClientItem{}
	err := db.QueryRowContext(
		ctxTimeout,
		"SELECT id FROM clients WHERE uuid = $1 AND date_delete IS NULL",
		uuid,
	).Scan(
		&item.ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
