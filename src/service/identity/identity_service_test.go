package identity_test

import (
	"database/sql"
	"encoding/json"
	"service/identity"
	"service/log/logfakes"
	"service/utils/sqltest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ = Describe("Identity Service Specs", func() {
	Context("identity service logic", func() {
		var (
			identityService *identity.ServiceObject
			fakeLog         *logfakes.FakeProdInterface
			db              *sql.DB
			mockDB          sqlmock.Sqlmock
		)

		BeforeEach(func() {
			var sqlmockErr error
			db, mockDB, sqlmockErr = sqlmock.New()
			Expect(sqlmockErr).ToNot(HaveOccurred())

			fakeLog = &logfakes.FakeProdInterface{}
		})

		Context("when a user fetches a record", func() {

			var (
				identityRow *identity.Row
				err         error
			)

			BeforeEach(func() {
				rightNow := time.Now()
				mockRows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "profile",
					"created_at", "updated_at"})
				mockRows = mockRows.AddRow("uuidv4", "test_first_name", "test_last_name",
					[]byte(`{"email": "test@gmail.com"}`), rightNow, rightNow)
				mockDB.ExpectQuery(
					`SELECT
						id, first_name, last_name, profile, created_at, updated_at
						FROM identity
						WHERE id = \$1;`).WithArgs("uuidv4").WillReturnRows(mockRows)
				identityService = identity.NewServiceObject(fakeLog, db)
			})

			JustBeforeEach(func() {
				identityRow, err = identityService.Fetch("uuidv4")
			})

			It("should return with no errors", func() {
				Expect(err).ToNot(HaveOccurred())
				mockErr := mockDB.ExpectationsWereMet()
				Expect(mockErr).ToNot(HaveOccurred())
			})

			It("fetch should return a single row with the specified id", func() {
				Expect(identityRow).ToNot(BeNil())

				Expect(identityRow.ID).To(Equal("uuidv4"))
				Expect(identityRow.FirstName).To(Equal("test_first_name"))
				Expect(identityRow.LastName).To(Equal("test_last_name"))

				expectedProfile := make(map[string]string)
				expectedProfile["email"] = "test@gmail.com"

				var key, value string
				switch v := identityRow.ProfileInfo.(type) {
				case map[string]interface{}:
					for s, b := range v {
						key = s
						value = b.(string)
					}
				}

				Expect(key).To(Equal("email"))
				Expect(value).To(Equal("test@gmail.com"))
			})

			It("should only log the proper things", func() {
				Expect(fakeLog.DebugCallCount()).To(Equal(1))
				Expect(fakeLog.InfoCallCount()).To(Equal(0))
				Expect(fakeLog.WarnCallCount()).To(Equal(0))
				Expect(fakeLog.ErrorCallCount()).To(Equal(0))
			})
		})

		Context("when a user creates an identity", func() {
			var (
				identityRow *identity.Row
				err         error
				result      sql.Result
			)

			BeforeEach(func() {
				mockResult := sqlmock.NewResult(0, 1)

				m := make(map[string]string)
				m["email"] = "test@gmail.com"
				rawJSON, _ := json.Marshal(&m)

				rightNow := time.Now()
				mockRows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "profile",
					"created_at", "updated_at"})
				mockRows = mockRows.AddRow("uuidv4", "test_first_name", "test_last_name",
					[]byte(`{"email": "test@gmail.com"}`), rightNow, rightNow)

				mockDB.ExpectExec("INSERT INTO identity").WithArgs(sqltest.AnyString{},
					"adam", "cobb", rawJSON,
					sqltest.AnyTime{}, sqltest.AnyTime{}).WillReturnResult(mockResult)
				mockDB.ExpectQuery(`SELECT id, first_name, last_name, profile, created_at, updated_at`).WithArgs(sqltest.AnyString{}).WillReturnRows(mockRows)
			})

			JustBeforeEach(func() {
				//Doing this will pass in a properly configured mock HERE. Do not do earlier!
				identityService = identity.NewServiceObject(fakeLog, db)
				identityRow, result, err = identityService.Create("uuidv7")
			})

			It("should return the result of the INSERT INTO and contain 1 row changed", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(result).ToNot(BeNil())
				Expect(result.LastInsertId()).To(Equal(int64(0)))
				Expect(result.RowsAffected()).To(Equal(int64(1)))
			})

			It("should return the created row", func() {
				Expect(identityRow).ToNot(BeNil())
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
