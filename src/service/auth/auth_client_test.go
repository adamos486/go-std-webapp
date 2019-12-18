package auth_test

import (
	"net/http"
	"net/http/httptest"
	"service/auth"
	"service/auth/authfakes"
	"service/auth/token/tokenfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Client Specs", func() {
	var (
		authClient *auth.Client
		fakeAuth   *authfakes.FakeInterface
		fakeToken  *tokenfakes.FakeInterface
		request    *http.Request
	)

	BeforeEach(func() {
		fakeToken = &tokenfakes.FakeInterface{}
		fakeAuth = &authfakes.FakeInterface{}
		authClient = auth.NewClient(fakeAuth, fakeToken)
	})

	Context("Authorize", func() {
		BeforeEach(func() {
			request = httptest.NewRequest("POST", "/identity", nil)
		})

		JustBeforeEach(func() {
			authClient.Authorize(request)
		})

		It("should pass authorize to an auth interface", func() {
			Expect(fakeAuth.AuthorizeCallCount()).To(Equal(1))
		})
	})

	Context("ValidateTokenHeader", func() {
		BeforeEach(func() {
			request = httptest.NewRequest("POST", "/identity", nil)
			request.Header.Set("token", "test-token-1abcdef")
		})

		JustBeforeEach(func() {
			authClient.ValidateTokenHeader(request)
		})

		It("should grab the token from headers", func() {
			Expect(fakeToken.ValidateTokenArgsForCall(0)).To(Equal("test-token-1abcdef"))
		})

		It("should pass validate to a token interface", func() {
			Expect(fakeToken.ValidateTokenCallCount()).To(Equal(1))
		})
	})
})
