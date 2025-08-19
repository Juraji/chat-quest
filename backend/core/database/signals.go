package database

import "juraji.nl/chat-quest/core/util/signals"

type MigratedEvent struct {
	FromVersion uint
	ToVersion   uint
}

// IsUp checks if this is an up migration event (if we're moving forward in version).
// Returns true if FromVersion < ToVersion, meaning we're migrating from a lower to higher version.
func (e *MigratedEvent) IsUp() bool {
	return e.FromVersion < e.ToVersion
}

// IsUpIncludingVersion checks if:
//  1. This is an "up" migration event (FromVersion < ToVersion)
//  2. The given version is within the range of versions that were migrated up.
//     Specifically, it checks if version > FromVersion AND version <= ToVersion.
//
// Returns true if both conditions are met.
func (e *MigratedEvent) IsUpIncludingVersion(version uint) bool {
	return e.IsUp() && version > e.FromVersion && version <= e.ToVersion
}

var MigrationsCompletedSignal = signals.New[MigratedEvent]()
