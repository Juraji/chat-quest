package chat_sessions

import (
	"time"
)

type TimeOfDay string

const (
	Midnight     TimeOfDay = "MIDNIGHT"      // (00:00–01:00)
	Night        TimeOfDay = "NIGHT"         // (01:00–05:59)
	EarlyMorning TimeOfDay = "EARLY_MORNING" // (06:00–08:59)
	Morning      TimeOfDay = "MORNING"       // (09:00–11:59)
	Noon         TimeOfDay = "NOON"          // (12:00–13:00)
	Afternoon    TimeOfDay = "AFTERNOON"     // (13:00–18:00)
	Evening      TimeOfDay = "EVENING"       // (18:00–22:00)
	LateNight    TimeOfDay = "LATE_NIGHT"    // (22:00–23:59)
	RealTime     TimeOfDay = "REAL_TIME"
)

func (t *TimeOfDay) IsValid() bool {
	if t == nil {
		return true
	}

	switch *t {
	case Midnight, Night,
		EarlyMorning, Morning,
		Noon, Afternoon,
		Evening, LateNight,
		RealTime:
		return true
	}

	return false
}

func (t *TimeOfDay) HumanFmtEn() string {
	if t == nil {
		return ""
	}
	switch *t {
	case Midnight:
		return "Midnight (00:00–01:00)"
	case Night:
		return "Night time (01:00–06:00)"
	case EarlyMorning:
		return "Early morning (06:00–09:00)"
	case Morning:
		return "Morning (09:00–11:59)"
	case Noon:
		return "Noon (12:00-13:00)"
	case Afternoon:
		return "Afternoon (13:00–18:00)"
	case Evening:
		return "Evening (18:00–22:00)"
	case LateNight:
		return "Late night (22:00–23:59)"
	case RealTime:
		return time.Now().Format("15:04")
	default:
		panic("invalid timeOfDay")
	}
}
