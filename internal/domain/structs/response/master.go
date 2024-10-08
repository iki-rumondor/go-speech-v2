package response

type Student struct {
	Uuid              string `json:"uuid"`
	Name              string `json:"name"`
	Nim               string `json:"nim"`
	Email             string `json:"email"`
	RegTimeString     string `json:"reg_time_string"`
	RegisterClassTime int64  `json:"register_class_time"`
}

type Class struct {
	Uuid              string     `json:"uuid"`
	Name              string     `json:"name"`
	Code              string     `json:"code" `
	Teacher           string     `json:"teacher"`
	TeacherUuid       string     `json:"teacher_uuid"`
	TeacherDepartment string     `json:"teacher_department" `
	Students          *[]Student `json:"students"`
}

type StudentClass struct {
	Join  bool   `json:"join"`
	Class *Class `json:"class"`
}

type Department struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type Video struct {
	Uuid         string `json:"uuid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	VideoName    string `json:"video_name"`
	SubtitleName string `json:"subtitle_name"`
	Teacher      string `json:"teacher"`
	CreatedAt    int64  `json:"created_at"`
}

type Book struct {
	Uuid        string    `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileName    string    `json:"file_name"`
	Teacher     string    `json:"teacher"`
	URLDropbox  string    `json:"url_dropbox"`
	CreatedAt   int64     `json:"created_at"`
	FlipBook    *FlipBook `json:"flipbook"`
}

type FlipBook struct {
	Url       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

type Note struct {
	Uuid      string `json:"uuid"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"created_at"`
}

type Material struct {
	Uuid         string `json:"uuid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	BookName     string `json:"book_name"`
	VideoUuid    string `json:"video_uuid"`
	VideoName    string `json:"video_name"`
	SubtitleName string `json:"subtitle_name"`
	CreatedAt    int64  `json:"created_at"`
	UpdateAt     int64  `json:"updated_at"`
	Class        *Class `json:"class"`
}
