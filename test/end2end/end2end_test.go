package end2end_test

import (
	. "github.com/onsi/ginkgo/v2"
	//. "github.com/onsi/gomega"
)

var _ = Describe("End2End", func() {
	Context("Testing item buying scenario", Serial, func() {
		var token string

		When("new user tries to register", Ordered, func() {
			It("app returns status 'OK' (200)", func() {

			})
			It("app returns a token", func() {

			})
			It("app increases user balance by 1000 coins", func() {

			})
			It("app has no user info", func() {

			})
		})

		When("registered user tries to buy an existed item and has enough coins", Ordered, func() {
			BeforeAll(func() {

			})

			It("app returns status 'OK' (200)", func() {

			})
			It("app decreases user balance by item price", func() {

			})
			It("app adds item to the user inventory", func() {

			})
		})

		When("registered user tries to buy an existed item and has not enough coins", Ordered, func() {
			BeforeAll(func() {

			})

			It("app returns status 'Bad request' (400)", func() {

			})
			It("app returns an error", func() {

			})
		})

		When("registered user tries to buy an unknown item", Ordered, func() {
			BeforeAll(func() {

			})

			It("app returns status 'Bad request' (400)", func() {

			})
			It("app returns an error", func() {

			})
		})

	})

	Context("Testing coins sending scenario", Serial, func() {

	})
})
