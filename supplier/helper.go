package supplier

// derefFloat64 safely dereferences a float64 pointer, returning 0 if the pointer is nil.
func derefFloat64(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0
}

// derefString safely dereferences a string pointer, returning an empty string if the pointer is nil.
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
