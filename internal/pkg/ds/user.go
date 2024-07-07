package ds

type (
	User struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		UpdateAt  string `json:"updated_at"`
		CreatedAt string `json:"created_at"`
	}

	UserToken struct {
		ID        string `json:"id"`
		UserID    string `json:"user_id"`
		Token     string `json:"token"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		IsActive  bool   `json:"is_active"`
	}
)
