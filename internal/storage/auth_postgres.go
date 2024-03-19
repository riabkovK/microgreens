package storage

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/riabkovK/microgreens/internal/domain"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (receiver *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, name, username, password_hash) VALUES ($1, $2, $3, $4) RETURNING id", usersTable)
	row := receiver.db.QueryRow(query, user.Email, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (receiver *AuthPostgres) GetByCredentials(email, passwordHash string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := receiver.db.Get(&user, query, email, passwordHash)

	return user, err
}

func (receiver *AuthPostgres) GetByRefreshToken(refreshToken string) (domain.User, error) {
	var user domain.User

	// Check if the refresh token has expired
	var refreshSession domain.RefreshSessionCheckResponse
	getRefreshSessionExpiresQuery := fmt.Sprintf(`SELECT rst.id, rst.user_id, rst.refresh_token, rst.expires_at
														 FROM %s AS rst
														 WHERE rst.refresh_token=$1`, refreshSessionsTable)

	err := receiver.db.Get(&refreshSession, getRefreshSessionExpiresQuery, refreshToken)
	if err != nil {
		return user, err
	}
	if refreshSession.ExpiresAt.Sub(time.Now()) < 0 {
		return user, domain.ErrRefreshTokenHasExpired
	}

	getUserQuery := fmt.Sprintf(`SELECT ut.id FROM %s AS ut
								 INNER JOIN %s AS rst ON ut.id = rst.user_id
							     WHERE rst.refresh_token=$1`, usersTable, refreshSessionsTable)
	err = receiver.db.Get(&user, getUserQuery, refreshToken)
	if err != nil {
		return user, err
	}
	return user, err
}

func (receiver *AuthPostgres) SetSession(userId int, session domain.Session) error {
	tx, err := receiver.db.Begin()
	if err != nil {
		return err
	}

	// Update table if refresh token is belonged to user
	updateRefreshTokenQuery := fmt.Sprintf(`UPDATE %s AS rst SET refresh_token=$1, expires_at=$2, created_at=$3 FROM %s AS ut
	                                              WHERE rst.user_id=ut.id AND rst.user_id=$4`,
		refreshSessionsTable, usersTable)

	result, err := tx.Exec(updateRefreshTokenQuery, session.RefreshToken, session.ExpiresAt, time.Now(), userId)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			return err1
		}
		return err
	}
	rows, _ := result.RowsAffected()

	// Refresh token does not belong to user
	if rows == 0 {
		createRefreshTokenQuery := fmt.Sprintf(`INSERT INTO %s (user_id, refresh_token, expires_at) VALUES ($1, $2, $3)
	                                                  RETURNING id`, refreshSessionsTable)
		_, err := tx.Exec(createRefreshTokenQuery, userId, session.RefreshToken, session.ExpiresAt)
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				return err1
			}
			return err
		}
	}

	// Update user last visit
	updateUserLastVisitQuery := fmt.Sprintf(`UPDATE %s AS ut SET last_visit_at=$1 FROM %s AS rst
													WHERE ut.id=rst.user_id and ut.id=$2`,
		usersTable, refreshSessionsTable)
	_, err = tx.Exec(updateUserLastVisitQuery, time.Now(), userId)
	if err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			return err1
		}
		return err
	}

	return tx.Commit()
}
