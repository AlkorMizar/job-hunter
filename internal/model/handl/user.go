package handl

type NewUser struct {
	Login    string   `json:"login" validate:"required,min=3,max=40"`
	Email    string   `json:"email" validate:"required,email"`
	Roles    []string `json:"roles" validate:"required"`
	Password string   `json:"password" validate:"required,min=5,max=40,excludesall= "`
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

type User struct {
	Login    string   `json:"login"    validate:"required,min=3,max=40"`
	Email    string   `json:"email"    validate:"required,email"`
	Roles    []string `json:"roles"    validate:"required"`
	FullName string   `json:"fullName" validate:"required,min=3,max=150"`
}

type UpdateInfo struct {
	Login    string `json:"login"    validate:"required_without_all=Email FullName,omitempty,min=3,max=40,alphanum"`
	Email    string `json:"email"    validate:"required_without_all=Login FullName,omitempty,email"`
	FullName string `json:"fullName" validate:"required_without_all=Login Email,omitempty,min=3,max=150"`
}

type Passwords struct {
	NewPassword  string `json:"newPassword"  validate:"required,min=5,max=40,excludesall= "`
	CurrPassword string `json:"curPassword" validate:"required,min=5,max=40,excludesall= "`
}
