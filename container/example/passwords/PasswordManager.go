package passwords

type PasswordManager struct {
}

//for each service we suggest to create a constructor function
func NewPasswordManager() PasswordManager {
	return PasswordManager{}
}

func (pm PasswordManager) GetLoginData() (login, pass string, err error) {
	//of course in real life here should be a more complex logic to extract sensitive password data
	return "no@mail.me", "1234567", nil
}

