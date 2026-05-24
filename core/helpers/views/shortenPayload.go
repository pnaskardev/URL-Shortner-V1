package views

//	type SampleRequest struct {
//		Name string `json:"name" validate:"required,min=3"`
//		Age  int    `json:"age" validate:"required,min=1"`
//	}

type ShortenRequest struct {
	Url string `json:"url" validate:"url"`
}
