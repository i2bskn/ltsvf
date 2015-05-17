package main

type Condition struct {
	filters    map[string]string
	keys       []string
	concurrent bool
	channel    chan int
}

func newCondition(filters map[string]string, keys []string, concurrent bool) *Condition {
	return &Condition{
		filters:    filters,
		keys:       keys,
		concurrent: concurrent,
		channel:    make(chan int),
	}
}

func (condition *Condition) copiedFilters() map[string]string {
	filters := make(map[string]string)
	for key, value := range condition.filters {
		filters[key] = value
	}
	return filters
}

func (condition *Condition) displayKey(target string) bool {
	if len(condition.keys) > 0 {
		for _, key := range condition.keys {
			if key == target {
				return true
			}
		}
		return false
	} else {
		return true
	}
}
