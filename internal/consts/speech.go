package consts

const SUBTITLE_FOLDER string = "internal/files/subtitle"
const VIDEOS_FOLDER string = "internal/files/videos"
const TEACHER_ROLE uint = 2
const STUDENT_ROLE uint = 3

func GetSubtitleFolder() string {
	return SUBTITLE_FOLDER
}
