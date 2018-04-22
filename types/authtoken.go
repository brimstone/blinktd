package types

type AuthToken struct {
	Pixels []int `json:"pixels"`
	Nbf    int   `json:"nbf"`
}
