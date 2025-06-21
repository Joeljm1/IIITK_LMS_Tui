package content

type CourseNameItem string

func (item CourseNameItem) Title() string {
	return string(item)
}
func (item CourseNameItem) Description() string { return "" }
func (item CourseNameItem) FilterValue() string {
	return ""
}
