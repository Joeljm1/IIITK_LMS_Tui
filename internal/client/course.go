package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type (
	Course struct {
		Id             int    `json:"id"`
		Fullname       string `json:"fullname"`
		Shortname      string `json:"shortname"`
		Coursecategory string `json:"coursecategory"`
	}
	CourseDetails struct {
		Error bool `json:"error"`
		Data  struct {
			Courses []Course `json:"courses"`
		} `json:"data"`
	}
	CourseList []Course
)

// Need to close the Response.Not done in this fn
func (lmsClient LMSCLient) PostData(url string, headers *http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		req.Header = *headers
	}
	resp, err := lmsClient.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// also writes attendance list to file
func (lmsClient *LMSCLient) GetCoursesFromLMS() (CourseList, error) {
	body := strings.NewReader(allCourseDetailReqBody)
	postUrl := fmt.Sprintf(courseDetailURL, lmsClient.Sesskey)
	headers := &http.Header{}
	headers.Set("Host", "lmsug23.iiitkottayam.ac.in")
	headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0")
	headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	headers.Set("Accept-Language", "en-US,en;q=0.5")
	headers.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	headers.Set("Content-Type", "application/json")
	headers.Set("X-Requested-With", "XMLHttpRequest")
	headers.Set("Origin", "https://lmsug23.iiitkottayam.ac.in")
	headers.Set("Sec-GPC", "1")
	headers.Set("Connection", "keep-alive")
	headers.Set("Referer", "https://lmsug23.iiitkottayam.ac.in/my/courses.php")
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Site", "same-origin")
	headers.Set("Pragma", "no-cache")
	headers.Set("Cache-Control", "no-cache")
	resp, err := lmsClient.PostData(postUrl, headers, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	courseDetails := []CourseDetails{}
	err = json.NewDecoder(resp.Body).Decode(&courseDetails)
	if err != nil {
		return nil, err
	}
	if len(courseDetails) == 0 {
		return nil, errors.New("no courses present")
	}
	if courseDetails[0].Error {
		return nil, ErrLMSServerErr
	}
	courses := courseDetails[0].Data.Courses
	/*for _, course := range courseDetails[0].Data.Courses {
		course := Course{
			Id:             course.Id,
			Fullname:       course.Fullname,
			Shortname:      course.Shortname,
			Coursecategory: course.Coursecategory,
		}
		courses = append(courses, course)
	}*/
	courseList := CourseList(courses)
	return courseList, nil
}
