package entity

type (
	User struct {
		ID       string `json:"id" example:"1"`
		Email    string `json:"email" example:"example@xyz.xyz"`
		Username string `json:"username" example:"John Doe"`
		Password string `json:"password" example:"Password@123"`
		IsAdmin  bool   `json:"isAdmin" example:"false"`
	}

	AuthResponse struct {
		Token        string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		User         User   `json:"user"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	RegisterRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Username string `json:"username" validate:"required,min=3,max=32"`
	}
)
