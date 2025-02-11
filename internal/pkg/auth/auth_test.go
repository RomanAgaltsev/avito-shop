package auth_test

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/RomanAgaltsev/avito-shop/internal/model"
	"github.com/RomanAgaltsev/avito-shop/internal/pkg/auth"
)

var _ = Describe("Auth", func() {
	Describe("Creating new JWT token", func() {
		var (
			ja        *jwtauth.JWTAuth
			secretKey string
			userName  string
		)

		JustBeforeEach(func() {
			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())
		})

		Context("When the secret key and user name are defined and correct", func() {
			BeforeEach(func() {
				secretKey = "secret"
				userName = "user"
			})

			It("can create new valid JWT token", func() {
				token, tokenString, err := auth.NewJWTToken(ja, userName)
				Expect(err).NotTo(HaveOccurred())
				Expect(tokenString).NotTo(BeEmpty())

				err = jwt.Validate(token, ja.ValidateOptions()...)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When the secret key and user name are undefined", func() {
			BeforeEach(func() {
				secretKey = ""
				userName = ""
			})

			It("cannot create new JWT token", func() {
				token, tokenString, err := auth.NewJWTToken(ja, userName)
				Expect(err).To(HaveOccurred())
				Expect(tokenString).To(BeEmpty())

				Expect(func() {
					err = jwt.Validate(token, ja.ValidateOptions()...)
				}).To(Panic())
			})
		})
	})

	Describe("Hashing password", func() {
		var (
			password string
			hash     string
			err      error
		)

		Context("When the password is defined", func() {
			BeforeEach(func() {
				password = "password"
			})

			It("can hash the password", func() {
				hash, err = auth.HashPassword(password)
				Expect(err).NotTo(HaveOccurred())
				Expect(hash).NotTo(BeEmpty())
			})
			It("can check the password hash", func() {
				result := auth.CheckPasswordHash(password, hash)
				Expect(result).To(BeTrue())
			})
		})

		Context("When the password is undefined", func() {
			BeforeEach(func() {
				password = ""
			})

			It("can hash the password", func() {
				hash, err = auth.HashPassword(password)
				Expect(err).NotTo(HaveOccurred())
				Expect(hash).NotTo(BeEmpty())
			})
			It("can check the password hash", func() {
				result := auth.CheckPasswordHash(password, hash)
				Expect(result).To(BeTrue())
			})
		})

		Context("When the password is longer than 72 bytes", func() {
			BeforeEach(func() {
				password = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFJHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFJHIJKLMNOPQRSTUVWXYZ"

			})

			It("cannot hash the password", func() {
				hash, err = auth.HashPassword(password)
				Expect(err).To(HaveOccurred())
				Expect(hash).To(BeEmpty())
			})
		})
	})

	Describe("Extracting user login from HTTP request", func() {
		const secretKeyEnc = "secret"

		var (
			ja          *jwtauth.JWTAuth
			request     *http.Request
			user        model.User
			emptyUser   model.User
			userName    string
			secretKey   string
			tokenString string
			err         error
		)

		JustBeforeEach(func() {
			ja = auth.NewAuth(secretKeyEnc)
			Expect(ja).ShouldNot(BeNil())

			_, tokenString, err = auth.NewJWTToken(ja, userName)
			Expect(err).NotTo(HaveOccurred())
			Expect(tokenString).NotTo(BeEmpty())

			request, err = http.NewRequest("", "", nil)
			Expect(err).NotTo(HaveOccurred())

			bearer := "Bearer " + tokenString
			request.Header.Add("Authorization", bearer)
		})

		Context("When the secret key is defined and right", func() {
			BeforeEach(func() {
				secretKey = secretKeyEnc
				userName = "user"
			})

			It("can get user from request", func() {
				user, err = auth.UserFromRequest(request, secretKey)
				Expect(err).NotTo(HaveOccurred())
				Expect(user.UserName).To(Equal(userName))
			})
		})

		Context("When the secret key is defined and wrong", func() {
			BeforeEach(func() {
				secretKey = "wrong key"
				userName = "user"

				emptyUser = model.User{}
			})

			It("cannot get user from request", func() {
				user, err = auth.UserFromRequest(request, secretKey)
				Expect(err).To(HaveOccurred())
				Expect(user).To(Equal(emptyUser))
			})
		})

		Context("When the secret key is undefined", func() {
			BeforeEach(func() {
				secretKey = ""
				userName = "user"

				emptyUser = model.User{}
			})

			It("cannot get user from request", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user).To(Equal(emptyUser))
			})
		})
	})
})
