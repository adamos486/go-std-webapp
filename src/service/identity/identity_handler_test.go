package identity_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"service/auth"
	"service/auth/authfakes"
	"service/auth/token/tokenfakes"
	"service/identity"
	"service/identity/identityfakes"
	"service/log/logfakes"
	"time"

	"github.com/go-chi/chi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ = Describe("Identity Handler Specs", func() {
	Context("identity routes", func() {
		var (
			identityHandler *identity.HandlerObject
			fakeService     *identityfakes.FakeServiceInterface //Fake
			fakeAuth        *authfakes.FakeInterface
			fakeToken       *tokenfakes.FakeInterface
			fakeAuthClient  *auth.Client
			fakeLog         *logfakes.FakeProdInterface
			router          *chi.Mux
			server          *httptest.Server
		)

		BeforeEach(func() {
			fakeLog = &logfakes.FakeProdInterface{}
			fakeService = &identityfakes.FakeServiceInterface{}
			fakeAuth = &authfakes.FakeInterface{}
			fakeToken = &tokenfakes.FakeInterface{}
			fakeAuthClient = auth.NewClient(fakeAuth, fakeToken)

			identityHandler = identity.NewHandlerObject(fakeLog, fakeService, fakeAuthClient)
		})

		Context("create identity routes", func() {

			BeforeEach(func() {
				router = chi.NewRouter()
				router.Post("/identity", identityHandler.CreateIdentity)
				server = httptest.NewServer(router)
			})

			Context("POST to /identity", func() {
				It("should not return an error when reading body", func() {
					//Expect(expectedError).ToNot(HaveOccurred())
				})

				Context("when it has a valid post body", func() {
					var (
						request       *http.Request
						recorder      *httptest.ResponseRecorder
						response      *http.Response
						expectedError error
						body          []byte
					)

					BeforeEach(func() {
						buffer := bytes.NewBuffer([]byte(`{"email": "test@gmail.com"}`))
						request = httptest.NewRequest("POST", server.URL+"/identity", buffer)
						request.Header.Set("token", "asdlkgjaskgsadgjadsglasdjgasjdglsjdasgd")
						recorder = httptest.NewRecorder()
					})

					JustBeforeEach(func() {
						fakeService.CreateReturns(&identity.Row{
							ID:          "test_id",
							FirstName:   "test_first",
							LastName:    "test_last",
							ProfileInfo: "test_profile_info",
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
						}, sqlmock.NewResult(0, 1), nil)

						router.ServeHTTP(recorder, request)
						response = recorder.Result()
						body, expectedError = ioutil.ReadAll(response.Body)
					})

					It("should have a token header to create an identity", func() {

					})

					It("returns a 200 OK and a response body", func() {
						Expect(recorder.Code).To(Equal(http.StatusOK))
						Expect(expectedError).ToNot(HaveOccurred())
						Expect(body).ToNot(BeNil())
						Expect(body).ToNot(HaveLen(0))
					})
				})

				Context("when it doesn't have a valid post body", func() {
					var (
						request       *http.Request
						recorder      *httptest.ResponseRecorder
						response      *http.Response
						expectedError error
						body          []byte
					)

					BeforeEach(func() {
						request = httptest.NewRequest("POST", server.URL+"/identity", nil)
						recorder = httptest.NewRecorder()
					})

					JustBeforeEach(func() {
						fakeService.CreateReturns(nil, nil, nil)
						router.ServeHTTP(recorder, request)
						response = recorder.Result()
						body, expectedError = ioutil.ReadAll(response.Body)
					})

					It("returns a 400 BAD REQUEST and a standard message", func() {
						Expect(expectedError).ToNot(HaveOccurred())
						Expect(body).ToNot(BeEmpty())
						//Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})
				})
			})
		})

		Context("when a user visits identity index", func() {

			It("identity index handler should return nothing", func() {
				Expect(fakeService.FetchCallCount()).To(Equal(0))
			})
		})

		Context("when a user hits identity with an id query param", func() {

			var (
				request       *http.Request
				recorder      *httptest.ResponseRecorder
				response      *http.Response
				expectedError error
				body          []byte
			)

			BeforeEach(func() {
				router = chi.NewRouter()
				router.Get("/identity", identityHandler.Handler)
				server = httptest.NewServer(router)
			})

			Context("when identity_service returns db errors, ensure they are logged correctly", func() {
				BeforeEach(func() {
					fakeService.FetchReturns(nil, errors.New("fetch error"))
				})

				JustBeforeEach(func() {
					request = httptest.NewRequest("GET", server.URL+"/", nil)
					recorder = httptest.NewRecorder()
					identityHandler.Handler(recorder, request)
					response = recorder.Result()
					body, expectedError = ioutil.ReadAll(response.Body)
					_ = expectedError
					_ = body
				})

				It("should log the errors with our log client and be proper", func() {
					Expect(fakeLog.ErrorCallCount()).To(Equal(1))
				})
			})

			Context("when id is defined and valid, return a single identity row", func() {
				BeforeEach(func() {
					rightNow := time.Now()
					testRow := identity.Row{
						ID:          "test_id",
						FirstName:   "test_first",
						LastName:    "test_last",
						ProfileInfo: "test_profile_info",
						CreatedAt:   rightNow,
						UpdatedAt:   rightNow,
					}
					fakeService.FetchReturns(&testRow, nil)
				})

				JustBeforeEach(func() {
					request = httptest.NewRequest("GET", server.URL+"/", nil)
					recorder = httptest.NewRecorder()
					identityHandler.Handler(recorder, request)
					response = recorder.Result()
					body, expectedError = ioutil.ReadAll(response.Body)
				})

				It("should fetch the record by that ID", func() {
					Expect(fakeService.FetchCallCount()).To(Equal(1))
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})
		Context("when a user wants to authorize an identity", func() {
			var (
				request  *http.Request
				recorder *httptest.ResponseRecorder
				response *http.Response
				body     []byte
			)

			BeforeEach(func() {
				router = chi.NewRouter()
				router.Get("/identity/auth", identityHandler.AuthIdentity)
				server = httptest.NewServer(router)
			})

			Context("when a user makes invalid requests", func() {

				BeforeEach(func() {
					request = httptest.NewRequest("POST", server.URL+"/auth", nil)
					recorder = httptest.NewRecorder()
				})

				Context("and is missing a post body", func() {

					JustBeforeEach(func() {
						identityHandler.AuthIdentity(recorder, request)
						response = recorder.Result()
						body, _ = ioutil.ReadAll(response.Body)
					})

					It("should respond with a bad request and a message", func() {
						Expect(string(body)).Should(Equal("bad request"))
					})
				})

				Context("and has a bad json post body", func() {
					It("should respond with a bad request and a message", func() {

					})
				})

				Context("and is missing an id or password", func() {
					It("should respond with a bad request and a message", func() {

					})
				})
			})

			Context("when a user has a valid identity", func() {
				It("should generate a new access token", func() {

				})

				It("the new access token should contain all of the necessary identity data", func() {

				})

				It("should return the proper response body and code", func() {

				})
			})
		})
	})
})
