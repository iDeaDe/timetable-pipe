package timetable

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Course struct {
	Number int     `json:"course"`
	Groups []Group `json:"groups"`
}

type Courses struct {
	Count int      `json:"count"`
	Items []Course `json:"courses"`
}
