package timetable

type Lesson struct {
	Name        string `json:"name"`
	TeacherName string `json:"teacher"`
	Cabinet     string `json:"cabinet"`
	Canceled    bool   `json:"canceled"`
}

type LessonInfo struct {
	Number int      `json:"number"`
	Start  string   `json:"start"`
	End    string   `json:"end"`
	Lesson []Lesson `json:"info"`
}

type Day struct {
	Name    string       `json:"day"`
	Lessons []LessonInfo `json:"lessons"`
}

type Week struct {
	WeekNumber int    `json:"week"`
	Days       [6]Day `json:"days"`
}
