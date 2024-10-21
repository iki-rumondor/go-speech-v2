package request

type Class struct {
	TeacherUuid string `json:"teacher_uuid" valid:"required~field dosen tidak ditemukan"`
	Name        string `json:"name" valid:"required~field nama tidak ditemukan"`
	Code        string `json:"code" valid:"required~field kode tidak ditemukan, stringlength(8|8)~field kode harus 8 karakter"`
}

type Department struct {
	Name string `json:"name" valid:"required~field nama tidak ditemukan"`
}

type Note struct {
	ClassUuid string `json:"class_uuid" valid:"required~field uuid class tidak ditemukan"`
	Title     string `json:"title" valid:"required~field judul tidak ditemukan"`
	Body      string `json:"body" valid:"required~field body tidak ditemukan"`
}

type Material struct {
	ClassUuid   string `form:"class_uuid" valid:"required~field uuid class tidak ditemukan"`
	Title       string `form:"title" valid:"required~field judul tidak ditemukan"`
	Description string `form:"description" valid:"required~field description tidak ditemukan"`
}

type UpdateMaterial struct {
	Title       string `json:"title" valid:"required~field judul tidak ditemukan"`
	Description string `json:"description" valid:"required~field description tidak ditemukan"`
}

type UpdateVideo struct {
	Title       string `json:"title" valid:"required~field judul tidak ditemukan"`
	Description string `json:"description" valid:"required~field description tidak ditemukan"`
}

type CreateTeacher struct {
	Name           string `json:"name" valid:"required~field nama tidak ditemukan"`
	Nidn           string `json:"nidn" valid:"required~field nidn tidak ditemukan"`
	DepartmentUuid string `json:"department_uuid" valid:"required~field description tidak ditemukan"`
}

type UpdateTeacher struct {
	Name           string `json:"name" valid:"required~field nama tidak ditemukan"`
	Nidn           string `json:"nidn" valid:"required~field nidn tidak ditemukan"`
	DepartmentUuid string `json:"department_uuid" valid:"required~field description tidak ditemukan"`
}

type CreateStudent struct {
	Name string `json:"name" valid:"required~field nama tidak ditemukan"`
	Nim  string `json:"nim" valid:"required~field nim tidak ditemukan"`
}

type UpdateStudent struct {
	Name string `json:"name" valid:"required~field nama tidak ditemukan"`
	Nim  string `json:"nim" valid:"required~field nim tidak ditemukan"`
}
