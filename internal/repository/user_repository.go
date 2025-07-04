package repository

import (
	"backend/internal/models"
	"backend/internal/utils"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func CreateUser(conn *pgxpool.Conn, user models.User) (string, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	var userID string
	query := `INSERT INTO "user" (username, email, password, phone_number) VALUES ($1, $2, $3, $4) RETURNING id`
	err = conn.QueryRow(context.Background(), query, user.Username, user.Email, hashedPassword, user.PhoneNumber).Scan(&userID)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return "", err
	}
	return userID, nil
}

func UpdateUser(conn *pgxpool.Conn, user models.User) error {

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	query := `UPDATE "user" SET username=$1, password=$2 WHERE id=$3`
	cmdTag, err := conn.Exec(context.Background(), query, user.Username, hashedPassword, user.ID)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
	}
	if cmdTag.RowsAffected() != 1 {
		return errors.New("failed to update user")
	}
	return nil
}

func GetUserById(conn *pgxpool.Pool, uuid string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, phone_number, email_verified FROM "user" WHERE id = $1`
	err := conn.QueryRow(context.Background(), query, uuid).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.PhoneNumber, &user.EmailVerified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(conn *pgxpool.Pool, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password FROM "user" WHERE email = $1`
	err := conn.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CheckUserExists(db *pgxpool.Conn, email string, phoneNumber string) (bool, error) {
	var exists bool
	err := db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM "user" 
                       WHERE email = $1 OR phone_number = $2)`, email, phoneNumber).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func UpdateEmailVerified(conn *pgxpool.Conn, userID string, verified bool) error {
	sql := `UPDATE "user" SET email_verified = $1 WHERE id = $2`
	cmdTag, err := conn.Exec(context.Background(), sql, verified, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != 1 {
		return errors.New("no user found or updated")
	}
	return nil
}
