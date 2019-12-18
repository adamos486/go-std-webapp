package main_test

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"service/auth"
	"service/auth/authfakes"
	"service/auth/basic"
	"service/auth/token/tokenfakes"
	"service/database/databasefakes"
	"service/handlers/index"
	"service/handlers/request"
	"service/log"
	"service/log/logfakes"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main Specs", func() {
	Context("index", func() {
		var (
			router       *chi.Mux
			server       *httptest.Server
			indexHandler *index.Index
			logFake      *logfakes.FakeProdInterface
			fakeDB       *databasefakes.FakeDBInterface
			authFake     *authfakes.FakeInterface
			tokenFake    *tokenfakes.FakeInterface
			db           *sql.DB
			mockDB       sqlmock.Sqlmock
			err          error
		)

		BeforeEach(func() {
			logFake = &logfakes.FakeProdInterface{}
			logClient := log.New(logFake)

			fakeDB = &databasefakes.FakeDBInterface{}
			//initialize a sqlmock
			db, mockDB, err = sqlmock.New()
			Expect(err).ToNot(HaveOccurred())

			indexHandler = index.New(logClient, fakeDB)
			router = chi.NewRouter()

			authFake = &authfakes.FakeInterface{}
			tokenFake = &tokenfakes.FakeInterface{}
			fakeAuthClient := auth.NewClient(authFake, tokenFake)

			router.Use(request.GenerateRequestIDMiddle)

			request.SetupLogger(logClient)
			router.Use(request.Logger)

			basic.SetupAuthMiddleware(fakeAuthClient, logClient)
			router.Use(basic.AuthMiddleware)

			server = httptest.NewServer(router)
		})

		Context("when a user hits the index route", func() {
			var (
				request       *http.Request
				recorder      *httptest.ResponseRecorder
				response      *http.Response
				expectedError error
				body          []byte
			)

			BeforeEach(func() {
				now := time.Now()
				mockRows := sqlmock.NewRows([]string{"id", "name", "description", "date_added"})
				mockRows = mockRows.AddRow(0, "test concert", "test description", now)
				mockDB.ExpectQuery("SELECT id, name, description FROM event;").WillReturnRows(mockRows)
				rows, _ := db.Query("SELECT id, name, description FROM event;")
				fakeDB.QueryReturns(rows, nil)
			})

			Context("when a user has the right credentials", func() {
				BeforeEach(func() {

					authFake.AuthorizeReturns("tony", "house", true)

					router.Get("/", indexHandler.Handler)
					request = httptest.NewRequest("GET", server.URL+"/", nil)
					request.SetBasicAuth("tony", "house")
					recorder = httptest.NewRecorder()

				})

				JustBeforeEach(func() {
					router.ServeHTTP(recorder, request)
					response = recorder.Result()
					body, expectedError = ioutil.ReadAll(response.Body)
				})
				It("should not have errored at all", func() {
					Expect(expectedError).ToNot(HaveOccurred())
				})

				It("should have the proper status code", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
				})

				It("should have the expected body", func() {
					var resStruct index.EventsResponse
					err := json.Unmarshal(body, &resStruct)
					Expect(err).ToNot(HaveOccurred())
					Expect(resStruct.Code).To(Equal(http.StatusOK))
					Expect(resStruct.List).To(HaveLen(1))
				})

				It("should call logging only in expected cases", func() {
					Expect(logFake.InfoCallCount()).To(Equal(2))
					Expect(logFake.ErrorCallCount()).To(Equal(0))
					Expect(logFake.DebugCallCount()).To(Equal(2))
				})
			})

			Context("when a user has incorrect credentials", func() {
				BeforeEach(func() {
					router.Get("/", indexHandler.Handler)
					request = httptest.NewRequest("GET", server.URL+"/", nil)
					request.SetBasicAuth("tony", "house1")
					authFake.AuthorizeReturns("tony", "house1", false)
					recorder = httptest.NewRecorder()
				})

				JustBeforeEach(func() {
					router.ServeHTTP(recorder, request)
					response = recorder.Result()
					body, expectedError = ioutil.ReadAll(response.Body)
				})

				It("should not have errored at all", func() {
					Expect(response.StatusCode).ToNot(Equal(http.StatusInternalServerError))
					Expect(expectedError).ToNot(HaveOccurred())
				})

				It("should have a proper 401 status code and message", func() {
					Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				})

				It("should call logging only in expected cases", func() {
					Expect(logFake.InfoCallCount()).To(Equal(2))
					Expect(logFake.ErrorCallCount()).To(Equal(0))
					Expect(logFake.DebugCallCount()).To(Equal(0))
				})
			})
		})
	})
})
