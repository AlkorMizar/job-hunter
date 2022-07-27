package handl

type NewUser struct {
	Login    string   `json:"login" validate:"required,min=3,max=40"`
	Email    string   `json:"email" validate:"required,email"`
	Roles    []string `json:"roles" validate:"required"`
	Password string   `json:"password" validate:"required,min=5,max=40"`
}

type AuthInfo struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=40"`
}

type JSONResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type UserInfo struct {
	ID    int                 `json:"userId"`
	Roles map[string]struct{} `json:"roles"`
}
