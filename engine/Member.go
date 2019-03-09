package engine

type Member struct {
	Id         uint64 `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Patronymic string `json:"patronymic"`
	Status     uint   `json:"status"`
}
