package dbase

import (
	"encoding/json"
	"time"
)

var cache = map[string]map[string]any{}

type Cache struct {
	Module      string
	Time        time.Time
	Interval    time.Duration
}

func (c *Cache) get(key any, get_value func()any, interval time.Duration, force bool) interface{} {
	if _, ok := cache[c.Module]; !ok {
		cache[c.Module] = map[string]any{}
	}
	lcache := cache[c.Module]
	key_tmp, _ := json.Marshal(key)
	lkey := string(key_tmp)

	if _, ok := lcache[lkey]; !ok || force  || time.Now().Sub(c.Time) > interval {
		val := get_value()
		lcache[lkey] = val
	}

	return lcache[lkey]
}

func (c *Cache) Get(key any, get_value func()any) interface{} {
	return c.get(key, get_value, c.Interval, false)
}

func (c *Cache) ForceGet(key any, get_value func()any) interface{} {
	return c.get(key, get_value, c.Interval, true)
}

func (c *Cache) GetWithInterval(key any, get_value func()any, interval time.Duration) interface{} {
	return c.get(key, get_value, interval, false)
}

func (c *Cache) ForceGetWithInterval(key any, get_value func()any, interval time.Duration) interface{} {
	return c.get(key, get_value, interval, true)
}
