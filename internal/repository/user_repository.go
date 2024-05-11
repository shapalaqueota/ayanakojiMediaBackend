package repository

import (
	"backend/internal/models"
	"backend/internal/utils"
	"context"
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

func GetUserById(conn *pgxpool.Pool, uuid string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password FROM "user" WHERE id = $1`
	err := conn.QueryRow(context.Background(), query, uuid).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
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
