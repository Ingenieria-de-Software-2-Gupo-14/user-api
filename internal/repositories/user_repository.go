package repositories

import (
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

func (db Database) GetUser(id int) (*models.User, error) {
	row := db.DB.QueryRow("SELECT * FROM users WHERE id = $1", id)

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

func (db Database) GetAllUsers() ([]models.User, error) {
	rows, err := db.DB.Query("SELECT id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description FROM users")
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

func (db Database) DeleteUser(id int) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) AddUser(user *models.User) (int, error) {
	r, err := db.DB.Exec("INSERT INTO users (username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description)
	if err != nil {
		return 0, err
	}

	if id, err := r.LastInsertId(); err != nil {
		return 0, err
	} else {
		return int(id), nil
	}
}

func (db Database) GetUserByEmail(email string) (*models.User, error) {
	row := db.DB.QueryRow("SELECT * FROM users WHERE email ILIKE $1", email)
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

func (db Database) ModifyUser(user *models.User) error {
	_, err := db.DB.Exec("UPDATE users SET username = $1, name= $2, surname=$3,  password=$4, email=$5, location=$6, admin=$7, blocked_user=$8, profile_photo=$9,description=$10 WHERE id = $11", &user.Username, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.BlockedUser, &user.ProfilePhoto, &user.Description, &user.Id)
	return err
}

//id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description
