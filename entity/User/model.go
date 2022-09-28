package User

type ModelUserList struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	MelliCode string `json:"melli_code"`
}


type RegisterForm struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type string `json:"type"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}