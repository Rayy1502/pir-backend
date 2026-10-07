package repository

import (
	"context"
	"database/sql"
	"errors"

	mssql "github.com/microsoft/go-mssqldb"

	"github.com/Rayy1502/pir-backend/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func scanUser(row *sql.Row) (*model.User, error) {
	var u model.User
	var uid mssql.UniqueIdentifier

	err := row.Scan(
		&uid,
		&u.NIP,
		&u.Nama,
		&u.PasswordHash,
		&u.Jabatan,
		&u.Divisi,
		&u.Unit,
		&u.Role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	u.ID = uid.String()
	return &u, nil
}

func (r *UserRepository) FindByNIP(ctx context.Context, nip string) (*model.User, error) {
	query := `
		SELECT id, nip, nama, password_hash, jabatan, divisi, unit, role, is_active, created_at, updated_at
		FROM users WHERE nip = @nip
	`
	row := r.DB.QueryRowContext(ctx, query, sql.Named("nip", nip))
	return scanUser(row)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, nip, nama, password_hash, jabatan, divisi, unit, role, is_active, created_at, updated_at
		FROM users WHERE id = @id
	`
	row := r.DB.QueryRowContext(ctx, query, sql.Named("id", id))
	return scanUser(row)
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (nip, nama, password_hash, jabatan, divisi, unit, role, is_active)
		OUTPUT INSERTED.id
		VALUES (@nip, @nama, @password_hash, @jabatan, @divisi, @unit, @role, @is_active)
	`
	var uid mssql.UniqueIdentifier
	err := r.DB.QueryRowContext(ctx, query,
		sql.Named("nip", u.NIP),
		sql.Named("nama", u.Nama),
		sql.Named("password_hash", u.PasswordHash),
		sql.Named("jabatan", u.Jabatan),
		sql.Named("divisi", u.Divisi),
		sql.Named("unit", u.Unit),
		sql.Named("role", u.Role),
		sql.Named("is_active", u.IsActive),
	).Scan(&uid)
	if err != nil {
		return err
	}
	u.ID = uid.String()
	return nil
}
