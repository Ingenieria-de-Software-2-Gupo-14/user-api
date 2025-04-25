package repositories

import (
	"context"
	"database/sql"
	"ing-soft-2-tp1/internal/errors"
	"ing-soft-2-tp1/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_USERNAME        = "testUser"
	TEST_NAME            = "testName"
	TEST_SURNAME         = "testSurname"
	TEST_PASSWORD        = "testPassword"
	TEST_EMAIL           = "testEmail"
	TEST_LOCATION        = "testLocation"
	TEST_ADMIN           = false
	TEST_BLOCKED         = false
	TEST_PROFILE_PICTURE = 0
	TEST_DESCRIPTION     = "testDesc"
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

	//rows := sqlmock.NewRows([]string{"id", "username", "name", "surname", "password", "email", "location", "admin", "bolcked_user", "profile_photo", "description"})

	mock.ExpectExec("INSERT INTO users").WithArgs(TEST_USERNAME, TEST_NAME, TEST_SURNAME, TEST_PASSWORD, TEST_EMAIL, TEST_LOCATION, TEST_ADMIN, TEST_BLOCKED, TEST_PROFILE_PICTURE, TEST_DESCRIPTION).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT id FROM users WHERE email = \$1`).WithArgs(TEST_EMAIL).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))

	ctx := context.Background()

	database := CreateUserRepo(db)

	user := models.User{
		Username:     TEST_USERNAME,
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     TEST_LOCATION,
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  TEST_DESCRIPTION,
	}

	id, err := database.AddUser(ctx, &user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, 1, id) //TODO: change adduser to no longer return the id, its not necessary
}

func TestDatabase_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1`).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "surname", "password", "email", "location", "admin", "bolcked_user", "profile_photo", "description"}).
			AddRow(1, TEST_USERNAME, TEST_NAME, TEST_SURNAME, TEST_PASSWORD, TEST_EMAIL, TEST_LOCATION, TEST_ADMIN, TEST_BLOCKED, TEST_PROFILE_PICTURE, TEST_DESCRIPTION))

	ctx := context.Background()

	database := CreateUserRepo(db)

	expectedUser := models.User{
		Id:           1,
		Username:     TEST_USERNAME,
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     TEST_LOCATION,
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  TEST_DESCRIPTION,
	}

	user, err := database.GetUser(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, *user)
}

func TestDatabase_GetUser_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1`).WithArgs(1).WillReturnError(sql.ErrNoRows)

	ctx := context.Background()

	database := CreateUserRepo(db)

	_, err = database.GetUser(ctx, 1)
	assert.Error(t, err, errors.ErrNotFound)
}

func TestDatabase_GetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	mock.ExpectQuery("SELECT id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description FROM users").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "surname", "password", "email", "location", "admin", "bolcked_user", "profile_photo", "description"}).
			AddRow(1, TEST_USERNAME, TEST_NAME, TEST_SURNAME, TEST_PASSWORD, TEST_EMAIL, TEST_LOCATION, TEST_ADMIN, TEST_BLOCKED, TEST_PROFILE_PICTURE, TEST_DESCRIPTION))

	ctx := context.Background()

	database := CreateUserRepo(db)

	expectedUser := models.User{
		Id:           1,
		Username:     TEST_USERNAME,
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     TEST_LOCATION,
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  TEST_DESCRIPTION,
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

	mock.ExpectQuery(`SELECT \* FROM users WHERE email ILIKE \$1`).WithArgs(TEST_EMAIL).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "name", "surname", "password", "email", "location", "admin", "bolcked_user", "profile_photo", "description"}).
			AddRow(1, TEST_USERNAME, TEST_NAME, TEST_SURNAME, TEST_PASSWORD, TEST_EMAIL, TEST_LOCATION, TEST_ADMIN, TEST_BLOCKED, TEST_PROFILE_PICTURE, TEST_DESCRIPTION))

	ctx := context.Background()

	database := CreateUserRepo(db)

	expectedUser := models.User{
		Id:           1,
		Username:     TEST_USERNAME,
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     TEST_LOCATION,
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  TEST_DESCRIPTION,
	}

	user, err := database.GetUserByEmail(ctx, TEST_EMAIL)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, *user)
}

func TestDatabase_GetUserByEmail_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT \* FROM users WHERE email ILIKE \$1`).WithArgs(TEST_EMAIL).WillReturnError(sql.ErrNoRows)

	ctx := context.Background()

	database := CreateUserRepo(db)

	_, err = database.GetUserByEmail(ctx, TEST_EMAIL)
	assert.Error(t, errors.ErrNotFound)
}

func TestDatabase_DeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	database := CreateUserRepo(db)

	err = database.DeleteUser(ctx, 1)
	assert.NoError(t, err)
}

func TestDatabase_BlockUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`UPDATE users SET blocked_user = true where id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	database := CreateUserRepo(db)

	err = database.BlockUser(ctx, 1)
	assert.NoError(t, err)
}

func TestDatabase_ModifyLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`UPDATE users SET location = \$1 where id = \$2`).WithArgs(TEST_LOCATION, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	database := CreateUserRepo(db)

	err = database.ModifyLocation(ctx, 1, TEST_LOCATION)
	assert.NoError(t, err)
}

func TestDatabase_ModifyUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`UPDATE users SET username = \$1, name= \$2, surname=\$3,  location=\$4, profile_photo=\$5,description=\$6 WHERE id = \$7`).WithArgs(TEST_USERNAME, TEST_NAME, TEST_SURNAME, TEST_LOCATION, TEST_PROFILE_PICTURE, TEST_DESCRIPTION, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	database := CreateUserRepo(db)

	user := models.User{
		Id:           1,
		Username:     TEST_USERNAME,
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     TEST_LOCATION,
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  TEST_DESCRIPTION,
	}

	err = database.ModifyUser(ctx, &user)
	assert.NoError(t, err)
}
