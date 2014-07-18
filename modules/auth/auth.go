package auth

// Web form interface.
type Form interface {
	Name(field string) string
}
