package repositories

import (
	"context"
	"database/sql"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/errors"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	_ "github.com/lib/pq"
)

const (
	badLoginAttemptWindow = "-15 minutes"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id int) error
	AddUser(ctx context.Context, user *models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ModifyUser(ctx context.Context, user *models.User) error
	ModifyLocation(ctx context.Context, id int, newLocation string) error
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
			u.id, u.name, u.surname, u.password, u.email, u.location, u.admin,
			u.profile_photo, u.description, u.phone, u.created_at, u.updated_at,
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
		&user.Admin, &user.ProfilePhoto, &user.Description, &user.Phone, &user.CreatedAt,
		&user.UpdatedAt, &user.Blocked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	// Simplificado: No calcular BadLoginAttempts ni Blocked aquí por rendimiento.
	// Estos campos tendrán su valor cero (0 y false).
	rows, err := db.DB.QueryContext(ctx, `
		SELECT id, name, surname, password, email, location, admin, profile_photo,
		       description, phone, created_at, updated_at
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
			&user.Admin, &user.ProfilePhoto, &user.Description, &user.Phone, &user.CreatedAt,
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
	// La consulta INSERT no necesita cambios relacionados con bad_login_attempts
	query := `
		INSERT INTO users (name, surname, password, email, location, admin, profile_photo, description, phone)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, query,
		&user.Name, &user.Surname, &user.Password, &user.Email, &user.Location, &user.Admin, &user.ProfilePhoto, &user.Description, &user.Phone,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT
			u.id, u.name, u.surname, u.password, u.email, u.location, u.admin,
			u.profile_photo, u.description, u.phone, u.created_at, u.updated_at,
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
	// Eliminada la columna bad_login_attempts del Scan
	err := row.Scan(
		&user.Id, &user.Name, &user.Surname, &user.Password, &user.Email, &user.Location,
		&user.Admin, &user.ProfilePhoto, &user.Description, &user.Phone, &user.CreatedAt,
		&user.UpdatedAt, &user.Blocked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db userRepository) ModifyUser(ctx context.Context, user *models.User) error {
	// No modificar bad_login_attempts ni blocked aquí directamente
	query := `
		UPDATE users SET name = $1, surname = $2, location = $3, profile_photo = $4, description = $5, phone = $6
		WHERE id = $7`
	_, err := db.DB.ExecContext(ctx, query,
		&user.Name, &user.Surname, &user.Location, &user.ProfilePhoto, &user.Description, &user.Phone, &user.Id,
	)
	// updated_at se actualiza por trigger
	return err
}

func (db userRepository) ModifyLocation(ctx context.Context, id int, newLocation string) error {
	// updated_at se actualizará automáticamente por el trigger
	_, err := db.DB.ExecContext(ctx, "UPDATE users SET location = $1 where id = $2", newLocation, id)
	return err
}

// Eliminar Implementación de métodos IncrementBadLoginAttempts y ResetBadLoginAttempts
// func (db userRepository) IncrementBadLoginAttempts(...) { ... }
// func (db userRepository) ResetBadLoginAttempts(...) { ... }
