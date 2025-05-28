package repositories

import (
	"context"
	"database/sql"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	_ "github.com/lib/pq"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id int) error
	AddUser(ctx context.Context, user *models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ModifyUser(ctx context.Context, user *models.User) error
	ModifyPassword(ctx context.Context, id int, password string) error
	AddNotificationToken(ctx context.Context, id int, text string) error
	GetUserNotificationsToken(ctx context.Context, id int) (models.Notifications, error)
	SetVerifiedTrue(ctx context.Context, id int) error
}

type userRepository struct {
	DB *sql.DB
}

// CreateUserRepo creates and returns a database
func CreateUserRepo(db *sql.DB) *userRepository {
	return &userRepository{DB: db}
}

func (db userRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT
			u.id, u.name, u.surname, u.password, u.email, u.location, u.role, u.verified,
			u.profile_photo, u.description, u.created_at, u.updated_at,
			EXISTS(
				SELECT 1 FROM blocked_users
				WHERE blocked_user_id = u.id
				AND (blocked_until IS NULL OR blocked_until > NOW())
			) AS blocked
		FROM users u
		WHERE u.id = $1
		LIMIT 1`

	row := db.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.Id, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location,
		&user.Role, &user.Verified, &user.ProfilePhoto, &user.Description, &user.CreatedAt,
		&user.UpdatedAt, &user.Blocked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	// Simplificado: No calcular BadLoginAttempts ni Blocked aquí por rendimiento.
	// Estos campos tendrán su valor cero (0 y false).
	rows, err := db.DB.QueryContext(ctx, `
		SELECT id, name, surname, password, email, location, role, profile_photo,
		       description, created_at, updated_at
		FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		// Scan sin bad_login_attempts ni blocked
		err := rows.Scan(
			&user.Id, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location,
			&user.Role, &user.ProfilePhoto, &user.Description, &user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db userRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (db userRepository) AddUser(ctx context.Context, user *models.User) (int, error) {
	query := `
		INSERT INTO users (name, surname, password, email, location, role, verified, profile_photo, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, query,
		&user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Role,
		&user.Verified, &user.ProfilePhoto, &user.Description,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT
			u.id, u.name, u.surname, u.password, u.email, u.location, u.role, u.verified,
			u.profile_photo, u.description, u.created_at, u.updated_at,
			EXISTS(
				SELECT 1 FROM blocked_users
				WHERE blocked_user_id = u.id
				AND (blocked_until IS NULL OR blocked_until > NOW())
			) AS blocked
		FROM users u
		WHERE u.email ILIKE $1
		LIMIT 1`

	row := db.DB.QueryRowContext(ctx, query, email)
	var user models.User

	err := row.Scan(
		&user.Id, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location,
		&user.Role, &user.Verified, &user.ProfilePhoto, &user.Description, &user.CreatedAt,
		&user.UpdatedAt, &user.Blocked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db userRepository) ModifyUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET name = $1, surname = $2, location = $3, profile_photo = $4, description = $5, verified = $6
		WHERE id = $7`
	_, err := db.DB.ExecContext(ctx, query,
		&user.Name, &user.Surname, &user.Location, &user.ProfilePhoto, &user.Description, &user.Verified, &user.Id,
	)

	return err
}

func (db userRepository) ModifyPassword(ctx context.Context, id int, password string) error {
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET password = $1 where id = $2", password, id)
	return err
}

func (db userRepository) AddNotificationToken(ctx context.Context, id int, text string) error {
	query := `
		INSERT INTO notifications (user_id, token)
		VALUES ($1, $2)`
	row := db.DB.QueryRowContext(ctx, query, id, text)
	return row.Err()
}

func (db userRepository) GetUserNotificationsToken(ctx context.Context, id int) (models.Notifications, error) {
	query := `
			SELECT token, created_time
			FROM notifications
			WHERE user_id = $1`
	rows, err := db.DB.QueryContext(ctx, query, id)
	if err != nil {
		return models.Notifications{}, err
	}

	var notifications models.Notifications

	for rows.Next() {
		var n models.Notification
		err := rows.Scan(&n.NotificationText, &n.CreatedTime)
		if err != nil {
			return models.Notifications{}, err
		}
		notifications.Notifications = append(notifications.Notifications, n)
	}

	if err := rows.Err(); err != nil {
		return models.Notifications{}, err
	}
	return notifications, nil
}

func (db userRepository) SetVerifiedTrue(ctx context.Context, id int) error {
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET verified = true where id = $1", id)
	return err
}

//id, username, name, surname,  password,email, location, admin, blocked_user, profile_photo,description
// Eliminar Implementación de métodos IncrementBadLoginAttempts y ResetBadLoginAttempts
// func (db userRepository) IncrementBadLoginAttempts(...) { ... }
// func (db userRepository) ResetBadLoginAttempts(...) { ... }
