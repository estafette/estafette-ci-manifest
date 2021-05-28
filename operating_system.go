package manifest

type OperatingSystem string

const (
	OperatingSystemUnknown OperatingSystem = ""
	OperatingSystemLinux   OperatingSystem = "linux"
	OperatingSystemWindows OperatingSystem = "windows"
)

// OperatingSystemArrayContains checks if an array contains a specific value
func OperatingSystemArrayContains(array []OperatingSystem, search OperatingSystem) bool {
	for _, v := range array {
		if v == search {
			return true
		}
	}
	return false
}
