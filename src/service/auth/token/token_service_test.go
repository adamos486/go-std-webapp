package token_test

import (
	"service/auth/token"
	"service/auth/token/tokenfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token Service Specs", func() {
	var (
		tokenService   *token.Service
		fakeTypedToken *tokenfakes.FakeInterface
	)

	BeforeEach(func() {
		fakeTypedToken = &tokenfakes.FakeInterface{}
		tokenService = token.NewService(fakeTypedToken)
	})

	Context("ValidateToken", func() {

		JustBeforeEach(func() {
			tokenService.ValidateToken("some test token")
		})

		It("should hit interface token handler only once", func() {
			Expect(fakeTypedToken.ValidateTokenCallCount()).To(Equal(1))
		})
	})

	Context("Generate", func() {
		JustBeforeEach(func() {
			m := make(map[string]interface{})
			m["email"] = "test@gmail.com"
			tokenService.Generate(m)
		})

		It("should hit interface token service only once", func() {
			Expect(fakeTypedToken.GenerateCallCount()).To(Equal(1))
		})
	})
})
