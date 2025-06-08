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

