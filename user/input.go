package user

// InputUser represents json body for users POST/PUT api
type InputUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required"`
	Password  string `json:"password" validate:"required,ascii,min=8,max=72"`
	Email     string `json:"email" validate:"required,email"`
	Country   string `json:"country" validate:"required,iso3166_1_alpha2"`
}
