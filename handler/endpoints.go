package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = "somesecretkey123"

func (s *Server) Registration(ctx echo.Context) error {

	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Invalid JSON"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	var registrationRequest generated.RegistrationRequest
	err = json.Unmarshal(b, &registrationRequest)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Invalid JSON"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	err = validateRegistrationRequest(registrationRequest)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	hashedPass, err := hashingAndSalting(registrationRequest.Password)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Invalid password"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	output, err := s.Repository.InsertUsersData(ctx.Request().Context(), repository.InsertUsersDataInput{
		PhoneNumber:    registrationRequest.PhoneNumber,
		FullName:       registrationRequest.FullName,
		HashedPassword: hashedPass,
	})
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Internal error"
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	var resp generated.RegistrationResponse
	resp.Id = int(output.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func validateRegistrationRequest(request generated.RegistrationRequest) error {
	if len(request.PhoneNumber) < 10+2 {
		return errors.New("phone number must be at least 10 characters")
	}
	if len(request.PhoneNumber) > 13+2 {
		return errors.New("phone number must be at most 13 characters")
	}
	if len(request.FullName) < 3 {
		return errors.New("full name must be at least 3 characters")
	}
	if len(request.FullName) > 60 {
		return errors.New("full name must be at most 60 characters")
	}
	if len(request.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(request.Password) > 64 {
		return errors.New("password must be at most 64 characters")
	}
	var isContainsUppercase bool
	var isContainSpecialChar bool
	for _, v := range request.Password {
		if unicode.IsUpper(v) {
			isContainsUppercase = true
		} else if !unicode.IsLetter(v) {
			isContainSpecialChar = true
		}
	}
	if !isContainsUppercase {
		return errors.New("password must be at least has 1 capital")
	}
	if !isContainSpecialChar {
		return errors.New("password must be at least has 1 special character")
	}
	return nil
}

func hashingAndSalting(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func matchingPass(pass, hashedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))

	return err == nil
}

func (s *Server) Login(ctx echo.Context) error {

	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Invalid JSON"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	var loginRequest generated.LoginRequest
	err = json.Unmarshal(b, &loginRequest)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Invalid JSON"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	output, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), repository.GetUserByPhoneNumberInput{
		PhoneNumber: loginRequest.PhoneNumber,
	})
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Internal error"
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	if !matchingPass(loginRequest.Password, output.HashedPass) {
		var resp generated.ErrorResponse
		resp.Message = "Wrong combination phone number and password"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	token, err := generateJwtToken(int(output.Id), jwtSecret)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Internal error"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	err = s.Repository.UpdateCountById(ctx.Request().Context(), repository.UpdateCountByIdInput{
		Id:    output.Id,
		Count: output.Count + 1,
	})
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Internal error"
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	var resp generated.LoginResponse
	resp.Id = int(output.Id)
	resp.Jwt = token
	return ctx.JSON(http.StatusOK, resp)
}

type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func generateJwtToken(id int, secret string) (string, error) {
	claims := jwt.MapClaims{
		"id": id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(tokenString string, secret string) (int32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, okClaim := claims["id"].(float64)
		if !okClaim {
			return 0, errors.New("id not found")
		}
		return int32(id), nil
	}

	return 0, nil
}

func (s *Server) GetProfile(ctx echo.Context) error {

	req := ctx.Request()
	headers := req.Header

	token := headers.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)

	userId, err := validateToken(token, jwtSecret)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "not authorized"
		return ctx.JSON(http.StatusForbidden, resp)
	}

	output, err := s.Repository.GetUserById(ctx.Request().Context(), repository.GetUserByIdInput{
		Id: userId,
	})
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Internal error"
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	var resp generated.GetProfileResponse
	resp.FullName = output.FullName
	resp.PhoneNumber = output.PhoneNumber
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {

	req := ctx.Request()
	headers := req.Header

	token := headers.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)

	userId, err := validateToken(token, jwtSecret)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "not authorized"
		return ctx.JSON(http.StatusForbidden, resp)
	}

	b, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Invalid JSON"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	var updateProfileRequest generated.UpdateProfileRequest
	err = json.Unmarshal(b, &updateProfileRequest)
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Invalid JSON"
		return ctx.JSON(http.StatusBadRequest, resp)
	}

	if updateProfileRequest.PhoneNumber != nil {
		var resp generated.ErrorResponse
		resp.Message = "Cannot change phone number"
		return ctx.JSON(http.StatusConflict, resp)
	}

	if updateProfileRequest.FullName == nil {
		return ctx.JSON(http.StatusOK, nil)
	}

	err = s.Repository.UpdateFullNameByIdInput(ctx.Request().Context(), repository.UpdateFullNameByIdInput{
		Id:       userId,
		FullName: *updateProfileRequest.FullName,
	})
	if err != nil {
		var resp generated.ErrorResponse
		resp.Message = "Internal error"
		return ctx.JSON(http.StatusInternalServerError, resp)
	}

	return ctx.JSON(http.StatusOK, nil)
}
