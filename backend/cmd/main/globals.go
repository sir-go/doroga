package main

const DefaultConfFile = "config.toml"

func StringSliceContains(sl []string, item string) bool {
	for _, i := range sl {
		if i == item {
			return true
		}
	}
	return false
}
