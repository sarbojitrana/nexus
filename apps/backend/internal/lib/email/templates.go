package email

type Template string

const (
	TemplateWecome          Template = "welcome"
	TemplateSignIn          Template = "signin"
	TemplatePasswordChanged Template = "password_changed"
)
