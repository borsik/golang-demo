package user

type InputUser struct {
	FirstName string `json:"first_name" binding:"-"`
	LastName  string `json:"last_name" binding:"-"`
	Nickname  string `json:"nickname" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Country   string `json:"country" binding:"-"`
}
