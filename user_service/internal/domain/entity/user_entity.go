package entity

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleGuest:
		return true
	}
	return false
}

type User struct {
	ID             int64           `db:"id"`
	Login          string          `db:"login"`
	Name           string          `db:"name"`
	UID            uuid.UUID       `db:"uid"`
	Email          string          `db:"email"`
	Password       string          `db:"password"`
	IsActive       bool            `db:"is_active"`
	LastLogin      sql.NullTime    `db:"last_login"`
	Role           Role            `db:"role"`
	ProfilePicture sql.NullString  `db:"profile_picture"`
	Phone          sql.NullString  `db:"phone"`
	Metadata       json.RawMessage `db:"metadata"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
	DeletedAt      sql.NullTime    `db:"deleted_at"`
}
