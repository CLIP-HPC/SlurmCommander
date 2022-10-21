package stats

import (
	"sort"
	"time"
)

// Return med,min,max
func Median(s []time.Duration) (time.Duration, time.Duration, time.Duration) {
	var ret time.Duration

	n := len(s)
	switch n {
	case 0:
		return 0, 0, 0
	case 1:
		return s[0], s[0], s[0]
	}

	n -= 1

	sort.Slice(s, func(i, j int) bool {
		if s[i] < s[j] {
			return true
		} else {
			return false
		}
	})

	if n%2 == 0 {
		ret = (s[n/2] + s[n/2+1]) / 2
	} else {
		ret = s[(n+1)/2]
	}
	return ret, s[0], s[n-1]
}

func Avg(s []time.Duration) time.Duration {
	var ret time.Duration

	n := len(s)
	if n == 0 {
		return 0
	}

	for _, v := range s {
		ret += v
	}
	ret = ret / time.Duration(n)

	return ret
}
