package handler

import (
	"GolangInternship/FMicroservice/internal/handler/mocks"
	"GolangInternship/FMicroservice/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUser_Signup(t *testing.T) {
	testInit()
	s := mocks.NewUserService(t)
	h := NewUserHandlerClassic(s)

	for _, user := range testValidData {
		s.On("Signup", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("*model.User")).
			Return("", "", &user, nil).
			Once()
		data, err := json.Marshal(user)
		require.NoError(t, err)
		reader := bytes.NewReader(data)

		req := httptest.NewRequest(http.MethodPost, "/", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		require.NoError(t, h.Signup(c))
		require.Equal(t, http.StatusCreated, rec.Code)
		//check, _ := json.Marshal(SignupResponse{&user, &TokenResponse{
		//	AccessToken:  "",
		//	RefreshToken: "",
		//}})
		//require.Equal(t, check, rec.Body.Bytes())
	}

	for _, user := range testNoValidData {
		data, err := json.Marshal(user)
		require.NoError(t, err)
		reader := bytes.NewReader(data)

		req := httptest.NewRequest(http.MethodPost, "/", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = h.Signup(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	}

	s.On("Signup", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("*model.User")).
		Return("", "", nil, fmt.Errorf("something went wrong")).
		Once()

	data, err := json.Marshal(testValidData[0])
	require.NoError(t, err)
	reader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = h.Signup(c)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestUserClassic_Login(t *testing.T) {
	testInit()
	s := mocks.NewUserService(t)
	h := NewUserHandlerClassic(s)

	for _, user := range testValidData {
		s.On("Login", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return("", "", nil).
			Once()
		data, err := json.Marshal(user)
		require.NoError(t, err)
		reader := bytes.NewReader(data)

		req := httptest.NewRequest(http.MethodPost, "/", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusOK, rec.Code)
	}

	s.On("Login", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", "", fmt.Errorf("something went wrong")).
		Once()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Login(c)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestUserClassic_Refresh(t *testing.T) {
	testInit()
	s := mocks.NewUserService(t)
	h := NewUserHandlerClassic(s)

	for _, user := range testValidData {
		s.On("Refresh", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"), mock.AnythingOfType("string")).
			Return("", "", nil).
			Once()
		data, err := json.Marshal(user)
		require.NoError(t, err)
		reader := bytes.NewReader(data)

		req := httptest.NewRequest(http.MethodPost, "/", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256,
			&service.CustomClaims{
				Login: user.Login,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
				}}))

		require.NoError(t, h.Refresh(c))
		require.Equal(t, http.StatusOK, rec.Code)
	}

	s.On("Refresh", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return("", "", fmt.Errorf("something went wrong")).
		Once()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256,
		&service.CustomClaims{
			Login: testValidData[0].Login,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			}}))

	err := h.Refresh(c)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestUserClassic_Update(t *testing.T) {
	testInit()
	s := mocks.NewUserService(t)
	h := NewUserHandlerClassic(s)

	for _, user := range testValidData {
		s.On("Update", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).
			Return(nil).
			Once()
		data, err := json.Marshal(user)
		require.NoError(t, err)
		reader := bytes.NewReader(data)

		req := httptest.NewRequest(http.MethodPost, "/", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256,
			&service.CustomClaims{
				Login: testValidData[0].Login,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
				}}))

		require.NoError(t, h.Update(c))
		require.Equal(t, http.StatusOK, rec.Code)
	}

	for _, user := range testNoValidData {
		data, err := json.Marshal(user)
		require.NoError(t, err)
		reader := bytes.NewReader(data)

		req := httptest.NewRequest(http.MethodPost, "/", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256,
			&service.CustomClaims{
				Login: testValidData[0].Login,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
				}}))

		err = h.Update(c)
		require.Error(t, err)
		require.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	}

	s.On("Update", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).
		Return(fmt.Errorf("something went wrong")).
		Once()

	data, err := json.Marshal(testValidData[0])
	require.NoError(t, err)
	reader := bytes.NewReader(data)

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256,
		&service.CustomClaims{
			Login: testValidData[0].Login,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			}}))

	err = h.Update(c)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestUserClassic_Delete(t *testing.T) {
	testInit()
	s := mocks.NewUserService(t)
	h := NewUserHandlerClassic(s)

	for _, user := range testValidData {
		s.On("Delete", mock.AnythingOfType("*context.emptyCtx"),
			mock.AnythingOfType("string")).
			Return(nil).
			Once()
		data, err := json.Marshal(user)
		require.NoError(t, err)
		reader := bytes.NewReader(data)

		req := httptest.NewRequest(http.MethodPost, "/", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256,
			&service.CustomClaims{
				Login: user.Login,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
				}}))

		require.NoError(t, h.Delete(c))
		require.Equal(t, http.StatusOK, rec.Code)
	}

	s.On("Delete", mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("string")).
		Return(fmt.Errorf("something went wrong")).
		Once()

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256,
		&service.CustomClaims{
			Login: testValidData[0].Login,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			}}))

	err := h.Delete(c)
	require.Error(t, err)
	require.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}
