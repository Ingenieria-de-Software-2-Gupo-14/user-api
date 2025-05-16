package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateDatabase(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	result := CreateUserRepo(db)

	assert.NotNil(t, result)
}

func TestDatabase_AddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	name := "Test"
	surname := "User"
	password := "password123"
	email := "test@example.com"
	location := "Test Location"
	role := "student"
	verified := false
	profilePicture := "test_profile.jpg"
	description := "Test description"

	// Match the exact query pattern
	mock.ExpectQuery(`INSERT INTO users \(name, surname, password, email, location, role, verified, profile_photo, description\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9\) RETURNING id`).
		WithArgs(name, surname, password, email, location, role, verified, &profilePicture, description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	ctx := context.Background()
	database := CreateUserRepo(db)

	user := models.User{
		Name:         name,
		Surname:      surname,
		Password:     password,
		Email:        email,
		Location:     location,
		Role:         role,
		Verified:     verified,
		ProfilePhoto: &profilePicture,
		Description:  description,
	}

	id, err := database.AddUser(ctx, &user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, 1, id)
}

func TestDatabase_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	id := 1
	name := "Test"
	surname := "User"
	password := "password123"
	email := "test@example.com"
	location := "Test Location"
	role := "student"
	verified := false
	profilePicture := "test_profile.jpg"
	description := "Test description"
	createdAt := time.Now()
	updatedAt := time.Now()
	blocked := false

	// Use the actual query pattern from the repository
	mock.ExpectQuery(`SELECT\s+u\.id,\s*u\.name,\s*u\.surname,\s*u\.password,\s*u\.email,\s*u\.location,\s*u\.role,\s*u\.verified,\s*u\.profile_photo,\s*u\.description,\s*u\.created_at,\s*u\.updated_at,\s*EXISTS`).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "surname", "password", "email", "location", "role", "verified",
			"profile_photo", "description", "created_at", "updated_at", "blocked"}).
			AddRow(id, name, surname, password, email, location, role, verified,
				&profilePicture, description, createdAt, updatedAt, blocked))

	ctx := context.Background()
	database := CreateUserRepo(db)

	expectedUser := models.User{
		Id:           id,
		Name:         name,
		Surname:      surname,
		Password:     password,
		Email:        email,
		Location:     location,
		Role:         role,
		Verified:     verified,
		ProfilePhoto: &profilePicture,
		Description:  description,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Blocked:      blocked,
	}

	user, err := database.GetUser(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, *user)
}

func TestDatabase_GetUser_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Use the actual query pattern
	mock.ExpectQuery(`SELECT\s+u\.id,\s*u\.name`).
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	database := CreateUserRepo(db)

	_, err = database.GetUser(ctx, 1)
	assert.Error(t, err, ErrNotFound)
}

func TestDatabase_GetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	id := 1
	name := "Test"
	surname := "User"
	password := "password123"
	email := "test@example.com"
	location := "Test Location"
	role := "student"
	profilePicture := "test_profile.jpg"
	description := "Test description"
	createdAt := time.Now()
	updatedAt := time.Now()

	// Use the actual query pattern
	mock.ExpectQuery(`SELECT id, name, surname, password, email, location, role, profile_photo, description, created_at, updated_at FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "surname", "password", "email", "location", "role",
			"profile_photo", "description", "created_at", "updated_at"}).
			AddRow(id, name, surname, password, email, location, role,
				&profilePicture, description, createdAt, updatedAt))

	ctx := context.Background()
	database := CreateUserRepo(db)

	expectedUser := models.User{
		Id:           id,
		Name:         name,
		Surname:      surname,
		Password:     password,
		Email:        email,
		Location:     location,
		Role:         role,
		ProfilePhoto: &profilePicture,
		Description:  description,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	var expectedUsers []models.User
	expectedUsers = append(expectedUsers, expectedUser)

	users, err := database.GetAllUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestDatabase_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	id := 1
	name := "Test"
	surname := "User"
	password := "password123"
	email := "test@example.com"
	location := "Test Location"
	role := "student"
	verified := false
	profilePicture := "test_profile.jpg"
	description := "Test description"
	createdAt := time.Now()
	updatedAt := time.Now()
	blocked := false

	// Use the actual query pattern
	mock.ExpectQuery(`SELECT\s+u\.id,\s*u\.name`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "surname", "password", "email", "location", "role", "verified",
			"profile_photo", "description", "created_at", "updated_at", "blocked"}).
			AddRow(id, name, surname, password, email, location, role, verified,
				&profilePicture, description, createdAt, updatedAt, blocked))

	ctx := context.Background()
	database := CreateUserRepo(db)

	expectedUser := models.User{
		Id:           id,
		Name:         name,
		Surname:      surname,
		Password:     password,
		Email:        email,
		Location:     location,
		Role:         role,
		Verified:     verified,
		ProfilePhoto: &profilePicture,
		Description:  description,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Blocked:      blocked,
	}

	user, err := database.GetUserByEmail(ctx, email)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, *user)
}

func TestDatabase_GetUserByEmail_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	email := "test@example.com"

	// Use the actual query pattern
	mock.ExpectQuery(`SELECT\s+u\.id,\s*u\.name`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	database := CreateUserRepo(db)

	_, err = database.GetUserByEmail(ctx, email)
	assert.Error(t, ErrNotFound)
}

func TestDatabase_DeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	database := CreateUserRepo(db)

	err = database.DeleteUser(ctx, 1)
	assert.NoError(t, err)
}

func TestDatabase_ModifyUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	id := 1
	name := "Test"
	surname := "User"
	password := "password123"
	email := "test@example.com"
	location := "Test Location"
	role := "student"
	verified := false
	profilePicture := "test_profile.jpg"
	description := "Test description"
	blocked := false

	mock.ExpectExec(`UPDATE users SET name = \$1, surname = \$2, location = \$3, profile_photo = \$4, description = \$5, verified = \$6 WHERE id = \$7`).
		WithArgs(name, surname, location, &profilePicture, description, verified, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	database := CreateUserRepo(db)

	user := models.User{
		Id:           id,
		Name:         name,
		Surname:      surname,
		Password:     password,
		Email:        email,
		Location:     location,
		Role:         role,
		Verified:     verified,
		ProfilePhoto: &profilePicture,
		Description:  description,
		Blocked:      blocked,
	}

	err = database.ModifyUser(ctx, &user)
	assert.NoError(t, err)
}

func TestDatabase_ModifyPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`UPDATE users SET password = \$1 where id = \$2`).WithArgs("test", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	database := CreateUserRepo(db)

	err = database.ModifyPassword(ctx, 1, "test")
	assert.NoError(t, err)
}
