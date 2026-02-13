package utils

type EmailPayload struct {
	Addresses *[]string
	Subject   string
	Body      string
}

type EmailHelloVars struct {
	Name string
}

type ResetPasswordVars struct {
	Link string
}

type WelcomeEmailVars struct {
	Username       string
	LoginURL       string
	UnsubscribeURL string
	Password       string
	AssetsBaseURL  string
}
