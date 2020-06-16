package timetable

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"strconv"
	"strings"
)

func ParseGroups(content io.Reader) (Courses, error) {
	var (
		err     error
		courses Courses
	)

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return Courses{}, err
	}

	doc.Find(".content > div:nth-child(3) > .spec-year-block-container > .spec-year-block").Each(func(_ int, course *goquery.Selection) {
		var currCourse Course
		currCourse.Number, err = strconv.Atoi(course.Find(".spec-year-name > span > b").Text())

		course.Find(".group-block").Each(func(i int, group *goquery.Selection) {
			var (
				gr     Group
				exists bool
			)

			idStr, exists := group.Attr("group_id")

			if exists {
				gr.Id, err = strconv.Atoi(idStr)
				if err != nil {
					log.Fatal(err)
				}

				gr.Name = group.Find("span[group_id] > a").Text()
			} else {
				gr.Id = 0
				gr.Name = "null"
			}

			currCourse.Groups = append(currCourse.Groups, gr)
		})

		courses.Items = append(courses.Items, currCourse)
	})

	courses.Count = len(courses.Items)

	return courses, nil
}

func ParseTimetable(content io.Reader) (Week, error) {
	var (
		err  error
		week Week
	)

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return Week{}, err
	}

	week.WeekNumber, err = strconv.Atoi(strings.Split(doc.Find(".weekHeader > span").Text(), " ")[0])
	if err != nil {
		return Week{}, err
	}

	table := doc.Find(".timetable > tbody > tr > td")

	table.Each(func(iDay int, day *goquery.Selection) {
		week.Days[iDay].Name = day.Find(".dayHeader > span").Text()

		day.Find(".lessonBlock").Each(func(_ int, lessonInfo *goquery.Selection) {
			var lInfo LessonInfo

			timeInfo := lessonInfo.Find(".lessonTimeBlock").Children().Nodes

			lInfo.Start = timeInfo[1].FirstChild.Data
			lInfo.End = timeInfo[2].FirstChild.Data
			lInfo.Number, err = strconv.Atoi(timeInfo[0].FirstChild.Data)

			lessonInfo.Find(".discBlock").Each(func(_ int, lesson *goquery.Selection) {
				var llesson Lesson

				if len(lesson.Children().Nodes) > 0 {
					llesson.Name = lesson.Find("span[title]").Text()
					llesson.Cabinet = lesson.Find(".discSubgroupClassroom > span").Text()
					llesson.TeacherName = lesson.Find(".discSubgroupTeacher > span").Text()
					llesson.Canceled = lesson.HasClass("cancelled")
				}

				lInfo.Lesson = append(lInfo.Lesson, llesson)
			})

			week.Days[iDay].Lessons = append(week.Days[iDay].Lessons, lInfo)
		})
	})

	return week, nil
}
