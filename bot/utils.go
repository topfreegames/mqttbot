package bot

// RouteIncludesTopic returns if the topic matches the given route
func RouteIncludesTopic(route []string, topic []string) bool {
	if len(route) == 0 {
		if len(topic) == 0 {
			return true
		}
		return false
	}

	if route[0] == "#" {
		return true
	}

	if len(topic) == 0 {
		return false
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return RouteIncludesTopic(route[1:], topic[1:])
	}

	return false
}
