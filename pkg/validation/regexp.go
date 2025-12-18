package validation

const (
	localeRegexp               = "^[A-Za-z]{2}" //"^[A-Za-z]{2,4}([_-][A-Za-z]{4})?([_-]([A-Za-z]{2}|[0-9]{3}))?$"
	slugRegexp                 = "^[a-z0-9]+(?:-[a-z0-9]+)*$"
	genderRegexp               = "^male$|^female$|^none$"
	phoneWithCountryCodeRegexp = "^\\+[0-9]{1,3}[0-9]{3}[0-9]{3}[0-9]{4}$" // +1-123-456-7890, +90-123-456-7890, +123-123-456-7890
)
