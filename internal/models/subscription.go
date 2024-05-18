package models

type SubscriptionType struct {
	slug string
}

var (
	UnknownSubscription = SubscriptionType{"unknown"}
	TrialSubscription   = SubscriptionType{"trial"}
)

func (k SubscriptionType) String() string {
	return k.slug
}

func NewSubscriptionType(kind string) SubscriptionType {
	switch kind {
	case TrialSubscription.String():
		return TrialSubscription
	default:
		return UnknownSubscription
	}
}

type SubscriptionDuration struct {
	slug string
}

var (
	UnknownSubscriptionDuration = SubscriptionDuration{"unknown"}
	OneWeekSubscriptionDuration = SubscriptionDuration{"1w"}
)

func (d SubscriptionDuration) String() string {
	return d.slug
}

func NewSubscriptionDuration(duration string) SubscriptionDuration {
	switch duration {
	default:
		return UnknownSubscriptionDuration
	}
}
