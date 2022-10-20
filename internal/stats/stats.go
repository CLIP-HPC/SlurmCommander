package stats

import (
	"sort"
	"time"
)

func Median(s []time.Duration) (time.Duration, time.Duration, time.Duration) {
	var ret time.Duration

	n := len(s)
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
	for _, v := range s {
		ret += v
	}
	ret = ret / time.Duration(n)

	return ret
}
