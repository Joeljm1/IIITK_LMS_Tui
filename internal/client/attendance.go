package client

import (
	"fmt"
	"io"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

const (
	attendanceURL = "https://lmsug23.iiitkottayam.ac.in/mod/attendance/view.php?id=%v"
	courseURL     = `https://lmsug23.iiitkottayam.ac.in/course/view.php?id=%v`
)

func (lmsClient *LMSCLient) GetAttendanceId(id int) (string, error) {
	reqURL := fmt.Sprintf(courseURL, id)
	resp, err := lmsClient.HttpClient.Get(reqURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	regExpr := regexp.MustCompile(`https://lmsug23\.iiitkottayam\.ac\.in/mod/attendance/view\.php\?id=(\d+)`)
	submatch := regExpr.FindSubmatch(b)
	if submatch == nil {
		return "", nil
	}
	return string(submatch[1]), nil
}

// Checks attendace of user selection is LMSClient.Choices.AllChoices
// and stores at map in LMSClient.Choices.AttendanceId
// Value not present if no attendance Id
func (lmsClient *LMSCLient) AttendanceForAllSelection(courseList CourseList) (*Choices, error) {
	aMap := make(map[string]string)
	var m sync.Mutex
	var g errgroup.Group
	for _, course := range courseList {
		g.Go(func() error {
			attId, err := lmsClient.GetAttendanceId(course.Id)
			if err != nil {
				return err
			}
			if attId != "" {
				m.Lock()
				aMap[course.Fullname] = attId
				m.Unlock()
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	choices := &Choices{
		AllChoices:   courseList,
		AttendanceId: aMap,
	}
	return choices, nil
}

type AttendanceStatus int

type Attendance struct {
	Date   string // prolly parse change it to a custom date type later
	Time   string
	Desc   string
	Status string
}

type OverallAttendance struct {
	Total      string
	Points     string
	Percentage string
}

const attendanceDetailsURL = "https://lmsug23.iiitkottayam.ac.in/mod/attendance/view.php?id=%v&view=5"

func (lmsCLient *LMSCLient) GetAttendance(id string) ([]Attendance, error) {
	URLWithId := fmt.Sprintf(attendanceDetailsURL, id)
	resp, err := lmsCLient.HttpClient.Get(URLWithId)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	ParsedTable := []Attendance{}
	doc.Find(".generaltable.attwidth.boxaligncenter > tbody > tr").Each(func(i int, s *goquery.Selection) {
		attRow := Attendance{}
		s.Find(".c0>nobr").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				attRow.Date = s.Text()
			case 1:
				attRow.Time = s.Text()
			}
		})
		attRow.Desc = s.Find(".c1>div").First().Text()
		attRow.Status = s.Find(".c2").First().Text()
		ParsedTable = append(ParsedTable, attRow)
	})
	return ParsedTable, nil
}

// Returns every sing attence of all courses
func (lmsClient *LMSCLient) AllAttendance() ([][]Attendance, error) {
	var allAttendance [][]Attendance
	var mut sync.Mutex
	var g errgroup.Group
	for _, id := range lmsClient.Choices.AttendanceId {
		g.Go(func() error {
			attend, err := lmsClient.GetAttendance(id)
			if err != nil {
				return err
			}
			mut.Lock()
			allAttendance = append(allAttendance, attend)
			return nil
		})
	}
	err := g.Wait()
	return allAttendance, err
}
