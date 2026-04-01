package auth

type Config struct {
	Secret            string
	InternalAuthToken string
	InternalLogToken  string
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
