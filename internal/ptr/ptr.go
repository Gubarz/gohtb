package ptr

// ClonePtr creates a new pointer to a copy of the value
// that the input pointer references
func Clone[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}
	clone := *ptr
	return &clone
}
