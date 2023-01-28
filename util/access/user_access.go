package access

type UserAccess int

const (
	NONE = iota
	BETA_MONTHLY
	BETA_YEARLY
	ALPHA_PERMANENT = 90
	ADMIN           = 99
)

func (ua UserAccess) String() string {
	return [...]string{"none", "beta_monthly", "beta_yearly", "alpha", "admin"}[ua]
}

func (ua UserAccess) EnumIndex() int {
	return int(ua)
}
