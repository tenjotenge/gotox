package verify

// StubVerifier is a stub implementation of Verifier for testing and scaffolding
type StubVerifier struct{}

// NewStubVerifier creates a new stub verifier
func NewStubVerifier() *StubVerifier {
	return &StubVerifier{}
}

// VerifySchedule is a stub implementation that always returns valid
func (s *StubVerifier) VerifySchedule(data interface{}) (*ValidationResult, error) {
	return &ValidationResult{Valid: true}, nil
}

// VerifyCommit is a stub implementation that always returns valid
func (s *StubVerifier) VerifyCommit(data interface{}) (*ValidationResult, error) {
	return &ValidationResult{Valid: true}, nil
}

// VerifyDataIntegrity is a stub implementation that always returns valid
func (s *StubVerifier) VerifyDataIntegrity(data interface{}) (*ValidationResult, error) {
	return &ValidationResult{Valid: true}, nil
}
