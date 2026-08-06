package email

func (c *Client) SendWelcomeEmail(to, firstName string) error {
	data := map[string]string{
		"UserFirstName": firstName,
	}

	return c.SendEmail(
		to,
		"Welcome to Nexus !!!",
		TemplateWecome,
		data,
	)
}

func (c *Client) SendSignInEmail(to, firstName, signInTime string) error {
	data := map[string]string{
		"UserFirstName": firstName,
		"SignInTime":    signInTime,
	}

	return c.SendEmail(
		to,
		"New sign-in to your Nexus account",
		TemplateSignIn,
		data,
	)
}

func (c *Client) SendPasswordChangedEmail(to, firstName, changedAt string) error {
	data := map[string]string{
		"UserFirstName": firstName,
		"ChangedAt":     changedAt,
	}

	return c.SendEmail(
		to,
		"Your Nexus password was changed",
		TemplatePasswordChanged,
		data,
	)
}
