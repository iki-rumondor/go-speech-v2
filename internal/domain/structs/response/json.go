package response

type Assignment struct {
	Uuid          string    `json:"uuid"`
	Submitted     bool      `json:"submitted"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Deadline      int64     `json:"deadline"`
	CreatedAt     int64     `json:"created_at"`
	UpdatedAt     int64     `json:"updated_at"`
	Class         *Class    `json:"class"`
	StudentAnswer *Answer   `json:"student_answer"`
	Answers       *[]Answer `json:"answers"`
}

type Answer struct {
	Uuid      string   `json:"uuid"`
	Ontime    bool     `json:"ontime"`
	Grade     int      `json:"grade"`
	Filename  string   `json:"filename"`
	Submitted bool     `json:"submitted"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
	Student   *Student `json:"student"`
}

type Notification struct {
	Uuid      string `json:"uuid"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Class     *Class `json:"class"`
}

type StudentInformation struct {
	UnreadNotification int `json:"unread_notification"`
}
