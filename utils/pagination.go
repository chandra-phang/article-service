package utils

const (
	DefaultPageLimit = 20
	MaxPageLimit     = 100
)

func SetLimit(limit int) int {
	if limit <= 0 {
		return DefaultPageLimit
	} else if limit > MaxPageLimit {
		return MaxPageLimit
	} else {
		return limit
	}
}

func SetOffset(page int, limit int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * limit
}
