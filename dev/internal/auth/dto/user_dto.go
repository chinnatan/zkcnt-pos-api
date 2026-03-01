package auth_dto

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	StoreID      string `json:"storeId"`
	StoreName    string `json:"storeName"`
	StoreAddress string `json:"storeAddress"`
}
