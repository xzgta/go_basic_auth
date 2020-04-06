package structs

type User struct {
	Id       string `form:"id_user" json:"id_user"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Token    string `form:"token" json:"token"`
}
