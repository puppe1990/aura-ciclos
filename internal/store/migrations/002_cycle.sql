-- up
ALTER TABLE users ADD COLUMN display_name TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN anonymous INTEGER NOT NULL DEFAULT 0 CHECK (anonymous IN (0, 1));

CREATE TABLE IF NOT EXISTS profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT NOT NULL DEFAULT '',
    anonymous INTEGER NOT NULL DEFAULT 0 CHECK (anonymous IN (0, 1)),
    mode TEXT NOT NULL DEFAULT 'cycle' CHECK (mode IN ('cycle', 'trying', 'pregnancy')),
    pregnancy_lmp TEXT,
    typical_cycle_days INTEGER NOT NULL DEFAULT 28,
    typical_period_days INTEGER NOT NULL DEFAULT 5,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS periods (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_on TEXT NOT NULL,
    ended_on TEXT,
    flow TEXT NOT NULL DEFAULT 'medio',
    UNIQUE (user_id, started_on)
);

CREATE TABLE IF NOT EXISTS daily_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    logged_on TEXT NOT NULL,
    symptom_ids TEXT NOT NULL DEFAULT '[]',
    mood TEXT NOT NULL DEFAULT '',
    energy INTEGER NOT NULL DEFAULT 0,
    sleep_hours REAL NOT NULL DEFAULT 0,
    water_glasses INTEGER NOT NULL DEFAULT 0,
    exercise_minutes INTEGER NOT NULL DEFAULT 0,
    sex TEXT NOT NULL DEFAULT '',
    bbt_celsius REAL,
    ovulation_test TEXT NOT NULL DEFAULT '',
    pregnancy_test TEXT NOT NULL DEFAULT '',
    notes TEXT NOT NULL DEFAULT '',
    UNIQUE (user_id, logged_on)
);

CREATE TABLE IF NOT EXISTS reminders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('pill', 'water', 'period')),
    hour INTEGER NOT NULL,
    minute INTEGER NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS maya_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('user', 'maya')),
    body TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
