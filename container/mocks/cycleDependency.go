package mocks

//UserProvider depends on RoleProvider
type UserProvider struct {
	roleProvider RoleProvider
}

//NewUserProvider UserProvider -> RoleProvider
func NewUserProvider(roleProvider RoleProvider) UserProvider {
	return UserProvider{roleProvider: roleProvider}
}

//RoleProvider depends on RightsProvider
type RoleProvider struct {
	rightsProvider RightsProvider
}

//NewRoleProvider RoleProvider -> RightsProvider
func NewRoleProvider(rightsProvider RightsProvider) RoleProvider {
	return RoleProvider{rightsProvider: rightsProvider}
}

//RightsProvider depends on nothing
type RightsProvider struct {
}

//NewRightsProvider RightsProvider ->
func NewRightsProvider() RightsProvider {
	return RightsProvider{}
}
