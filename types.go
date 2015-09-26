package main

import "time"

type configRollOut struct {
	groups []string
	filter string
}

type result struct {
	instance   string
	launchTime time.Time
}

type results []result

func (r results) Len() int {
	return len(r)
}

func (r results) Less(i, j int) bool {
	return r[i].launchTime.Before(r[j].launchTime)
}

func (r results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r results) Values() []string {
	v := make([]string, len(r))
	for idx := range r {
		v[idx] = r[idx].instance
	}
	return v
}
