package randomticker

import (
	"math/rand"
	"time"
)

// DefaultMinimumDuration is used as the minimum duration for for Tick()
const DefaultMinimumDuration = time.Second

// A Ticker holds a channel that delivers `ticks' of a clock at random intervals within the specified range.
type Ticker struct {
	C      <-chan time.Time
	minDur time.Duration
	maxDur time.Duration
	dead   bool
}

// NewRandomTicker returns a new Ticker containing a channel that will send the current time after a random duration within the specified minimum and maximum durations.
func NewRandomTicker(minDur time.Duration, maxDur time.Duration) *Ticker {
	t := &Ticker{
		dead:   true,
		minDur: minDur,
		maxDur: maxDur,
	}
	return t
}

func (t *Ticker) runFunc(c chan time.Time) {
	rand.Seed(time.Now().UnixNano())
	for {
		dur := generateDuration(t.minDur, t.maxDur)
		time.Sleep(dur)
		if t.dead {
			close(c)
			return
		}
		c <- time.Now()
	}
}

func generateDuration(minDur time.Duration, maxDur time.Duration) time.Duration {
	r := rand.Float64()
	out := time.Duration(r*float64(maxDur-minDur)) + minDur
	return out
}

//Start starts the ticker. It also initializes the channel, so don't look try to use it before starting the ticker.
func (t *Ticker) Start() {
	t.dead = false
	c := make(chan time.Time, 1)
	t.C = c
	go t.runFunc(c)
}

//Tick is a convenience function that returns the underlying channel of the ticker configured with a predefined minimum duration. Just like time.Tick(), it can leak memory
func Tick(maxDur time.Duration) <-chan time.Time {
	ticker := NewRandomTicker(time.Duration(0), maxDur)
	ticker.Start()
	return ticker.C
}

//Stop stops the ticker
func (t *Ticker) Stop() {
	t.dead = true
}
