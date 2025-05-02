package repositories

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"ing-soft-2-tp1/internal/errors"
	"ing-soft-2-tp1/internal/models"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

// CreateUserRepo creates and returns a database
func CreateUserRepo(db *sql.DB) *Database {
	return &Database{DB: db}
}

func (db Database) GetUser(ctx context.Context, id int) (*models.User, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description FROM users WHERE id = $1", id)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
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
	_, errInsert := db.DB.ExecContext(ctx, "INSERT INTO users (username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if errInsert != nil {
		pqErr, ok := errInsert.(*pq.Error)
		if ok {
			if pqErr.Code == "23505" {
				return 0, errors.ErrEmailInUsed
			}
		}
		return 0, errInsert
	}

	var id int
	errSearch := db.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", user.Email).Scan(&id)
	if errSearch != nil {
		return 0, errSearch
	}

	return id, nil
}

func (db Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description FROM users WHERE email ILIKE $1", email)
	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db Database) ModifyUser(ctx context.Context, user *models.User) error {
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET username = $1, name= $2, surname=$3,  location=$4, profile_photo=$5,description=$6 WHERE id = $7", &user.Username, &user.Name, &user.Surname, &user.Location, &user.ProfilePhoto, &user.Description, &user.Id)
	return err
}

func (db Database) ModifyLocation(ctx context.Context, id int, newLocation string) error {
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET location = $1 where id = $2", newLocation, id)
	return err
}

func (db Database) BlockUser(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET blocked_user = true where id = $1", id)
	return err
}

func (db Database) GetUserPrivacy(ctx context.Context, id int) (*models.UserPrivacy, error) {
	row := db.DB.QueryRowContext(ctx, "SELECT name_privacy, surname_privacy, email_privacy, location_privacy, description_privacy FROM users WHERE id = $1", id)

	var userPrivacy models.UserPrivacy
	err := row.Scan(&userPrivacy.Name, &userPrivacy.Surname, &userPrivacy.Email, &userPrivacy.Location, &userPrivacy.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &userPrivacy, nil
}

func (db Database) ModifyPrivacy(ctx context.Context, id int, privacy models.UserPrivacy) error {
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET name_privacy = $2, surname_privacy = $3, email_privacy = $4, location_privacy = $5, description_privacy = $6 where id = $1", id, privacy.Name, privacy.Surname, privacy.Email, privacy.Location, privacy.Description)
	return err
}

//id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description
