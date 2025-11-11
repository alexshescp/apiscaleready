package domain

// Progress describes minimal progress reporting capabilities required by
// runners. Implementations may be backed by a real progress bar or act as a
// no-op in environments where visual progress does not make sense (e.g. HTTP
// handlers or tests).
type Progress interface {
	Add(int) error
	Finish() error
}

// NoopProgress is a Progress implementation that ignores all updates. It is
// helpful for automated environments where progress visualisation would be
// noise.
type NoopProgress struct{}

// Add implements the Progress interface.
func (NoopProgress) Add(int) error { return nil }

// Finish implements the Progress interface.
func (NoopProgress) Finish() error { return nil }

// NewNoopProgress returns a Progress implementation that performs no action.
func NewNoopProgress() Progress {
	return NoopProgress{}
}
