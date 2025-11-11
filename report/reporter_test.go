package report

import (
	"testing"
	"time"
)

func TestToHuman(t *testing.T) {
	cases := map[time.Duration]string{
		1500 * time.Microsecond: "1.50ms",
		2 * time.Second:         "2.000s",
		500 * time.Nanosecond:   "0.50µs",
	}

	for input, expected := range cases {
		if got := toHuman(input); got != expected {
			t.Fatalf("toHuman(%v) = %s, expected %s", input, got, expected)
		}
	}
}
