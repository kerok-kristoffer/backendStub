package util

const (
	USD = "USD"
	EUR = "EUR"
	SEK = "SEK"
)

// sample "enum?" to demonstrate custom validators for api calls seen in #14 in tut
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, SEK:
		return true
	}
	return false
}
