package models

type KindRole string

const (
	UnknownRole  KindRole = "unknown"
	CustomerRole KindRole = "customer"
	StaffRole    KindRole = "staff"
)

func (k KindRole) String() string {
	return string(k)
}

func NewKindRole(kind string) KindRole {
	switch kind {
	case CustomerRole.String():
		return CustomerRole
	case StaffRole.String():
		return StaffRole
	default:
		return UnknownRole
	}
}
