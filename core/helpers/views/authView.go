package views

//	type SampleRequest struct {
//		Name string `json:"name" validate:"required,min=3"`
//		Age  int    `json:"age" validate:"required,min=1"`
//	}

type AuthSignInPayload struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=4"`
}
