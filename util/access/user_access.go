package access

type UserAccess int

const (
	NONE = iota
	BETA_MONTHLY
	BETA_YEARLY
)

func (ua UserAccess) String() string {
	return [...]string{"none", "beta_monthly", "beta_yearly"}[ua]
}

func (ua UserAccess) EnumIndex() int {
	return int(ua)
}
