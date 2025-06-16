package repository

import (
	//"fmt"
	//"time"

	"errors"
	"log"

	//"github.com/PushinMax/lesta-tf-idf-go/internal/service"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	db *sqlx.DB
}

func newAuthApi(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) Authentication(login, password string) (uuid.UUID, error) {
	var user User
	if err := r.db.Get(&user, "select * from users where login = $1", &login); err != nil {
		log.Println(err.Error())
		return uuid.UUID{}, errors.New("incorrect username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err.Error())
		return uuid.UUID{}, errors.New("incorrect username or password")
	}
	return user.ID, nil
}

func (r *AuthRepo) Register(login, password string) error {
	var user User
	if err := r.db.Get(&user, "select * from users where login = $1", &login); err == nil {
		return errors.New("choose another login/password")
	}
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.New("choose another login/password")
	}

	jti := uuid.New().String()


	if _, err := r.db.Exec(
		"INSERT INTO users (id, login, password_hash) values($1, $2, $3)", 
		&jti, &login, string(password_hash),
	); err != nil {
		log.Println(err)
		return errors.New("choose another login/password")
	}
	return nil
}

func (r *AuthRepo) ChangePassword(id, password string) error {
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.New("choose another password")
	}
	if _, err := r.db.Exec(
		"UPDATE users SET password_hash = $1 WHERE id = $2",
		&password_hash, &id,
	); err != nil {
		return errors.New("password change denied")
	}
	return nil
}

func (r *AuthRepo) SetRefreshToken(userID uuid.UUID, token string) error {
	token_hash, err := bcrypt.GenerateFromPassword([]byte(token[:72]), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.New("invalid token")
	}
	if _, err := r.db.Exec(
		"UPDATE users SET token_hash = $1 WHERE id = $2",
		&token_hash, &userID,
	); err != nil {
		log.Println(err.Error())
		return errors.New("failed to set refresh token")
	}
	return nil
}


func (r *AuthRepo) CheckAndChangeRefreshToken(userID uuid.UUID, token_old, token_new string) error {
	var user User
	if err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", userID); err != nil {
		log.Println(err.Error())
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Token), []byte(token_old[:72])); err != nil {
		log.Println(err.Error())
		return errors.New("invalid old token")
	}

	token_new_hash, err := bcrypt.GenerateFromPassword([]byte(token_new[:72]), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.New("invalid new token")
	}

	if _, err := r.db.Exec(
		"UPDATE users SET token_hash = $1, updated_at = NOW() WHERE id = $2",
		&token_new_hash, &userID,
	); err != nil {
		log.Println(err.Error())
		return errors.New("failed to update refresh token")
	}

	return nil
}

func (r *AuthRepo) Logout(id string) error {
	if _, err := r.db.Exec(
		"UPDATE users SET token_hash = NULL, updated_at = NULL WHERE id = $1",
		&id,
	); err != nil {
		log.Println(err.Error())
		return errors.New("logout failed")
	}
	return nil
}

func (r *AuthRepo) DeleteUser(id string) error {
	if _, err := r.db.Exec(
		"DELETE FROM users WHERE id = $1",
		&id,
	); err != nil {
		log.Println(err.Error())
		return errors.New("failed to delete user")
	}
	return nil
}

