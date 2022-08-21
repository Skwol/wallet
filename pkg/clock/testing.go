package clock

import "time"

type Fake struct {
	now time.Time
}

func NewFake(now time.Time) SettableClock {
	return &Fake{
		now: now,
	}
}

func (f Fake) Now() time.Time {
	return f.now
}

func (f *Fake) SetTime(t time.Time) {
	f.now = t
}
