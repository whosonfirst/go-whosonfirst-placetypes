package placetypes

// COMMON_ROLE defines the string label for the "common" placetype role.
const COMMON_ROLE string = "common"

// OPTIONAL_ROLE defines the string label for the "optional" placetype role.
const OPTIONAL_ROLE string = "optional"

// COMMON_OPTIONAL_ROLE defines the string label for the "common_optional" placetype role.
const COMMON_OPTIONAL_ROLE string = "common_optional"

// CUSTOM_ROLE defines the string label for the "custom" placetype role.
const CUSTOM_ROLE string = "custom"

// AllRoles returns a list of all the known placetype roles.
func AllRoles() []string {

	roles := []string{
		COMMON_ROLE,
		OPTIONAL_ROLE,
		COMMON_OPTIONAL_ROLE,
		CUSTOM_ROLE,
	}

	return roles
}
