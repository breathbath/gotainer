package mocks


type UserProvider struct {
	roleProvider RoleProvider
}

func NewUserProvider(roleProvider RoleProvider) UserProvider {
	return UserProvider {roleProvider:roleProvider}
}

type RoleProvider struct {
	rightsProvider RightsProvider
}

func NewRoleProvider(rightsProvider RightsProvider) RoleProvider {
	return RoleProvider{rightsProvider: rightsProvider}
}

type RightsProvider struct {
}

func NewRightsProvider() RightsProvider {
	return RightsProvider{}
}
