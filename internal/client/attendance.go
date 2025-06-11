package client

import (
	"fmt"
	"io"
	"regexp"
	"sync"

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
