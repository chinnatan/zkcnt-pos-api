package auth_dto

type SignUpDTO struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Role            string `json:"role"`
	StoreName       string `json:"storeName"`
	StoreAddress    string `json:"storeAddress"`
}
