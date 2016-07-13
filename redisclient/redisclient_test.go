// arkadiko
// https://github.com/topfreegames/arkadiko
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package redisclient

import (
	"testing"

	. "github.com/franela/goblin"
	"github.com/garyburd/redigo/redis"
	. "github.com/onsi/gomega"
)

func TestRedisClient(t *testing.T) {
	g := Goblin(t)

	// special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Redis Client", func() {
		g.It("It should set and get without error", func() {
			rc := GetRedisClient("localhost", 4444, "")
			_, err := rc.Pool.Get().Do("set", "teste", 1)
			Expect(err).To(BeNil())
			res, err := redis.Int(rc.Pool.Get().Do("get", "teste"))
			Expect(err).To(BeNil())
			g.Assert(res).Equal(1)
		})
	})
}
