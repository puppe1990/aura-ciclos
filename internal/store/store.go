package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/puppe1990/amarra-cais/pkg/cais/devlog"
	"github.com/puppe1990/amarra-cais/pkg/cais/session"
	caissqlite "github.com/puppe1990/amarra-cais/pkg/cais/sqlite"
	"github.com/puppe1990/amarra-cais/pkg/cais/sqllog"
	_ "modernc.org/sqlite"

	"github.com/puppe1990/aura-ciclos/internal/models"
)

var ErrEmailTaken = errors.New("email already registered")

type Store interface {
	FindUserByEmail(email string) (models.User, error)
	FindUser(id int64) (models.User, error)
	CreateUser(email, passwordHash string) (int64, error)
	MarkAnonymous(id int64, anonymous bool) error
	CreatePasswordResetToken(userID int64) (string, error)
	ClearPasswordResetTokens(userID int64) error
	FindPasswordResetUserID(token string) (int64, bool)
	ResetPasswordWithToken(token, passwordHash string) error
	ChangePassword(id int64, passwordHash string) error
	GetProfile(userID int64) (models.Profile, error)
	SaveProfile(p models.Profile) error
	UpsertPeriod(p models.Period) error
	ListPeriods(userID int64) ([]models.Period, error)
	DeletePeriod(userID, id int64) error
	UpsertDailyLog(l models.DailyLog) error
	GetDailyLog(userID int64, day time.Time) (models.DailyLog, error)
	ListDailyLogs(userID int64, since time.Time) ([]models.DailyLog, error)
	UpsertReminder(r models.Reminder) error
	ListReminders(userID int64) ([]models.Reminder, error)
	DeleteReminder(userID, id int64) error
	AddMayaMessage(userID int64, role, body string) error
	ListMayaMessages(userID int64, limit int) ([]models.MayaMessage, error)
	Sessions() session.Store
	Ping() error
	DB() *sql.DB
	Close() error
}

type SQLiteStore struct {
	db *sqllog.DB
}

func NewSQLiteStore(dsn string, env string) (*SQLiteStore, error) {
	if dsn != ":memory:" {
		dir := filepath.Dir(dsn)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create db dir: %w", err)
		}
	}

	db, err := sql.Open("sqlite", caissqlite.DSN(dsn))
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := caissqlite.Configure(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("configure sqlite: %w", err)
	}

	if err := applyMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	cfg := sqllog.ConfigForEnv(env)
	if cfg.Enabled {
		cfg.Writer = devlog.MirrorDefault(os.Stdout)
	}
	wrapped := sqllog.Wrap(db, cfg)
	if err := seedAuthData(wrapped.Raw(), env); err != nil {
		_ = wrapped.Close()
		return nil, err
	}
	return &SQLiteStore{db: wrapped}, nil
}

func seedAuthData(db *sql.DB, env string) error {
	if env != "development" {
		return nil
	}
	if err := session.EnsureSQLiteSchema(db); err != nil {
		return err
	}
	hash, err := session.HashPassword("password")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT OR IGNORE INTO users (email, password_hash, is_admin) VALUES (?, ?, 1)", "demo@example.com", hash)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT OR IGNORE INTO profiles (user_id, display_name)
		SELECT id, 'Sofia' FROM users WHERE email = 'demo@example.com'`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE users SET display_name = 'Sofia' WHERE email = 'demo@example.com' AND display_name = ''`)
	if err != nil {
		return err
	}
	starts := []struct{ start, end string }{
		{"2026-05-04", "2026-05-08"},
		{"2026-06-01", "2026-06-05"},
		{"2026-06-29", "2026-07-03"},
		{"2026-07-27", "2026-07-31"},
		{"2026-08-24", "2026-08-28"},
		{"2026-09-21", "2026-09-25"},
		{"2026-10-19", "2026-10-23"},
	}
	for _, p := range starts {
		_, err = db.Exec(`INSERT OR IGNORE INTO periods (user_id, started_on, ended_on, flow)
			SELECT id, ?, ?, 'medio' FROM users WHERE email = 'demo@example.com'`, p.start, p.end)
		if err != nil {
			return err
		}
	}
	return nil
}

func scanUser(row interface{ Scan(dest ...any) error }) (models.User, error) {
	var u models.User
	var anon int
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &anon, &u.CreatedAt)
	if err != nil {
		return models.User{}, err
	}
	u.Anonymous = anon == 1
	return u, nil
}

func (s *SQLiteStore) FindUserByEmail(email string) (models.User, error) {
	u, err := scanUser(s.db.QueryRow(
		"SELECT id, email, password_hash, display_name, anonymous, created_at FROM users WHERE email = ?",
		email,
	))
	if err != nil {
		return models.User{}, fmt.Errorf("find user: %w", err)
	}
	return u, nil
}

func (s *SQLiteStore) FindUser(id int64) (models.User, error) {
	u, err := scanUser(s.db.QueryRow(
		"SELECT id, email, password_hash, display_name, anonymous, created_at FROM users WHERE id = ?",
		id,
	))
	if err != nil {
		return models.User{}, fmt.Errorf("find user: %w", err)
	}
	return u, nil
}

func (s *SQLiteStore) CreateUser(email, passwordHash string) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, passwordHash,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return 0, ErrEmailTaken
		}
		return 0, fmt.Errorf("create user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("create user id: %w", err)
	}
	if _, err := s.db.Exec("INSERT INTO profiles (user_id) VALUES (?)", id); err != nil {
		return 0, fmt.Errorf("create profile: %w", err)
	}
	return id, nil
}

func (s *SQLiteStore) MarkAnonymous(id int64, anonymous bool) error {
	flag := 0
	if anonymous {
		flag = 1
	}
	if _, err := s.db.Exec("UPDATE users SET anonymous = ? WHERE id = ?", flag, id); err != nil {
		return fmt.Errorf("mark anonymous user: %w", err)
	}
	if _, err := s.db.Exec("UPDATE profiles SET anonymous = ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ?", flag, id); err != nil {
		return fmt.Errorf("mark anonymous profile: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ChangePassword(id int64, passwordHash string) error {
	tx, err := s.db.Raw().Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.Exec("UPDATE users SET password_hash = ? WHERE id = ?", passwordHash, id); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM sessions WHERE user_id = ?", id); err != nil {
		return fmt.Errorf("revoke sessions: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit change password: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Sessions() session.Store {
	return session.NewSQLiteStore(s.db.Raw())
}

func (s *SQLiteStore) Ping() error {
	return s.db.Raw().Ping()
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.db.Raw()
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
