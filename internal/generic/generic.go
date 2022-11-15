package generic

import "sort"

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
