package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(c string) bool {
	switch c {
	case USD, EUR, CAD:
		return true
	}

	return false
}
