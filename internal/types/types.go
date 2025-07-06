package types

type Student struct {
	ID        int64 `json:"id"`  
	Name      string `validate:"required" json:"name"`
	Age       int    `validate:"required" json:"age"`
	Email     string `validate:"required" json:"email"`
	City	  string `validate:"required" json:"city"`
}