package clock

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func New() Clock {
	return &realClock{}
}

func (*realClock) Now() time.Time {
	return time.Now()
}

type FakeClock struct {
	Clock
	now time.Time
}

func NewFake() Clock {
	return &FakeClock{now: time.Now()}
}

func NewFakeClockWithTime(now time.Time) Clock {
	return &FakeClock{now: now}
}

func NewFakeClockWithTimeS(now string) Clock {
	n, err := time.Parse(time.RFC3339, now)
	if err != nil {
		panic(err)
	}
	return &FakeClock{now: n}
}

func (f *FakeClock) Set(now time.Time) Clock {
	f.now = now
	return f
}

func (f *FakeClock) Now() time.Time {
	return f.now
}
