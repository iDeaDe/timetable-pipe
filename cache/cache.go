package cache

import (
	"log"
	"os"
	"runtime"
	"time"
)

type Item struct {
	ExpireTime time.Time
	Data       interface{}
}

type Store struct {
	Items map[string]*Item
}

func cacheLogger() *log.Logger {
	return log.New(os.Stdout, "[Cache]", log.LstdFlags)
}

func (s *Store) Set(id string, item *Item) {
	if s.Items == nil {
		s.Items = make(map[string]*Item)
	}

	s.Items[id] = item
	cacheLogger().Printf("Added item \"%s\"\n", id)
}

func (s *Store) Remove(id string) {
	delete(s.Items, id)
	cacheLogger().Printf("Removed item \"%s\"\n", id)
}

func (s *Store) IsEmpty() bool {
	return len(s.Items) == 0
}

func (s *Store) Clear() {
	for id := range s.Items {
		s.Remove(id)
	}
	cacheLogger().Println("Store is empty")
}

func (s *Store) Refresh() {
	var collectGarbage = false
	cacheLogger().Println("Refreshing...")
	now := time.Now()
	for id, item := range s.Items {
		if item.ExpireTime.Before(now) {
			s.Remove(id)
			collectGarbage = true
		}
	}

	if collectGarbage {
		cacheLogger().Println("Collecting garbage")
		runtime.GC()
	}
}
