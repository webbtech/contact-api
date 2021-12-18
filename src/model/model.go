package model

type ContactRequest struct {
	Email     *string `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Message   *string `json:"message"`
	Phone     *string `json:"phone"`
	Type      *string `json:"type"`
}
