package main

import (
	"encoding/json"
	"github.com/ideade/timetable-pipe/cache"
	"github.com/ideade/timetable-pipe/timetable"
	"log"
	"net/http"
	"strconv"
	"time"
)

type GroupsHandler struct {
	cache     *cache.Store
	ttHandler timetable.Handler
}

func (s *GroupsHandler) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	s.cache.Refresh()

	if s.cache.IsEmpty() {
		week, err := s.ttHandler.GetGroups()
		if err != nil {
			response.WriteHeader(500)
		}

		s.cache.Set("groups", &cache.Item{
			ExpireTime: time.Now().Add(time.Hour * 72),
			Data:       week,
		})
	}

	jsonResponse, err := json.Marshal(s.cache.Items["groups"].Data)
	if err != nil {
		response.WriteHeader(500)
		log.Fatal(err)
	}

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	if _, err = response.Write(jsonResponse); err != nil {
		return
	}
}

type TimetableHandler struct {
	cache     *cache.Store
	ttHandler timetable.Handler
}

func (s *TimetableHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	group, uWeek := 0, 0

	// Проверяем необходимые параметры запроса
	if gr, err := strconv.Atoi(request.URL.Query().Get("group")); err == nil {
		group = gr
	} else {
		response.WriteHeader(400)
		return
	}

	if request.URL.Query().Get("week") != "" {
		if w, err := strconv.Atoi(request.URL.Query().Get("week")); err == nil {
			uWeek = w
		} else {
			response.WriteHeader(400)
			return
		}
	}
	// Получили оба параметра

	s.cache.Refresh()

	ciId := strconv.Itoa(group) + "-" + strconv.Itoa(uWeek)

	if s.cache.IsEmpty() || s.cache.Items[ciId] == nil {
		groupTimetable, err := s.ttHandler.GetTimetable(group, uWeek)
		if err != nil {
			response.WriteHeader(500)
		}

		s.cache.Set(ciId, &cache.Item{
			ExpireTime: time.Now().Add(time.Minute * 3),
			Data:       groupTimetable,
		})
	}

	jsonResponse, err := json.Marshal(s.cache.Items[ciId].Data)
	if err != nil {
		response.WriteHeader(500)
		log.Fatal(err)
	}

	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	if _, err = response.Write(jsonResponse); err != nil {
		return
	}
}
