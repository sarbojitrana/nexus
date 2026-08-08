package service

import (
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
)

// primaryEmailFromClerkUser extracts the primary email address from a Clerk user,
// falling back to the first email on record if no primary is set.
func primaryEmailFromClerkUser(u *clerk.User) (string, error) {
	if len(u.EmailAddresses) == 0 {
		return "", fmt.Errorf("clerk user %s has no email addresses", u.ID)
	}

	if u.PrimaryEmailAddressID != nil {
		for _, email := range u.EmailAddresses {
			if email.ID == *u.PrimaryEmailAddressID {
				return email.EmailAddress, nil
			}
		}
	}

	return u.EmailAddresses[0].EmailAddress, nil
}

// generatePlaceholderUsername derives a unique placeholder username from a Clerk ID
// so the mini-user row can be created before the user picks a real one.
func generatePlaceholderUsername(clerkID string) string {
	suffix := clerkID
	if len(suffix) > 10 {
		suffix = suffix[len(suffix)-10:]
	}
	return "user_" + strings.ToLower(suffix)
}

// displayNameFromClerkUser builds a best-effort display name from Clerk's profile
// data, falling back to the email local-part and finally a generic default.
func displayNameFromClerkUser(u *clerk.User, email string) string {
	first := ""
	if u.FirstName != nil {
		first = strings.TrimSpace(*u.FirstName)
	}

	last := ""
	if u.LastName != nil {
		last = strings.TrimSpace(*u.LastName)
	}

	if name := strings.TrimSpace(first + " " + last); name != "" {
		return name
	}

	if at := strings.Index(email, "@"); at > 0 {
		return email[:at]
	}

	return "New User"
}

// clerkImageURL returns Clerk's hosted avatar, or nil when the account has no
// image -- HasImage is false for the auto-generated initials placeholder, which
// isn't worth mirroring.
func clerkImageURL(u *clerk.User) *string {
	if !u.HasImage || u.ImageURL == nil || *u.ImageURL == "" {
		return nil
	}
	return u.ImageURL
}
