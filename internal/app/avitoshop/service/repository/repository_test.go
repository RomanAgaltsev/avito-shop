package repository_test

import (
	"context"
	"errors"
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
		err                 error
		errSomethingStrange error

		ctx context.Context
		bo  *backoff.ExponentialBackOff

		mockPool pgxmock.PgxPoolIface
		repo     *repository.Repository

		rowID int32

		username string
		password string

		user       model.User
		expectUser model.User
	)

	BeforeEach(func() {
		errSomethingStrange = errors.New("something strange")

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

	Context("Calling SendCoins method", func() {
		BeforeEach(func() {
			username = "user"
			password = "password"

			user = model.User{
				UserName: username,
				Password: password,
			}
		})

		When("balance is enough to send", func() {
			BeforeEach(func() {
				rowID = 1

				var balanceSender int32 = 900
				var balanceReceiver int32 = 1100
				var amount int32 = 100
				var toUser model.User = model.User{
					UserName: "user1",
					Password: "password1",
				}

				mockPool.ExpectBegin()

				rsUpdateSender := pgxmock.NewRows([]string{"balance"}).AddRow(balanceSender)
				mockPool.ExpectQuery("UPDATE balance SET .+").WithArgs(username, -amount).WillReturnRows(rsUpdateSender).Times(1)

				rsCreateSender := pgxmock.NewRows([]string{"id"}).AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO history .+ VALUES .+").WithArgs(username, "", toUser.UserName, amount).WillReturnRows(rsCreateSender).Times(1)

				rsUpdateReceiver := pgxmock.NewRows([]string{"balance"}).AddRow(balanceReceiver)
				mockPool.ExpectQuery("UPDATE balance SET .+").WithArgs(toUser.UserName, amount).WillReturnRows(rsUpdateReceiver).Times(1)

				rsCreateReceiver := pgxmock.NewRows([]string{"id"}).AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO history .+ VALUES .+").WithArgs(toUser.UserName, username, "", amount).WillReturnRows(rsCreateReceiver).Times(1)

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()

				err = repo.SendCoins(ctx, bo, user, toUser, int(amount))
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns nil error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("balance is not enough to send", func() {
			BeforeEach(func() {
				rowID = 1

				var balanceSender int32 = 90
				var amount int32 = 100
				var toUser model.User = model.User{
					UserName: "user1",
					Password: "password1",
				}

				mockPool.ExpectBegin()

				rsUpdateSender := pgxmock.NewRows([]string{"balance"}).AddRow(-balanceSender)
				mockPool.ExpectQuery("UPDATE balance SET .+").WithArgs(username, -amount).WillReturnRows(rsUpdateSender).Times(1)

				mockPool.ExpectRollback()

				err = repo.SendCoins(ctx, bo, user, toUser, int(amount))
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns negative balance error", func() {
				Expect(err).To(Equal(repository.ErrNegativeBalance))
			})
		})
	})

	Context("Calling BuyItem method", func() {
		BeforeEach(func() {
			username = "user"
			password = "password"

			user = model.User{
				UserName: username,
				Password: password,
			}
		})

		When("balance is enough to buy", func() {
			BeforeEach(func() {
				rowID = 1

				var balance int32 = 900
				var itemType string = "book"
				var itemPrice int32 = 50

				var item model.InventoryItem = model.InventoryItem{
					Type:     itemType,
					Quantity: 1,
				}

				rsGet := pgxmock.NewRows([]string{"id", "type", "price"}).AddRow(rowID, itemType, itemPrice)
				mockPool.ExpectQuery("SELECT .+ FROM merch .+").WithArgs(itemType).WillReturnRows(rsGet).Times(1)

				mockPool.ExpectBegin()

				rsWithdraw := pgxmock.NewRows([]string{"balance"}).AddRow(balance)
				mockPool.ExpectQuery("UPDATE balance SET .+").WithArgs(username, -itemPrice).WillReturnRows(rsWithdraw).Times(1)

				rsCreate := pgxmock.NewRows([]string{"id"}).AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO inventory .+ VALUES .+").WithArgs(username, itemType).WillReturnRows(rsCreate).Times(1)

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()

				err = repo.BuyItem(ctx, bo, user, item)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns nil error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("balance is not enough to buy", func() {
			BeforeEach(func() {
				rowID = 1

				var balance int32 = 90
				var itemType string = "book"
				var itemPrice int32 = 100

				var item model.InventoryItem = model.InventoryItem{
					Type:     itemType,
					Quantity: 1,
				}

				rsGet := pgxmock.NewRows([]string{"id", "type", "price"}).AddRow(rowID, itemType, itemPrice)
				mockPool.ExpectQuery("SELECT .+ FROM merch .+").WithArgs(itemType).WillReturnRows(rsGet).Times(1)

				mockPool.ExpectBegin()

				rsWithdraw := pgxmock.NewRows([]string{"balance"}).AddRow(-balance)
				mockPool.ExpectQuery("UPDATE balance SET .+").WithArgs(username, -itemPrice).WillReturnRows(rsWithdraw).Times(1)

				mockPool.ExpectRollback()

				err = repo.BuyItem(ctx, bo, user, item)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns negative balance error", func() {
				Expect(err).To(Equal(repository.ErrNegativeBalance))
			})
		})
	})

	Context("Calling GetBalance method", func() {
		BeforeEach(func() {
			username = "user"
			password = "password"

			user = model.User{
				UserName: username,
				Password: password,
			}
		})

		When("there is no error", func() {
			BeforeEach(func() {
				rowID = 1

				var balance int32 = 1000

				rs := pgxmock.NewRows([]string{"id", "username", "coins"}).AddRow(rowID, username, balance)
				mockPool.ExpectQuery("SELECT .+ FROM balance .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns a balance and nil error", func() {
				result, err := repo.GetBalance(ctx, bo, user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(Equal(int(1000)))
			})
		})

		When("there is an error", func() {
			BeforeEach(func() {
				rowID = 0

				rs := pgxmock.NewRows([]string{"id"}).AddRow(rowID).RowError(int(rowID), errSomethingStrange)
				mockPool.ExpectQuery("SELECT .+ FROM balance .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns zero balance and an error", func() {
				result, err := repo.GetBalance(ctx, bo, user)
				Expect(err).Should(HaveOccurred())
				Expect(result).To(Equal(int(0)))
			})
		})
	})

	Context("Calling GetInventory method", func() {
		BeforeEach(func() {
			username = "user"
			password = "password"

			user = model.User{
				UserName: username,
				Password: password,
			}
		})

		When("there is no error", func() {
			BeforeEach(func() {
				rowID = 1

				var itemType string = "book"
				var itemQuantity int64 = 1

				rs := pgxmock.NewRows([]string{"type", "quantity"}).AddRow(itemType, itemQuantity)
				mockPool.ExpectQuery("SELECT .+ FROM inventory .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an inventory and nil error", func() {
				result, err := repo.GetInventory(ctx, bo, user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).ShouldNot(BeNil())
				Expect(result).Should(HaveLen(1))
				Expect(result[0].Type).To(Equal("book"))
				Expect(result[0].Quantity).To(Equal(1))
			})
		})

		When("there is an error", func() {
			BeforeEach(func() {
				rowID = 0

				rs := pgxmock.NewRows([]string{"id"}).AddRow(rowID).RowError(int(rowID), errSomethingStrange)
				mockPool.ExpectQuery("SELECT .+ FROM inventory .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an empty inventory and an error", func() {
				result, err := repo.GetInventory(ctx, bo, user)
				Expect(err).Should(HaveOccurred())
				Expect(result).Should(BeNil())
			})
		})
	})

	Context("Calling GetHistory method", func() {
		BeforeEach(func() {
			username = "user"
			password = "password"

			user = model.User{
				UserName: username,
				Password: password,
			}
		})

		When("there is no error", func() {
			BeforeEach(func() {
				rowID = 1

				var fromUser string = "user1"
				var toUser string = ""
				var amount int32 = 100
				var sentAt time.Time = time.Now()

				rs := pgxmock.NewRows([]string{"id", "username", "fromuser", "touser", "amount", "sentat"}).AddRow(rowID, username, fromUser, toUser, amount, sentAt)
				mockPool.ExpectQuery("SELECT .+ FROM history .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an inventory and nil error", func() {
				result, err := repo.GetHistory(ctx, bo, user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Received).Should(HaveLen(1))
				Expect(result.Received[0].FromUser).Should(Equal("user1"))
				Expect(result.Received[0].Amount).Should(Equal(100))
				Expect(result.Sent).Should(HaveLen(0))
			})
		})

		When("there is an error", func() {
			BeforeEach(func() {
				rowID = 0

				rs := pgxmock.NewRows([]string{"id"}).AddRow(rowID).RowError(int(rowID), errSomethingStrange)
				mockPool.ExpectQuery("SELECT .+ FROM history .+").WithArgs(username).WillReturnRows(rs).Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an empty inventory and an error", func() {
				result, err := repo.GetHistory(ctx, bo, user)
				Expect(err).Should(HaveOccurred())
				Expect(result.Received).Should(HaveLen(0))
				Expect(result.Sent).Should(HaveLen(0))
			})
		})
	})
})
