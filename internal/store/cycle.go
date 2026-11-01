package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/puppe1990/aura-ciclos/internal/models"
)

func (s *SQLiteStore) GetProfile(userID int64) (models.Profile, error) {
	var p models.Profile
	var anon int
	var lmp sql.NullString
	err := s.db.QueryRow(
		`SELECT user_id, display_name, anonymous, mode, pregnancy_lmp, typical_cycle_days, typical_period_days, updated_at
		 FROM profiles WHERE user_id = ?`,
		userID,
	).Scan(&p.UserID, &p.DisplayName, &anon, &p.Mode, &lmp, &p.TypicalCycleDays, &p.TypicalPeriodDays, &p.UpdatedAt)
	if err != nil {
		return models.Profile{}, fmt.Errorf("get profile: %w", err)
	}
	p.Anonymous = anon == 1
	if lmp.Valid {
		day, parseErr := parseDay(lmp.String)
		if parseErr != nil {
			return models.Profile{}, fmt.Errorf("get profile lmp %q: %w", lmp.String, parseErr)
		}
		p.PregnancyLMP = &day
	}
	return p, nil
}

func (s *SQLiteStore) SaveProfile(p models.Profile) error {
	var lmp any
	if p.PregnancyLMP != nil {
		lmp = formatDay(*p.PregnancyLMP)
	}
	anon := 0
	if p.Anonymous {
		anon = 1
	}
	if p.Mode == "" {
		p.Mode = models.ModeCycle
	}
	if p.TypicalCycleDays == 0 {
		p.TypicalCycleDays = 28
	}
	if p.TypicalPeriodDays == 0 {
		p.TypicalPeriodDays = 5
	}
	_, err := s.db.Exec(
		`INSERT INTO profiles (user_id, display_name, anonymous, mode, pregnancy_lmp, typical_cycle_days, typical_period_days, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(user_id) DO UPDATE SET
		   display_name = excluded.display_name,
		   anonymous = excluded.anonymous,
		   mode = excluded.mode,
		   pregnancy_lmp = excluded.pregnancy_lmp,
		   typical_cycle_days = excluded.typical_cycle_days,
		   typical_period_days = excluded.typical_period_days,
		   updated_at = CURRENT_TIMESTAMP`,
		p.UserID, p.DisplayName, anon, string(p.Mode), lmp, p.TypicalCycleDays, p.TypicalPeriodDays,
	)
	if err != nil {
		return fmt.Errorf("save profile: %w", err)
	}
	_, err = s.db.Exec("UPDATE users SET display_name = ?, anonymous = ? WHERE id = ?", p.DisplayName, anon, p.UserID)
	if err != nil {
		return fmt.Errorf("save profile user: %w", err)
	}
	return nil
}

