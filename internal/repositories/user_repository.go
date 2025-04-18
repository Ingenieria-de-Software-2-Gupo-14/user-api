package repositories

import (
	"context"
	"database/sql"
	"ing-soft-2-tp1/internal/models"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

// CreateDatabase creates and returns a database
func CreateDatabase(db *sql.DB) *Database {
	return &Database{DB: db}
}

func (db Database) GetUser(ctx context.Context, id int) (*models.User, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db Database) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := db.DB.QueryContext(ctx, "SELECT id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db Database) DeleteUser(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) AddUser(ctx context.Context, user *models.User) (int, error) {
	_, err := db.DB.ExecContext(ctx, "INSERT INTO users (username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		return 0, err
	}

	var id int
	err = db.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE email ILIKE $1", email)
	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db Database) ModifyUser(ctx context.Context, user *models.User) error {
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET username = $1, name= $2, surname=$3,  password=$4, email=$5, location=$6, admin=$7, blocked_user=$8, profile_photo=$9,description=$10 WHERE id = $11", &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description, &user.Id)
	return err
}

//id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description
