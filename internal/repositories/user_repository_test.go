package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestDatabase_MakeTeacher(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`UPDATE users SET role = \$1 where id = \$2`).WithArgs("teacher", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()

	database := CreateUserRepo(db)

	err = database.MakeTeacher(ctx, 1)
	assert.NoError(t, err)
}

func TestUserRepository_AddNotificationToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()
	userId := 1
	token := "123456"

	mock.ExpectQuery(`INSERT INTO notifications \(user_id, token\) VALUES \(\$1, \$2\)`).
		WithArgs(userId, token).
		WillReturnRows(sqlmock.NewRows([]string{}))

	err = repo.AddNotificationToken(ctx, userId, token)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserNotificationsToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()
	userId := 1
	now := time.Now()

	mock.ExpectQuery(`SELECT token, created_time FROM notifications WHERE user_id = \$1`).
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"token", "created_time"}).
			AddRow("token1", now).
			AddRow("token2", now))

	result, err := repo.GetUserNotificationsToken(ctx, userId)
	assert.NoError(t, err)
	assert.Len(t, result.NotificationTokens, 2)
	assert.Equal(t, "token1", result.NotificationTokens[0].NotificationToken)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetUserNotificationsToken_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"token", "created_time"}).
		AddRow("tok", time.Now())
	rows.RowError(0, ErrNotFound)

	mock.ExpectQuery(`SELECT token, created_time FROM notifications WHERE user_id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	_, err = repo.GetUserNotificationsToken(ctx, 1)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_SetVerifiedTrue(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()
	userID := 1

	mock.ExpectExec(`UPDATE users SET verified = true where id = \$1`).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.SetVerifiedTrue(ctx, userID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_AddPasswordResetToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()

	userID := 1
	email := "user@example.com"
	token := "123456"
	expiration := time.Now().Add(1 * time.Hour)

	mock.ExpectExec(`INSERT INTO password_reset \(user_id, email, token, token_expiration\) VALUES \(\$1, \$2, \$3, \$4\)`).
		WithArgs(userID, email, token, expiration).
		WillReturnResult(sqlmock.NewResult(1, 1)) // simulate 1 row inserted

	err = repo.AddPasswordResetToken(ctx, userID, email, token, expiration)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetPasswordResetTokenInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()

	token := "reset-token"
	expTime := time.Now()
	mail := "user@example.com"
	userId := 1

	mock.ExpectQuery(`SELECT email, user_id, token_expiration, used FROM password_reset WHERE token = \$1`).
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"email", "user_id", "token_expiration", "used"}).
			AddRow("user@example.com", userId, expTime, false))

	result, err := repo.GetPasswordResetTokenInfo(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, mail, result.Email)
	assert.Equal(t, userId, result.UserId)
	assert.Equal(t, expTime, result.Exp)
	assert.False(t, result.Used)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_SetPasswordTokenUsed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()
	token := "123456"

	mock.ExpectExec(`UPDATE password_reset SET used = true where token = \$1`).
		WithArgs(token).
		WillReturnResult(sqlmock.NewResult(0, 1)) // simulate success

	err = repo.SetPasswordTokenUsed(ctx, token)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_SetNotificationPreference(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()

	preference := models.NotificationPreferenceRequest{
		NotificationType:       "email_notifications",
		NotificationPreference: true,
	}
	userID := 1

	expectedQuery := fmt.Sprintf("UPDATE users SET %s = \\$2 WHERE id = \\$1", preference.NotificationType)

	mock.ExpectExec(expectedQuery).
		WithArgs(userID, preference.NotificationPreference).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SetNotificationPreference(ctx, userID, preference)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestUserRepository_CheckPreference(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()
	userID := 1
	notificationType := "exam_notifications"

	query := fmt.Sprintf("SELECT %s FROM users WHERE id = \\$1", notificationType)
	mock.ExpectQuery(query).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{notificationType}).AddRow(true))

	result, err := repo.CheckPreference(ctx, userID, notificationType)
	assert.NoError(t, err)
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetNotificationPreference(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateUserRepo(db)
	ctx := context.Background()
	userID := 1
	examNotification := true
	homeworkNotification := false
	socialNotification := true

	mock.ExpectQuery(`SELECT exam_notification, homework_notification, social_notification FROM users WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{
			"exam_notification", "homework_notification", "social_notification",
		}).AddRow(examNotification, homeworkNotification, socialNotification))

	result, err := repo.GetNotificationPreference(ctx, userID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, examNotification, result.ExamNotification)
	assert.Equal(t, homeworkNotification, result.HomeworkNotification)
	assert.Equal(t, socialNotification, result.SocialNotification)
	assert.NoError(t, mock.ExpectationsWereMet())
}
