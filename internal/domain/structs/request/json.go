package request

type Assignment struct {
	ClassUuid   string `json:"class_uuid" valid:"required~field uuid kelas tidak ditemukan"`
	Title       string `json:"title" valid:"required~field judul tidak ditemukan"`
	Description string `json:"description" valid:"required~field description tidak ditemukan"`
	Deadline    int64  `json:"deadline" valid:"required~field deadline tidak ditemukan"`
}

type UpdateAssignment struct {
	Title       string `json:"title" valid:"required~field judul tidak ditemukan"`
	Description string `json:"description" valid:"required~field description tidak ditemukan"`
}

type Grading struct {
	Grade int `json:"grade" valid:"required~field nilai tidak ditemukan, int~tipe data tidak valid, range(1|100)~range nilai dari 1 - 100"`
}

type ReadNotification struct {
	NotificationUuid string `json:"notification_uuid"`
}
