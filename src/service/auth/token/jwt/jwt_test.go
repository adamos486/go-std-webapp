package jwt_test

import (
	"service/auth/token/jwt"

	jwt2 "github.com/dgrijalva/jwt-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWT Service Specs", func() {
	var jwtService *jwt.Service
	var validTestToken string
	var inputMap map[string]interface{}

	BeforeEach(func() {
		jwtService = jwt.NewService()

		inputMap = make(map[string]interface{})
		inputMap["email"] = "test@gmail.com"
		var err error
		validTestToken, err = jwtService.Generate(inputMap)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("ValidateToken logic", func() {

		Context("when an emtpy string is passed in", func() {
			var (
				expanded   interface{}
				isValid    bool
				executeErr error
			)

			JustBeforeEach(func() {
				expanded, isValid, executeErr = jwtService.ValidateToken("")
			})

			It("should return a nil interface and a false and an error", func() {
				Expect(expanded).To(BeNil())
				Expect(isValid).To(BeFalse())
				Expect(executeErr).To(HaveOccurred())
				Expect(executeErr.Error()).To(Equal("cannot validate an empty token"))
			})
		})

		Context("when an invalid token is passed in", func() {
			var (
				isValid    bool
				executeErr error
			)

			JustBeforeEach(func() {
				_, isValid, executeErr = jwtService.ValidateToken("alskdaskdglsadkasdg")
			})

			It("should return nil, false, and the proper error", func() {
				Expect(executeErr).To(HaveOccurred())
				Expect(isValid).To(BeFalse())
			})
		})

		Context("when a valid token is passed in", func() {
			var (
				expanded   interface{}
				isValid    bool
				executeErr error
			)

			BeforeEach(func() {
			})

			JustBeforeEach(func() {
				expanded, isValid, executeErr = jwtService.ValidateToken(validTestToken)
			})

			It("should return an expanded token object", func() {
				Expect(executeErr).ToNot(HaveOccurred())
				Expect(expanded).ToNot(BeNil())
				Expect(isValid).To(BeTrue())
			})

			It("should have our identity claims bound to it", func() {
				claims := expanded.(*jwt2.Token).Claims
				identityClaims, ok := claims.(*jwt.IdentityClaims)
				Expect(ok).To(BeTrue())
				Expect(identityClaims.Email).To(Equal("test@gmail.com"))
			})

			It("should have an ExpiresAt field", func() {
				claims := expanded.(*jwt2.Token).Claims
				identityClaims, ok := claims.(*jwt.IdentityClaims)
				Expect(ok).To(BeTrue())
				Expect(identityClaims.ExpiresAt).ToNot(BeNil())
			})
		})
	})

	Context("Generate Logic", func() {
		var tokenString string
		var executeErr error

		Context("when an empty input map is supplied", func() {

			JustBeforeEach(func() {
				m := make(map[string]interface{})
				tokenString, executeErr = jwtService.Generate(m)
			})

			It("should return an empty token string and an error", func() {
				Expect(tokenString).To(BeEmpty())
				Expect(executeErr).To(HaveOccurred())
				Expect(executeErr.Error()).To(Equal("can't generate a token with required claims"))
			})
		})

		Context("when a nil input map is supplied", func() {
			JustBeforeEach(func() {
				tokenString, executeErr = jwtService.Generate(nil)
			})

			It("should return an empty token string and an error", func() {
				Expect(tokenString).To(BeEmpty())
				Expect(executeErr).To(HaveOccurred())
				Expect(executeErr.Error()).To(Equal("can't generate a token with required claims"))
			})
		})
	})
})
