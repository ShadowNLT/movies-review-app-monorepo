package users

import (
	"cinepulse.nlt.net/internal/data/shared"
	"cinepulse.nlt.net/internal/data/users/inputs"
	usersShared "cinepulse.nlt.net/internal/data/users/shared"
	"cinepulse.nlt.net/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID            int64                `json:"id"`
	Email         string               `json:"email"`
	Password      usersShared.Password `json:"-"`
	ProfileHandle string               `json:"profile_handle"`
	Location      string               `json:"location"`
	DateOfBirth   time.Time            `json:"date_of_birth"`
	IsProtected   bool                 `json:"is_protected"`
	IsActivated   bool                 `json:"is_activated"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	Version       int                  `json:"version"`
}

type CreatedUserOutput struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

type UserModel struct {
	DB *sql.DB
}

var (
	ErrDuplicateEmail         = errors.New("duplicate email")
	ErrDuplicateProfileHandle = errors.New("duplicate profile handle")
)

// Validation

func ValidateUser(v *validator.Validator, u *User) {
	// ProfileHandle
	usersShared.ValidateProfileHandle(v, u.ProfileHandle)

	// E-mail
	usersShared.ValidateEmail(v, u.Email)

	// Password
	if u.Password.Plaintext != nil {
		usersShared.ValidatePasswordPlaintext(v, *u.Password.Plaintext)
	}

	// IF the Password hash is ever nil, we need to panic as this should never happen
	if u.Password.Hash == nil {
		panic("missing Password hash for user")
	}

	// Location
	v.RequiredString(u.Location, "location")

	// Date of birth
	v.AddErrorIfNot(usersShared.IsAtLeast13YearsOld(u.DateOfBirth), "date_of_birth", "must be at least 13 years old")
	v.AddErrorIfNot(usersShared.IsAtMost100YearsOld(u.DateOfBirth), "date_of_birth", "must be at most 100 years old")
}

type UserSearchByProperty = string

var (
	ID    UserSearchByProperty = "id"
	Email UserSearchByProperty = "email"
)

func (m UserModel) Insert(user *inputs.CreateUserInput) (*CreatedUserOutput, error) {
	query := `
         INSERT INTO users (email, password_hash, handle, location, date_of_birth)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at, version`
	args := []any{user.Email, user.Password.Hash, user.ProfileHandle, user.Location, user.DateOfBirth}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var createdUser CreatedUserOutput
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&createdUser.ID, &createdUser.CreatedAt, &createdUser.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail

		case err.Error() == `pq: duplicate key value violates unique constraint "users_handle_key"`:
			return nil, ErrDuplicateProfileHandle
		default:
			return nil, err
		}
	}
	return &createdUser, nil
}

func (m UserModel) GetByEmailOrId(property UserSearchByProperty, value any) (*User, error) {
	var propertyQueryString string
	switch property {
	case ID:
		propertyQueryString = "id"
		_, ok := value.(int64)
		if !ok {
			// We should panic here because this should never be the case
			panic("invalid value for ID")
		}

	case Email:
		propertyQueryString = "email"
		_, ok := value.(string)
		if !ok {
			// We should panic here because this should never be the case
			panic("invalid value for ID")
		}
	default:
		// We should panic here because we shall never be in this case
		panic("unknown property type")
	}

	query := fmt.Sprintf(`
         SELECT id, email, password_hash, handle, location, date_of_birth, is_protected, is_activated, created_at, updated_at, version
         FROM users
         WHERE %s = $1`, propertyQueryString)

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, value).Scan(
		&user.ID,
		&user.Email,
		&user.Password.Hash,
		&user.ProfileHandle,
		&user.Location,
		&user.DateOfBirth,
		&user.IsProtected,
		&user.IsActivated,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, shared.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(input *inputs.UpdateUserInput, version int, userId int64) error {
	var fields []string
	var args []any
	argPos := 1
	incrementVersion := false

	if input.Email != nil {
		fields = append(fields, fmt.Sprintf("email = $%d", argPos))
		args = append(args, *input.Email)
		argPos++
		incrementVersion = true
	}

	if input.ProfileHandle != nil {
		fields = append(fields, fmt.Sprintf("handle = $%d", argPos))
		args = append(args, *input.ProfileHandle)
		argPos++
		incrementVersion = true
	}

	if input.Location != nil {
		fields = append(fields, fmt.Sprintf("location = $%d", argPos))
		args = append(args, *input.Location)
		argPos++
		incrementVersion = true
	}

	if input.IsProtected != nil {
		fields = append(fields, fmt.Sprintf("is_protected = $%d", argPos))
		args = append(args, *input.IsProtected)
		argPos++
	}

	fields = append(fields, "updated_at = now()")
	if incrementVersion {
		fields = append(fields, "version = version + 1")
	}

	query := fmt.Sprintf(`
							 UPDATE users SET %s WHERE id = $%d AND version = $%d
							 RETURNING id, email, handle, location, date_of_birth, is_protected, is_activated, created_at, updated_at, version
							 `, strings.Join(fields, ", "), argPos, argPos+1)
	args = append(args, userId, version)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.ProfileHandle,
		&user.Location,
		&user.DateOfBirth,
		&user.IsProtected,
		&user.IsActivated,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail

		case err.Error() == `pq: duplicate key value violates unique constraint "users_handle_key"`:
			return ErrDuplicateProfileHandle

		case errors.Is(err, sql.ErrNoRows):
			return shared.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
