package timetable

import (
	"net/http"
	"net/url"
	"strconv"
)

const (
	scheme = "http"
	domain = "is.krmt.edu.ru"
	path   = "/blocks/manage_groups/website/"
)

type Handler struct {
	RequestURI url.URL
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ReqUrl() {
	var timetableUrl url.URL
	timetableUrl.Host = domain
	timetableUrl.Scheme = scheme
	timetableUrl.Path = path
	h.RequestURI = timetableUrl
}

func (h *Handler) GetGroups() (Courses, error) {
	reqUrl := h.RequestURI
	reqUrl.Path = reqUrl.Path + "list.php"
	qVal := reqUrl.Query()
	qVal.Add("id", "1")
	reqUrl.RawQuery = qVal.Encode()

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return Courses{}, err
	}

	return ParseGroups(resp.Body)
}

func (h *Handler) GetTimetable(group, week int) (Week, error) {
	reqUrl := h.RequestURI
	reqUrl.Path = reqUrl.Path + "view.php"
	qVal := reqUrl.Query()
	qVal.Add("dep", "1")
	qVal.Add("gr", strconv.Itoa(group))
	if week != 0 {
		qVal.Add("week", strconv.Itoa(week))
	}
	reqUrl.RawQuery = qVal.Encode()

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return Week{}, err
	}

	return ParseTimetable(resp.Body)
}
