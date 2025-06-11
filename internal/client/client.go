package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
)

const (
	loginUrl        = "https://lmsug23.iiitkottayam.ac.in/login/index.php"
	courseDetailURL = "https://lmsug23.iiitkottayam.ac.in/lib/ajax/service.php?sesskey=%v&info=core_course_get_enrolled_courses_by_timeline_classification"
	// change classification to in progress if needed in url
	allCourseDetailReqBody = `[{"index":0,"methodname":"core_course_get_enrolled_courses_by_timeline_classification","args":{"offset":0,"limit":0,"classification":"all","sort":"ul.timeaccess desc","customfieldname":"","customfieldvalue":""}}]`
	courseDetail           = `https://lmsug23.iiitkottayam.ac.in/course/view.php?id=%v`
)

var (
	ErrIncorrectCredentials = errors.New("incorrect credentials")
	ErrNoSessKey            = errors.New("cannot find sesskey")
	ErrLMSServerErr         = errors.New("LMS server error")
)

type LMSCLient struct {
	HttpClient     *http.Client
	Sesskey        string
	Choices        *Choices // making it pointer may error later
	RecivedChoices bool
}

// Returns a client with logged in cookies.
// ErrIncorrectCredentials is returned if username or password
// is incorrect.
func NewClient(username, password string) (*LMSCLient, error) { // need to ssee valid info
	cookieJar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List, // nil is fine as only one site but just in case
	})
	if err != nil { // will never happen (see the New fn code)
		return nil, err
	}
	client := http.Client{
		Jar: cookieJar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(loginUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status not 200 but %v", resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing from doc:%v\n", err.Error())
	}
	token, present := doc.Find(`input[name="logintoken"]`).First().Attr("value")
	if !present {
		return nil, errors.New("token not found")
	}
	values := url.Values{}
	values.Set("logintoken", token)
	values.Set("username", username)
	values.Set("password", password)
	resp, err = client.PostForm(loginUrl, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respURL := resp.Request.URL.String()
	if respURL == "https://lmsug23.iiitkottayam.ac.in/login/index.php" {
		return nil, ErrIncorrectCredentials
	}
	// Use regex to extract sesskey
	re := regexp.MustCompile(`"sesskey":"(.+?)"`)
	matches := re.FindSubmatch(html) // matches can be nil
	if matches == nil {
		return nil, ErrNoSessKey
	}
	lmsclient := &LMSCLient{
		HttpClient:     &client,
		Sesskey:        string(matches[1]),
		RecivedChoices: false,
	}
	return lmsclient, nil
}
