package middleware

// Custom errors
var (
	ErrInvalidToken  = &TokenError{message: "invalid token"}
	ErrMissingSecret = &TokenError{message: "missing jwt secret"}
)

type TokenError struct {
	message string
}

func (e *TokenError) Error() string {
	return e.message
}
