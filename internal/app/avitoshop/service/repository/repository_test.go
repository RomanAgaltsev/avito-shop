package repository_test

import (
	"context"
	//"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock/v4"

	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/repository"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

var _ = Describe("Repository", func() {
	var (
		err error
		//errSomethingStrange error

		ctx context.Context
		bo  *backoff.ExponentialBackOff

		mockPool pgxmock.PgxPoolIface
		repo     *repository.Repository

		rowID int32

		username string
		password string
		//userCreatedAt time.Time

		user       model.User
		expectUser model.User
	)

	BeforeEach(func() {
		//errSomethingStrange = errors.New("something strange")

		ctx = context.Background()

		bo = backoff.NewExponentialBackOff()
		bo.InitialInterval = 50 * time.Millisecond
		bo.RandomizationFactor = 0.1
		bo.Multiplier = 2.0
		bo.MaxInterval = 1 * time.Second
		bo.MaxElapsedTime = 2 * time.Second
		bo.Reset()

		mockPool, err = pgxmock.NewPool()
		Expect(err).ShouldNot(HaveOccurred())

		repo, err = repository.New(mockPool)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		mockPool.Close()
	})

	Context("Calling CreateUser method", func() {
		BeforeEach(func() {
			username = "user"
			password = "password"

			user = model.User{
				UserName: username,
				Password: password,
			}
		})

		When("user doesn't exist", func() {
			BeforeEach(func() {
				rowID = 1

				rs := pgxmock.NewRows([]string{"id"}).AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO users .+ VALUES .+").WithArgs(username, password).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns nil error and a user", func() {
				expectUser, err = repo.CreateUser(ctx, bo, user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectUser.UserName).To(Equal(username))
				Expect(expectUser.Password).To(Equal(password))
			})
		})

		When("user already exist", func() {
			BeforeEach(func() {
				rowID = 0

				createdAt := time.Now()

				rsCreate := pgxmock.NewRows([]string{"id"}).AddRow(rowID).RowError(int(rowID), &pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation})
				mockPool.ExpectQuery("INSERT INTO users .+ VALUES .+").WithArgs(username, password).WillReturnRows(rsCreate).Times(1)

				rsGet := pgxmock.NewRows([]string{"id", "username", "password", "createdat"}).AddRow(rowID, username, password, createdAt)
				mockPool.ExpectQuery("SELECT .+ FROM users .+").WithArgs(username).WillReturnRows(rsGet).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns data conflict error", func() {
				expectUser, err = repo.CreateUser(ctx, bo, user)
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(repository.ErrConflict))
				Expect(expectUser.UserName).To(Equal(username))
				Expect(expectUser.Password).To(Equal(password))
			})
		})

	})

	Context("Calling CreateBalance method", func() {
		BeforeEach(func() {
			username = "user"
			password = "password"

			user = model.User{
				UserName: username,
				Password: password,
			}
		})

		When("balance doesn't exist", func() {
			BeforeEach(func() {
				rowID = 1

				rs := pgxmock.NewRows([]string{"id"}).AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO balance .+ VALUES .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns nil error", func() {
				err = repo.CreateBalance(ctx, bo, user)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("balance exists", func() {
			BeforeEach(func() {
				rowID = 0

				rs := pgxmock.NewRows([]string{"id"}).AddRow(rowID).RowError(int(rowID), &pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation})
				mockPool.ExpectQuery("INSERT INTO balance .+ VALUES .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an error", func() {
				err = repo.CreateBalance(ctx, bo, user)
				Expect(err).Should(HaveOccurred())
			})
		})

	})

})
