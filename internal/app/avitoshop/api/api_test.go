package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.uber.org/mock/gomock"

	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/api"
	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/repository"
	"github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop/service/shop"
	"github.com/RomanAgaltsev/avito-shop/internal/config"
	"github.com/RomanAgaltsev/avito-shop/internal/mock"
	"github.com/RomanAgaltsev/avito-shop/internal/model"
	"github.com/RomanAgaltsev/avito-shop/internal/pkg/auth"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

var _ = Describe("Handler", func() {
	var (
		err                 error
		errSomethingStrange error

		cfg *config.Config

		server *ghttp.Server

		service shop.Service
		ctrl    *gomock.Controller
		repo    *mock.MockRepository

		handler *api.Handler

		endpoint string

		user      model.User
		userBytes []byte

		coinsSending      model.CoinsSending
		coinsSendingBytes []byte

		ja *jwtauth.JWTAuth

		expectAuthResponse model.AuthResponse
		expectBalance      int
		expectInventory    []model.InventoryItem
		expectHistory      model.CoinsHistory

		username   string
		toUsername string
		secretKey  string
		token      string
		itemType   string
	)

	BeforeEach(func() {
		errSomethingStrange = errors.New("something strange")

		cfg, err = config.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).ShouldNot(BeNil())

		server = ghttp.NewServer()

		ctrl = gomock.NewController(GinkgoT())
		Expect(ctrl).ShouldNot(BeNil())

		repo = mock.NewMockRepository(ctrl)
		Expect(repo).ShouldNot(BeNil())

		service, err = shop.NewService(repo, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(service).ShouldNot(BeNil())

		handler = api.NewHandler(cfg, service)
		Expect(handler).ShouldNot(BeNil())
	})

	AfterEach(func() {
		server.Close()
	})

	Context("Receiving request at the /api/auth endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/auth"
			server.AppendHandlers(handler.Auth)
		})

		When("the method is POST, content type is right and payload is right", func() {
			BeforeEach(func() {
				user = model.User{
					UserName: "user",
					Password: "password",
				}

				expectAuthResponse = model.AuthResponse{}

				userBytes, err = json.Marshal(&user)
				Expect(err).ShouldNot(HaveOccurred())

				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil).Times(1)
				repo.EXPECT().CreateBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			})

			It("returns status 'OK' (200) and a token", func() {
				response, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				err = json.NewDecoder(response.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).NotTo(BeEmpty())
				Expect(expectAuthResponse.Errors).To(BeEmpty())
			})
		})

		When("the method is POST, content type is right but payload is wrong", func() {
			BeforeEach(func() {
				user = model.User{
					UserName: "user",
					Password: "",
				}

				expectAuthResponse = model.AuthResponse{}

				userBytes, err = json.Marshal(user)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Bad request' (400) and no token", func() {
				response, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))

				err = json.NewDecoder(response.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).To(BeEmpty())
				Expect(expectAuthResponse.Errors).NotTo(BeEmpty())
			})
		})

		When("the method is POST, content type is wrong and payload is right", func() {
			BeforeEach(func() {
				user = model.User{
					UserName: "user",
					Password: "password",
				}

				expectAuthResponse = model.AuthResponse{}

				userBytes, err = json.Marshal(user)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Bad request' (400) and no token", func() {
				response, err := http.Post(server.URL()+endpoint, ContentTypeText, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))

				err = json.NewDecoder(response.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).To(BeEmpty())
				Expect(expectAuthResponse.Errors).NotTo(BeEmpty())
			})
		})

		When("the method is POST, request is right but user already exists and password is correct", func() {
			BeforeEach(func() {
				user = model.User{
					UserName: "user",
					Password: "password",
				}

				expectAuthResponse = model.AuthResponse{}

				userBytes, err = json.Marshal(user)
				Expect(err).ShouldNot(HaveOccurred())

				hash, err := auth.HashPassword(user.Password)
				Expect(err).ShouldNot(HaveOccurred())

				user.Password = hash

				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, repository.ErrConflict).Times(1)
			})

			It("returns status 'OK' (200) and a token", func() {
				response, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				err = json.NewDecoder(response.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).NotTo(BeEmpty())
				Expect(expectAuthResponse.Errors).To(BeEmpty())
			})
		})

		When("the method is POST, request is right but user already exists and password is incorrect", func() {
			BeforeEach(func() {
				user = model.User{
					UserName: "user",
					Password: "wrong password",
				}

				expectAuthResponse = model.AuthResponse{}

				userBytes, err = json.Marshal(user)
				Expect(err).ShouldNot(HaveOccurred())

				hash, err := auth.HashPassword("password")
				Expect(err).ShouldNot(HaveOccurred())

				user.Password = hash

				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, repository.ErrConflict).Times(1)
			})

			It("returns status 'Unauthorized' (401) and no token", func() {
				response, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusUnauthorized))

				err = json.NewDecoder(response.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).To(BeEmpty())
				Expect(expectAuthResponse.Errors).NotTo(BeEmpty())
			})
		})

		When("everything is right with the request, but something has gone wrong with the user creation", func() {
			BeforeEach(func() {
				user = model.User{
					UserName: "user",
					Password: "password",
				}

				expectAuthResponse = model.AuthResponse{}

				userBytes, err = json.Marshal(user)
				Expect(err).ShouldNot(HaveOccurred())

				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.User{}, errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				response, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))

				err = json.NewDecoder(response.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).To(BeEmpty())
				Expect(expectAuthResponse.Errors).NotTo(BeEmpty())
			})
		})

		When("everything is right with the request, but something has gone wrong with the user balance creation", func() {
			BeforeEach(func() {
				user = model.User{
					UserName: "user",
					Password: "password",
				}

				expectAuthResponse = model.AuthResponse{}

				userBytes, err = json.Marshal(user)
				Expect(err).ShouldNot(HaveOccurred())

				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(user, nil).Times(1)
				repo.EXPECT().CreateBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				response, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))

				err = json.NewDecoder(response.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).To(BeEmpty())
				Expect(expectAuthResponse.Errors).NotTo(BeEmpty())
			})
		})
	})

	Context("Receiving request at the /api/sendCoin endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/sendCoin"
			server.AppendHandlers(handler.SendCoins)

			secretKey = "secret"
			username = "user"

			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())

			_, token, err = auth.NewJWTToken(ja, username)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())
		})

		When("the method is POST and balance enough to send", func() {
			BeforeEach(func() {
				toUsername = "user1"

				coinsSending = model.CoinsSending{
					ToUser: toUsername,
					Amount: 100,
				}

				coinsSendingBytes, err = json.Marshal(&coinsSending)
				Expect(err).ShouldNot(HaveOccurred())

				repo.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			})

			It("returns status 'OK' (200)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(coinsSendingBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))
			})
		})

		When("the method is POST and balance is not enough to send", func() {
			BeforeEach(func() {
				toUsername = "user1"

				coinsSending = model.CoinsSending{
					ToUser: toUsername,
					Amount: 100,
				}

				coinsSendingBytes, err = json.Marshal(&coinsSending)
				Expect(err).ShouldNot(HaveOccurred())

				repo.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.ErrNegativeBalance).Times(1)
			})

			It("returns status 'Bad request' (400)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(coinsSendingBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))
			})
		})

		When("the method is POST and payload is invalid", func() {
			BeforeEach(func() {
				toUsername = ""

				coinsSending = model.CoinsSending{
					ToUser: toUsername,
					Amount: 100,
				}

				coinsSendingBytes, err = json.Marshal(&coinsSending)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Bad request' (400)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(coinsSendingBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))
			})
		})

		When("the method is POST, everything is right with the request, but something has gone wrong with service", func() {
			BeforeEach(func() {
				toUsername = "user1"

				coinsSending = model.CoinsSending{
					ToUser: toUsername,
					Amount: 100,
				}

				coinsSendingBytes, err = json.Marshal(&coinsSending)
				Expect(err).ShouldNot(HaveOccurred())

				repo.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(coinsSendingBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("Receiving request at the /api/buy/ endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/buy/"
			server.AppendHandlers(handler.BuyItem)

			secretKey = "secret"
			username = "user"

			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())

			_, token, err = auth.NewJWTToken(ja, username)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())
		})

		When("the method is GET and everything is right", func() {
			BeforeEach(func() {
				itemType = "book"

				repo.EXPECT().BuyItem(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			})

			It("returns status 'OK' (200)", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint+itemType, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))
			})
		})

		When("the method is GET and item is empty", func() {
			BeforeEach(func() {
				itemType = ""
			})

			It("returns status 'Bad request' (400)", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint+itemType, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))
			})
		})
	})

	Context("Receiving request at the /api/info endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/info"
			server.AppendHandlers(handler.Info)

			secretKey = "secret"
			username = "user"

			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())

			_, token, err = auth.NewJWTToken(ja, username)
			Expect(err).NotTo(HaveOccurred())
			Expect(token).NotTo(BeEmpty())
		})

		When("the method is GET and there is an info to return", func() {
			BeforeEach(func() {
				expectBalance = 100
				expectInventory = []model.InventoryItem{
					{Type: "book", Quantity: 1},
				}
				expectHistory = model.CoinsHistory{
					Received: []model.CoinsReceiving{
						{FromUser: "user1", Amount: 100},
					},
					Sent: []model.CoinsSending{
						{ToUser: "user2", Amount: 200},
					},
				}

				repo.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectBalance, nil).Times(1)
				repo.EXPECT().GetInventory(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectInventory, nil).Times(1)
				repo.EXPECT().GetHistory(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectHistory, nil).Times(1)
			})

			It("returns status 'OK' (200) and an info", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var info model.Info
				err = json.NewDecoder(response.Body).Decode(&info)
				Expect(info.Coins).To(Equal(expectBalance))
				Expect(info.Inventory).Should(HaveLen(len(expectInventory)))
				Expect(info.CoinsHistory.Received).Should(HaveLen(len(expectHistory.Received)))
				Expect(info.CoinsHistory.Sent).Should(HaveLen(len(expectHistory.Sent)))
			})
		})

		When("the method is GET and there is no info to return", func() {
			BeforeEach(func() {
				expectBalance = 0
				expectInventory = []model.InventoryItem{}
				expectHistory = model.CoinsHistory{
					Received: []model.CoinsReceiving{},
					Sent:     []model.CoinsSending{},
				}

				repo.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectBalance, nil).Times(1)
				repo.EXPECT().GetInventory(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectInventory, nil).Times(1)
				repo.EXPECT().GetHistory(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectHistory, nil).Times(1)
			})

			It("returns status 'OK' (200) and no info", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var info model.Info
				err = json.NewDecoder(response.Body).Decode(&info)
				Expect(info.Coins).To(Equal(expectBalance))
				Expect(info.Inventory).Should(HaveLen(len(expectInventory)))
				Expect(info.CoinsHistory.Received).Should(HaveLen(len(expectHistory.Received)))
				Expect(info.CoinsHistory.Sent).Should(HaveLen(len(expectHistory.Sent)))
			})
		})

		When("the method is GET, but something has gone wrong with the service", func() {
			BeforeEach(func() {
				expectBalance = 0

				repo.EXPECT().GetBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectBalance, errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Add("Authorization", "Bearer "+token)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})

	})

})
