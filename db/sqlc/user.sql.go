// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    hash,
    full_name,
    email
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, username, hash, full_name, email, created_at, password_changed_at
`

type CreateUserParams struct {
	Username string `json:"username"`
	Hash     string `json:"hash"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Hash,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Hash,
		&i.FullName,
		&i.Email,
		&i.CreatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, hash, full_name, email, created_at, password_changed_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Hash,
		&i.FullName,
		&i.Email,
		&i.CreatedAt,
		&i.PasswordChangedAt,
	)
	return i, err
}
