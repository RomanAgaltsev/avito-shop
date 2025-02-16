package end2end_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/RomanAgaltsev/avito-shop/internal/model"
)

const contentTypeJSON = "application/json"

var _ = Describe("End2End", Serial, func() {
	var tokenUser1 string
	var tokenUser2 string

	Context("Testing item buying scenario", Serial, func() {
		var httpClient http.Client

		When("new user tries to register", Ordered, func() {
			var username string
			var password string

			BeforeAll(func() {
				httpClient = http.Client{}

				username = "user1"
				password = "password1"
			})

			It("app returns status 'OK' (200), a token and no error", func() {
				user := model.User{
					UserName: username,
					Password: password,
				}
				reqBytes, _ := json.Marshal(user)

				response, err := httpClient.Post(
					fmt.Sprintf("%s/api/auth", serverAddr),
					contentTypeJSON,
					bytes.NewReader(reqBytes),
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var authResponse model.AuthResponse
				err = json.NewDecoder(response.Body).Decode(&authResponse)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(authResponse.Token).NotTo(BeEmpty())
				Expect(authResponse.Errors).To(BeEmpty())

				tokenUser1 = authResponse.Token
			})
			It("app increases user balance by 1000 coins and has no user inventory and coins history", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/info", serverAddr),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser1)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var info model.Info
				err = json.NewDecoder(response.Body).Decode(&info)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(info.Coins).To(Equal(1000))
				Expect(info.Inventory).Should(HaveLen(0))
				Expect(info.CoinsHistory.Received).Should(HaveLen(0))
				Expect(info.CoinsHistory.Sent).Should(HaveLen(0))
			})
		})

		When("registered user tries to buy an existed item and has enough coins", Ordered, func() {
			var item string
			var itemPrice int

			BeforeAll(func() {
				httpClient = http.Client{}

				item = "book"
				itemPrice = 50
			})

			It("app returns status 'OK' (200)", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/buy/%s", serverAddr, item),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser1)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))
			})
			It("app decreases user balance by item price and adds item to the user inventory", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/info", serverAddr),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser1)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var info model.Info
				err = json.NewDecoder(response.Body).Decode(&info)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(info.Coins).To(Equal(1000 - itemPrice))
				Expect(info.Inventory).Should(HaveLen(1))
				Expect(info.Inventory[0].Type).To(Equal(item))
				Expect(info.Inventory[0].Quantity).To(Equal(1))
				Expect(info.CoinsHistory.Received).Should(HaveLen(0))
				Expect(info.CoinsHistory.Sent).Should(HaveLen(0))
			})
		})

		When("registered user tries to buy an existed item and has not enough coins", Ordered, func() {
			var item string

			BeforeAll(func() {
				httpClient = http.Client{}

				item = "pink-hoody"
			})

			It("app returns status 'Bad request' (400) and an error", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/buy/%s", serverAddr, item),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser1)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				response, err = http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))

				var authResponse model.AuthResponse
				err = json.NewDecoder(response.Body).Decode(&authResponse)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(authResponse.Errors).NotTo(BeEmpty())
				Expect(authResponse.Errors).To(Equal("Not enough coins"))
			})
		})

		When("registered user tries to buy an unknown item", Ordered, func() {
			var item string

			BeforeAll(func() {
				httpClient = http.Client{}

				item = "tomato"
			})

			It("app returns status 'Bad request' (400) and an error", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/buy/%s", serverAddr, item),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser1)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))

				var authResponse model.AuthResponse
				err = json.NewDecoder(response.Body).Decode(&authResponse)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(authResponse.Errors).NotTo(BeEmpty())
				Expect(authResponse.Errors).To(Equal("Unkown merch"))
			})
		})
	})

	Context("Testing coins sending scenario", Serial, func() {
		var httpClient http.Client

		When("new user tries to register", Ordered, func() {
			var username string
			var password string

			BeforeAll(func() {
				httpClient = http.Client{}

				username = "user2"
				password = "password2"
			})

			It("app returns status 'OK' (200), a token and no error", func() {
				user := model.User{
					UserName: username,
					Password: password,
				}
				reqBytes, _ := json.Marshal(user)

				response, err := httpClient.Post(
					fmt.Sprintf("%s/api/auth", serverAddr),
					contentTypeJSON,
					bytes.NewReader(reqBytes),
				)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var authResponse model.AuthResponse
				err = json.NewDecoder(response.Body).Decode(&authResponse)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(authResponse.Token).NotTo(BeEmpty())
				Expect(authResponse.Errors).To(BeEmpty())

				tokenUser2 = authResponse.Token
			})
			It("app increases user balance by 1000 coins and has no user inventory and coins history", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/info", serverAddr),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser2)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var info model.Info
				err = json.NewDecoder(response.Body).Decode(&info)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(info.Coins).To(Equal(1000))
				Expect(info.Inventory).Should(HaveLen(0))
				Expect(info.CoinsHistory.Received).Should(HaveLen(0))
				Expect(info.CoinsHistory.Sent).Should(HaveLen(0))
			})
		})

		When("registered user tries to send coins and has enough amount", Ordered, func() {
			var fromUsername string
			var toUsername string
			var coins int

			BeforeAll(func() {
				httpClient = http.Client{}

				fromUsername = "user2"
				toUsername = "user1"
				coins = 100
			})

			It("app returns status 'OK' (200)", func() {
				coinsSending := model.CoinsSending{
					ToUser: toUsername,
					Amount: coins,
				}

				reqBytes, _ := json.Marshal(coinsSending)

				request, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/api/sendCoin", serverAddr),
					bytes.NewReader(reqBytes),
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser2)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))
			})
			It("app decreases user balance and adds transaction to the user that sends history", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/info", serverAddr),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser2)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var info model.Info
				err = json.NewDecoder(response.Body).Decode(&info)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(info.Coins).To(Equal(1000 - coins))
				Expect(info.CoinsHistory.Sent).Should(HaveLen(1))
				Expect(info.CoinsHistory.Sent[0].ToUser).Should(Equal(toUsername))
				Expect(info.CoinsHistory.Sent[0].Amount).Should(Equal(coins))
			})
			It("app increases user balance by given amount and adds transaction to the user that receives history", func() {
				request, err := http.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/api/info", serverAddr),
					nil,
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser1)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var info model.Info
				err = json.NewDecoder(response.Body).Decode(&info)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(info.Coins).To(Equal(550))
				Expect(info.CoinsHistory.Received).Should(HaveLen(1))
				Expect(info.CoinsHistory.Received[0].FromUser).Should(Equal(fromUsername))
				Expect(info.CoinsHistory.Received[0].Amount).Should(Equal(coins))
			})
		})

		When("registered user tries to send coins and has not enough amount", Ordered, func() {
			var toUsername string
			var coins int

			BeforeAll(func() {
				httpClient = http.Client{}

				toUsername = "user1"
				coins = 2000
			})

			It("app returns status 'Bad request' (400) and an error", func() {
				coinsSending := model.CoinsSending{
					ToUser: toUsername,
					Amount: coins,
				}

				reqBytes, _ := json.Marshal(coinsSending)

				request, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/api/sendCoin", serverAddr),
					bytes.NewReader(reqBytes),
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser2)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))

				var authResponse model.AuthResponse
				err = json.NewDecoder(response.Body).Decode(&authResponse)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(authResponse.Errors).NotTo(BeEmpty())
				Expect(authResponse.Errors).To(Equal("Not enough coins"))
			})
		})

		When("registered user tries to send coins and receiving user does not exist", Ordered, func() {
			var toUsername string
			var coins int

			BeforeAll(func() {
				httpClient = http.Client{}

				toUsername = "user3"
				coins = 20
			})

			It("app returns status 'Bad request' (400) and an error", func() {
				coinsSending := model.CoinsSending{
					ToUser: toUsername,
					Amount: coins,
				}

				reqBytes, _ := json.Marshal(coinsSending)

				request, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/api/sendCoin", serverAddr),
					bytes.NewReader(reqBytes),
				)
				request.Header.Add("Authorization", "Bearer "+tokenUser2)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))

				var authResponse model.AuthResponse
				err = json.NewDecoder(response.Body).Decode(&authResponse)
				DeferCleanup(response.Body.Close)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(authResponse.Errors).NotTo(BeEmpty())
				Expect(authResponse.Errors).To(Equal("Unkown user to send coins"))
			})
		})
	})
})
