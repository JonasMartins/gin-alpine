package users

type PaginatedUsersResult struct {
	Users             []*User
	TotalItems        uint32
	TotalPages        uint32
	TotalItemsPerPage uint32
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	RoleID   *int32
}
type CreateUserResult struct {
	ID int
}

func NewCreateUserInput() *CreateUserInput {
	return &CreateUserInput{}
}

type UpdateUserInput struct {
	ID       int
	Name     *string
	Email    *string
	Password *string
}

type UpdateUserAdminInput struct {
	ID          int
	Name        *string
	Email       *string
	Password    *string
	RoleID      *int
	ToBeDeleted *bool
	Enabled     *bool
}
