package entities

type User struct {
	ID       string `db:id json:id`
	Username string `db:username json:username`
	Password string `db:password json:password`
	Email    string `db:email json:email`
}

type RegisterUserDto struct {
	Username string `json:username`
	Password string `json:password`
	Email    string `json:email`
}

func (this *RegisterUserDto) IsValid() bool {
	if this.Username == "" {
		return false
	}
	if this.Password == "" {
		return false
	}
	if this.Email == "" {
		return false
	}

	return true
}

type LoginUserDto struct {
	Username string `json:username`
	Password string `json:password`
}

func (this *LoginUserDto) IsValid() bool {
	if this.Username == "" {
		return false
	}
	if this.Password == "" {
		return false
	}

	return true
}
