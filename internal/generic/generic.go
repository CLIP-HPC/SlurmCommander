package generic

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type CountItemSlice []CountItem

type CountItem struct {
	Name  string
	Count uint
	Total uint
}

//type CountItemMap map[string]uint
type CountItemMap map[string]*CountItem

func SortItemMapBySel(what string, m *CountItemMap) CountItemSlice {
	var ret = CountItemSlice{}
	//ret := make(CountItemSlice, len(*m))
	for k, v := range *m {
		ret = append(ret, CountItem{
			Name:  k,
			Count: v.Count,
			Total: v.Total,
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		switch what {
		case "Count":
			if ret[i].Count > ret[j].Count {
				return true
			}
		case "Name":
			if ret[i].Name < ret[j].Name {
				return true
			}
		}
		return false
	})

	return ret
}

func Top5(src CountItemSlice) CountItemSlice {
	var ret CountItemSlice
	for i, v := range src {
		if i < 5 {
			ret = append(ret, v)
		}
	}
	return ret
}

func HumanizeDuration(t time.Duration, l *log.Logger) string {
	var ret string

	// total seconds
	s := int64(t.Seconds())

	// days
	d := s / (24 * 60 * 60)
	s = s % (24 * 60 * 60)

	// hours
	h := s / 3600
	s = s % 3600

	// minutes
	m := s / 60
	s = s % 60

	ret += fmt.Sprintf("%.2d-%.2d:%.2d:%.2d", d, h, m, s)

	l.Printf("Humanized %f to %q\n", t.Seconds(), ret)
	return ret
}

// Generate statistics string, vertical.
func GenCountStrVert(cnt map[string]uint, l *log.Logger) string {
	var (
		scr string
	)

	sm := make([]struct {
		name string
		val  uint
	}, 0)

	// place map to slice
	for k, v := range cnt {
		sm = append(sm, struct {
			name string
			val  uint
		}{name: k, val: uint(v)})
	}

	// sort first by name
	sort.Slice(sm, func(i, j int) bool {
		if sm[i].name < sm[j].name {
			return true
		} else {
			return false
		}
	})
	// then sort by numbers
	sort.Slice(sm, func(i, j int) bool {
		if sm[i].val > sm[j].val {
			return true
		} else {
			return false
		}
	})

	// print it out
	//scr = "Count: "
	for _, v := range sm {
		scr += fmt.Sprintf("%-15s: %d\n", v.name, v.val)
	}
	scr += "\n\n"

	return scr
}

// Generate statistics string, horizontal.
func GenCountStr(cnt map[string]uint, l *log.Logger) string {
	var (
		scr string
	)

	sm := make([]struct {
		name string
		val  uint
	}, 0)

	// place map to slice
	for k, v := range cnt {
		sm = append(sm, struct {
			name string
			val  uint
		}{name: k, val: uint(v)})
	}

	// sort it
	sort.Slice(sm, func(i, j int) bool {
		if sm[i].name < sm[j].name {
			return true
		} else {
			return false
		}
	})

	// print it out
	scr = "Count: "
	for _, v := range sm {
		scr += fmt.Sprintf("%s: %d ", v.name, v.val)
	}
	scr += "\n\n"

	return scr
}
