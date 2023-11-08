package yauthorization

type GoogleClaim struct {
	Aud           string `json:"aud"`
	Azp           string `json:"azp"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Exp           int    `json:"exp"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Iat           int    `json:"iat"`
	Iss           string `json:"iss"`
	Jti           string `json:"jti"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Nbf           int    `json:"nbf"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
}
