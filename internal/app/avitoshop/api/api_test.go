package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	//"github.com/go-chi/jwtauth/v5"
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

		//ja *jwtauth.JWTAuth

		expectAuthResponse model.AuthResponse
		//expectInfo  model.Info

		//username  string
		//secretKey string
		//tokenString string
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
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))

				err = json.NewDecoder(resp.Body).Decode(&expectAuthResponse)
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
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

				err = json.NewDecoder(resp.Body).Decode(&expectAuthResponse)
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
				resp, err := http.Post(server.URL()+endpoint, ContentTypeText, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

				err = json.NewDecoder(resp.Body).Decode(&expectAuthResponse)
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
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))

				err = json.NewDecoder(resp.Body).Decode(&expectAuthResponse)
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

			It("returns status 'OK' (401) and no token", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))

				err = json.NewDecoder(resp.Body).Decode(&expectAuthResponse)
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
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusInternalServerError))

				err = json.NewDecoder(resp.Body).Decode(&expectAuthResponse)
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
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(userBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusInternalServerError))

				err = json.NewDecoder(resp.Body).Decode(&expectAuthResponse)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectAuthResponse.Token).To(BeEmpty())
				Expect(expectAuthResponse.Errors).NotTo(BeEmpty())
			})
		})
	})
})
