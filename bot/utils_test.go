package bot

import (
	"testing"
)

func TestRoute(t *testing.T) {
	route := []string{"a", "b"}
	topic := []string{"a", "b"}
	if !RouteIncludesTopic(route, topic) {
		t.Fail()
	}

	topic = []string{"a", "b", "c"}
	if RouteIncludesTopic(route, topic) {
		t.Fail()
	}

	route = []string{"a", "b", "c"}
	topic = []string{"a", "b"}
	if RouteIncludesTopic(route, topic) {
		t.Fail()
	}

	route = []string{"a", "b", "c"}
	topic = []string{"a", "b", "d"}
	if RouteIncludesTopic(route, topic) {
		t.Fail()
	}
}

func TestWildcard(t *testing.T) {
	route := []string{"#"}
	topic := []string{"anything", "really"}
	if !RouteIncludesTopic(route, topic) {
		t.Fail()
	}

	route = []string{"+", "level"}
	topic = []string{"something", "level"}
	if !RouteIncludesTopic(route, topic) {
		t.Fail()
	}
}
