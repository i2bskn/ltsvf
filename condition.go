package main

type Condition struct {
	filters map[string]string
}

func newCondition(filters map[string]string) *Condition {
	return &Condition{
		filters: filters,
	}
}

func (condition *Condition) copiedFilters() map[string]string {
	filters := make(map[string]string)
	for key, value := range condition.filters {
		filters[key] = value
	}
	return filters
}
