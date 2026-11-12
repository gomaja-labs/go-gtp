package ie

// NewPGWChangeInfo creates a grouped PGW Change Info IE (TS 29.274 §8.145).
func NewPGWChangeInfo(ies ...*IE) *IE { return NewGroupedIE(PGWChangeInfo, ies...) }