func (s *SQLiteStore) UpsertPeriod(p models.Period) error {
	var ended any
	if p.EndedOn != nil {
		ended = formatDay(*p.EndedOn)
	}
	if p.Flow == "" {
		p.Flow = "medio"
	}
	_, err := s.db.Exec(
		`INSERT INTO periods (user_id, started_on, ended_on, flow) VALUES (?, ?, ?, ?)
		 ON CONFLICT(user_id, started_on) DO UPDATE SET ended_on = excluded.ended_on, flow = excluded.flow`,
		p.UserID, formatDay(p.StartedOn), ended, p.Flow,
	)
	if err != nil {
		return fmt.Errorf("upsert period: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListPeriods(userID int64) ([]models.Period, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, started_on, ended_on, flow FROM periods WHERE user_id = ? ORDER BY started_on ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list periods: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.Period
	for rows.Next() {
		var p models.Period
		var start string
		var end sql.NullString
		if err := rows.Scan(&p.ID, &p.UserID, &start, &end, &p.Flow); err != nil {
			return nil, fmt.Errorf("scan period: %w", err)
		}
		day, err := parseDay(start)
		if err != nil {
			return nil, err
		}
		p.StartedOn = day
		if end.Valid {
			e, err := parseDay(end.String)
			if err != nil {
				return nil, err
			}
			p.EndedOn = &e
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) DeletePeriod(userID, id int64) error {
	_, err := s.db.Exec("DELETE FROM periods WHERE user_id = ? AND id = ?", userID, id)
	if err != nil {
		return fmt.Errorf("delete period: %w", err)
	}
	return nil
}

func (s *SQLiteStore) UpsertDailyLog(l models.DailyLog) error {
	ids, err := json.Marshal(l.SymptomIDs)
	if err != nil {
		return fmt.Errorf("marshal symptoms: %w", err)
	}
	if l.SymptomIDs == nil {
		ids = []byte("[]")
	}
	var bbt any
	if l.BBTCelsius != nil {
		bbt = *l.BBTCelsius
	}
	_, err = s.db.Exec(
		`INSERT INTO daily_logs (
			user_id, logged_on, symptom_ids, mood, energy, sleep_hours, water_glasses,
			exercise_minutes, sex, sex_times, orgasm, bbt_celsius, ovulation_test, pregnancy_test, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(user_id, logged_on) DO UPDATE SET
			symptom_ids = excluded.symptom_ids,
			mood = excluded.mood,
			energy = excluded.energy,
			sleep_hours = excluded.sleep_hours,
			water_glasses = excluded.water_glasses,
			exercise_minutes = excluded.exercise_minutes,
			sex = excluded.sex,
			sex_times = excluded.sex_times,
			orgasm = excluded.orgasm,
			bbt_celsius = excluded.bbt_celsius,
			ovulation_test = excluded.ovulation_test,
			pregnancy_test = excluded.pregnancy_test,
			notes = excluded.notes`,
		l.UserID, formatDay(l.LoggedOn), string(ids), l.Mood, l.Energy, l.SleepHours, l.WaterGlasses,
		l.ExerciseMinutes, l.Sex, l.SexTimes, l.Orgasm, bbt, l.OvulationTest, l.PregnancyTest, l.Notes,
	)
	if err != nil {
		return fmt.Errorf("upsert daily log: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetDailyLog(userID int64, day time.Time) (models.DailyLog, error) {
	row := s.db.QueryRow(
		`SELECT id, user_id, logged_on, symptom_ids, mood, energy, sleep_hours, water_glasses,
		        exercise_minutes, sex, sex_times, orgasm, bbt_celsius, ovulation_test, pregnancy_test, notes
		 FROM daily_logs WHERE user_id = ? AND logged_on = ?`,
		userID, formatDay(day),
	)
	l, err := scanDailyLog(row)
	if err != nil {
		return models.DailyLog{}, fmt.Errorf("get daily log: %w", err)
	}
	return l, nil
}

func (s *SQLiteStore) ListDailyLogs(userID int64, since time.Time) ([]models.DailyLog, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, logged_on, symptom_ids, mood, energy, sleep_hours, water_glasses,
		        exercise_minutes, sex, sex_times, orgasm, bbt_celsius, ovulation_test, pregnancy_test, notes
		 FROM daily_logs WHERE user_id = ? AND logged_on >= ? ORDER BY logged_on ASC`,
		userID, formatDay(since),
	)
	if err != nil {
		return nil, fmt.Errorf("list daily logs: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.DailyLog
	for rows.Next() {
		l, err := scanDailyLog(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) UpsertReminder(r models.Reminder) error {
	if r.ID == 0 {
		_, err := s.db.Exec(
			`INSERT INTO reminders (user_id, kind, hour, minute, enabled) VALUES (?, ?, ?, ?, ?)`,
			r.UserID, r.Kind, r.Hour, r.Minute, boolInt(r.Enabled),
		)
		if err != nil {
			return fmt.Errorf("insert reminder: %w", err)
		}
		return nil
	}
	_, err := s.db.Exec(
		`UPDATE reminders SET kind = ?, hour = ?, minute = ?, enabled = ? WHERE id = ? AND user_id = ?`,
		r.Kind, r.Hour, r.Minute, boolInt(r.Enabled), r.ID, r.UserID,
	)
	if err != nil {
		return fmt.Errorf("update reminder: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListReminders(userID int64) ([]models.Reminder, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, kind, hour, minute, enabled FROM reminders WHERE user_id = ? ORDER BY hour, minute`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list reminders: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.Reminder
	for rows.Next() {
		var r models.Reminder
		var enabled int
		if err := rows.Scan(&r.ID, &r.UserID, &r.Kind, &r.Hour, &r.Minute, &enabled); err != nil {
			return nil, fmt.Errorf("scan reminder: %w", err)
		}
		r.Enabled = enabled == 1
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) DeleteReminder(userID, id int64) error {
	_, err := s.db.Exec("DELETE FROM reminders WHERE user_id = ? AND id = ?", userID, id)
	if err != nil {
		return fmt.Errorf("delete reminder: %w", err)
	}
	return nil
}

func (s *SQLiteStore) AddMayaMessage(userID int64, role, body string) error {
	_, err := s.db.Exec(
		`INSERT INTO maya_messages (user_id, role, body) VALUES (?, ?, ?)`,
		userID, role, body,
	)
	if err != nil {
		return fmt.Errorf("add maya message: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ListMayaMessages(userID int64, limit int) ([]models.MayaMessage, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Query(
		`SELECT id, user_id, role, body, created_at FROM maya_messages
		 WHERE user_id = ? ORDER BY id DESC LIMIT ?`,
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list maya messages: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var newestFirst []models.MayaMessage
	for rows.Next() {
		var m models.MayaMessage
		if err := rows.Scan(&m.ID, &m.UserID, &m.Role, &m.Body, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan maya message: %w", err)
		}
		newestFirst = append(newestFirst, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	out := make([]models.MayaMessage, len(newestFirst))
	for i, m := range newestFirst {
		out[len(newestFirst)-1-i] = m
	}
	return out, nil
}

func scanDailyLog(row interface{ Scan(dest ...any) error }) (models.DailyLog, error) {
	var l models.DailyLog
	var day, symptoms string
	var bbt sql.NullFloat64
	err := row.Scan(
		&l.ID, &l.UserID, &day, &symptoms, &l.Mood, &l.Energy, &l.SleepHours, &l.WaterGlasses,
		&l.ExerciseMinutes, &l.Sex, &l.SexTimes, &l.Orgasm, &bbt, &l.OvulationTest, &l.PregnancyTest, &l.Notes,
	)
	if err != nil {
		return models.DailyLog{}, err
	}
	logged, err := parseDay(day)
	if err != nil {
		return models.DailyLog{}, err
	}
	l.LoggedOn = logged
	if symptoms != "" {
		if err := json.Unmarshal([]byte(symptoms), &l.SymptomIDs); err != nil {
			return models.DailyLog{}, fmt.Errorf("symptoms json %q: %w", symptoms, err)
		}
	}
	if bbt.Valid {
		v := bbt.Float64
		l.BBTCelsius = &v
	}
	return l, nil
}

func formatDay(t time.Time) string {
	return t.UTC().Format(time.DateOnly)
}

func parseDay(s string) (time.Time, error) {
	t, err := time.ParseInLocation(time.DateOnly, s, time.UTC)
	if err != nil {
		return time.Time{}, fmt.Errorf("day %q: want YYYY-MM-DD", s)
	}
	return t, nil
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
