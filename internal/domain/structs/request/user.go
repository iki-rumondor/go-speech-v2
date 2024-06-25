package request

type User struct {
	Name           string `json:"name" valid:"required~field nama tidak ditemukan"`
	Email          string `json:"email" valid:"required~field email tidak ditemukan, email~Gunakan email yang valid"`
	Password       string `json:"password" valid:"required~field password tidak ditemukan"`
	RoleID         string `json:"role_id" valid:"required~field role tidak ditemukan"`
	Nip            string `json:"nip"`
	Nim            string `json:"nim"`
	DepartmentUuid string `json:"department_uuid"`
}

type ClassRequest struct {
	ClassUuid string `json:"class_uuid" valid:"required~field kelas tidak ditemukan"`
}

type StatusClassReq struct {
	Status uint `json:"status" valid:"required~field status tidak ditemukan"`
}

type SignIn struct {
	Email    string `json:"email" valid:"required~field email tidak ditemukan, email~Gunakan email yang valid"`
	Password string `json:"password" valid:"required~field password tidak ditemukan"`
}
