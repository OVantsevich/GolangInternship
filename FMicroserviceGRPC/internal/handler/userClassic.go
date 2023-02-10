package handler

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"

	"GolangInternship/FMicroserviceGRPC/internal/model"
	"GolangInternship/FMicroserviceGRPC/internal/service"
	pr "GolangInternship/FMicroserviceGRPC/proto"

	"github.com/sirupsen/logrus"
)

// UserClassicService service interface for user handler
//
//go:generate mockery --name=UserClassicService --case=underscore --output=./mocks
type UserClassicService interface {
	Signup(ctx context.Context, user *model.User) (string, string, *model.User, error)
	Login(ctx context.Context, login, password string) (string, string, error)
	Refresh(ctx context.Context, login, userRefreshToken string) (string, string, error)
	Update(ctx context.Context, login string, user *model.User) error
	Delete(ctx context.Context, login string) error

	GetByLogin(ctx context.Context, login string) (*model.User, error)
}

// FileService service interface for user handler
//
//go:generate mockery --name=FileService --case=underscore --output=./mocks
type FileService interface {
	StoreFile(fileType string, fileData bytes.Buffer) (string, error)
}

// max Image Size
const maxImageSize = 1 << 25

// UserClassic handler
type UserClassic struct {
	pr.UnimplementedUserServiceServer
	s      UserClassicService
	file   FileService
	jwtKey string
}

// NewUserHandlerClassic new user handler
func NewUserHandlerClassic(s UserClassicService, file FileService, key string) *UserClassic {
	return &UserClassic{s: s, file: file, jwtKey: key}
}

// Signup handler signup
func (h *UserClassic) Signup(ctx context.Context, request *pr.SignupRequest) (response *pr.SignupResponse, err error) {
	user := &model.User{
		Login:    request.Login,
		Email:    request.Email,
		Password: request.Password,
		Name:     request.Name,
		Age:      int(request.Age),
	}

	var userResponse *model.User
	response = &pr.SignupResponse{}
	response.AccessToken, response.RefreshToken, userResponse, err = h.s.Signup(ctx, user)
	if err != nil {
		err = fmt.Errorf("userHandler - Signup - Signup: %w", err)
		logrus.Error(err)
		return
	}
	response.User = &pr.User{
		Login:    userResponse.Login,
		Email:    userResponse.Email,
		Password: userResponse.Password,
		Name:     userResponse.Name,
		Age:      int32(userResponse.Age),
		Role:     userResponse.Role,
	}

	return
}

// Login handler login
func (h *UserClassic) Login(ctx context.Context, request *pr.LoginRequest) (response *pr.LoginResponse, err error) {
	response = &pr.LoginResponse{}
	response.AccessToken, response.RefreshToken, err = h.s.Login(ctx, request.Login, request.Password)
	if err != nil {
		err = fmt.Errorf("userHandler - Login - Login: %w", err)
		logrus.Error(err)
		return
	}

	return
}

// Refresh handler refresh
func (h *UserClassic) Refresh(ctx context.Context, request *pr.RefreshRequest) (response *pr.RefreshResponse, err error) {
	response = &pr.RefreshResponse{}
	response.AccessToken, response.RefreshToken, err = h.s.Refresh(ctx, request.Login, request.RefreshToken)
	if err != nil {
		err = fmt.Errorf("userHandler - Refresh - Refresh: %w", err)
		logrus.Error(err)
		return
	}

	return
}

// Update handler update
func (h *UserClassic) Update(ctx context.Context, request *pr.UpdateRequest) (response *pr.UpdateResponse, err error) {
	var claims = ctx.Value("user").(*service.CustomClaims)

	user := &model.User{
		Email: request.Email,
		Name:  request.Name,
		Age:   int(request.Age),
	}
	response = &pr.UpdateResponse{}
	err = h.s.Update(ctx, claims.Login, user)
	if err != nil {
		err = fmt.Errorf("userHandler - Update - Update: %w", err)
		logrus.Error(err)
		return
	}
	response.Login = claims.Login

	return
}

// Delete handler delete
func (h *UserClassic) Delete(ctx context.Context, _ *pr.Request) (response *pr.DeleteResponse, err error) {
	var claims = ctx.Value("user").(*service.CustomClaims)

	response = &pr.DeleteResponse{}
	err = h.s.Delete(ctx, claims.Login)
	if err != nil {
		err = fmt.Errorf("userHandler - Delete - Delete: %w", err)
		logrus.Error(err)
		return
	}
	response.Login = claims.Login

	return
}

// UserByLogin handler user by login
func (h *UserClassic) UserByLogin(ctx context.Context, request *pr.UserByLoginRequest) (response *pr.UserByLoginResponse, err error) {
	var claims = ctx.Value("user").(*service.CustomClaims)

	if claims.Role != "admin" {
		err = fmt.Errorf("access denied")
		logrus.Error(err)
		return
	}

	response = &pr.UserByLoginResponse{}
	var user *model.User
	user, err = h.s.GetByLogin(ctx, request.Login)
	if err != nil {
		err = fmt.Errorf("userHandler - UserByLogin - GetByLogin: %w", err)
		logrus.Error(err)
		return
	}
	response.User = &pr.User{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Age:      int32(user.Age),
		Role:     user.Role,
	}

	return
}

// Upload handler upload
func (h *UserClassic) Upload(request pr.UserService_UploadServer) (err error) {
	req, err := request.Recv()
	if err != nil {
		err = status.Errorf(codes.Unknown, "cannot receive image info: %v", err)
		logrus.Error(err)
		return err
	}
	fileType := req.GetInfo().GetFileType()

	startTime := time.Now()
	logrus.Infof("start time:%s", startTime)

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		req, err = request.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			err = status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err)
			logrus.Error(err)
			return err
		}

		chunk := req.GetChunk()
		size := len(chunk)

		imageSize += size
		if imageSize > maxImageSize {
			err = status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize)
			logrus.Error(err)
			return err
		}
		_, err = imageData.Write(chunk)
		if err != nil {
			err = status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
			logrus.Error(err)
			return err
		}

		imageID, err := h.file.StoreFile(fileType, imageData)
		if err != nil {
			err = status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
			logrus.Error(err)
			return err
		}

		res := &pr.UploadResponse{
			Id: imageID,
		}

		err = request.SendAndClose(res)
		if err != nil {
			err = status.Errorf(codes.Unknown, "cannot send response: %v", err)
			logrus.Error(err)
			return err
		}

		return nil
	}
	return nil
}

// Download handler upload
func (h *UserClassic) Download(ctx context.Context, request *pr.DownloadRequest) (response *pr.DownloadResponse, err error) {

	return
}
