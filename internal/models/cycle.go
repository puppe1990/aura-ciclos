package models

import "time"

type Mode string

const (
	ModeCycle     Mode = "cycle"
	ModeTrying    Mode = "trying"
	ModePregnancy Mode = "pregnancy"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	DisplayName  string
	Anonymous    bool
	CreatedAt    time.Time
}

type Profile struct {
	UserID            int64
	DisplayName       string
	Anonymous         bool
	Mode              Mode
	PregnancyLMP      *time.Time
	TypicalCycleDays  int
	TypicalPeriodDays int
	UpdatedAt         time.Time
}

type Period struct {
	ID        int64
	UserID    int64
	StartedOn time.Time
	EndedOn   *time.Time
	Flow      string
}

type DailyLog struct {
	ID              int64
	UserID          int64
	LoggedOn        time.Time
	SymptomIDs      []string
	Mood            string
	Energy          int
	SleepHours      float64
	WaterGlasses    int
	ExerciseMinutes int
	Sex             string
	BBTCelsius      *float64
	OvulationTest   string
	PregnancyTest   string
	Notes           string
}

type Reminder struct {
	ID      int64
	UserID  int64
	Kind    string
	Hour    int
	Minute  int
	Enabled bool
}

type MayaMessage struct {
	ID        int64
	UserID    int64
	Role      string
	Body      string
	CreatedAt time.Time
}
