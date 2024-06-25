package request

type Class struct {
	Name string `json:"name" valid:"required~field nama tidak ditemukan"`
	Code string `json:"code" valid:"required~field kode tidak ditemukan, stringlength(8|8)~field kode harus 8 karakter"`
}

type Department struct {
	Name string `json:"name" valid:"required~field nama tidak ditemukan"`
}

type Note struct {
	ClassUuid string `json:"class_uuid" valid:"required~field uuid class tidak ditemukan"`
	Title     string `json:"title" valid:"required~field judul tidak ditemukan"`
	Body      string `json:"body" valid:"required~field body tidak ditemukan"`
}

type UpdateVideo struct {
	Title       string `json:"title" valid:"required~field judul tidak ditemukan"`
	Description string `json:"description" valid:"required~field description tidak ditemukan"`
}
