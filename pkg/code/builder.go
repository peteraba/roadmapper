package code

// Builder is a factory for codes
type Builder struct {
}

// NewFromString creates a new code from a string representation (typically a URL path)
func (b Builder) NewFromString(s string) (Code, error) {
	return newCode64FromString(s)
}

// NewFromID creates a new code from a number
func (b Builder) NewFromID(id uint64) (Code, error) {
	return Code64(id), nil
}

// New creates a new random code
func (b Builder) New() Code {
	return NewCode64()
}

// NewCodeBuilder creates a new Builder instance
func NewCodeBuilder() Builder {
	return Builder{}
}
