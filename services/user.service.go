package services

import (
	"errors"
	"fmt"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type UserLoginPayload struct {
	Email    string
	Password string
	IP       string
	Device   string
}

type UserUpdatePayload struct {
	Role string `json:"role"`
}

type UserRegisterPayload struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DateOfBirth    string `json:"date_of_birth"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id"`
}

type IUserService interface {
	Pagination(pageOptions *core.PageOptions) ([]models.User, *core.PageResponse, core.IError)
	Login(payload *UserLoginPayload) (*models.User, *models.UserAccessToken, core.IError)
	Logout(token string) core.IError
	Update(id string, payload *UserUpdatePayload) (*models.User, core.IError)
	Verify(verifyToken string, password string) (*models.User, core.IError)
	Register(payload *UserRegisterPayload) (*models.User, core.IError)
	Find(id string) (*models.User, core.IError)
	ChangePassword(payload *UserPasswordChangePayload) core.IError
	Delete(id string) core.IError
	CheckVerifyToken(token string) core.IError
	ResetPassword(id string) core.IError
}
type userService struct {
	ctx        core.IContext
	orgService IOrganizationService
}

func NewUserService(ctx core.IContext) IUserService {
	return &userService{ctx: ctx, orgService: NewOrganizationService(ctx)}
}

func (s userService) Login(payload *UserLoginPayload) (*models.User, *models.UserAccessToken, core.IError) {
	user := &models.User{}
	err := s.ctx.DB().First(user, "email = ? AND status = ?", payload.Email, consts.UserStatusActive).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, s.ctx.NewError(err, emsgs.AuthEmailOrPasswordInvalid, user)
	}
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	isValidPassword := helpers.ComparePassword(utils.GetString(user.Password), payload.Password)
	if !isValidPassword {
		return nil, nil, s.ctx.NewError(err, emsgs.AuthEmailOrPasswordInvalid, user)
	}

	userToken := models.NewUserAccessToken(user.ID)
	userToken.Info, err = helpers.SetJSONValue(userToken.Info, "ip", payload.IP)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.InternalServerError)
	}
	userToken.Info, err = helpers.SetJSONValue(userToken.Info, "device", payload.Device)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.InternalServerError)
	}

	err = s.ctx.DB().Create(userToken).Error
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError, user, userToken)
	}

	return user, userToken, nil
}

func (s userService) Update(id string, payload *UserUpdatePayload) (*models.User, core.IError) {
	user, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if payload.Role != "" {
		user.Role = payload.Role
	}

	err := s.ctx.DB().Updates(user).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(id)
}

func (s userService) Register(payload *UserRegisterPayload) (*models.User, core.IError) {
	user := models.NewUser()
	user.Email = payload.Email
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Role = payload.Role
	user.DateOfBirth = payload.DateOfBirth
	user.OrganizationID = payload.OrganizationID
	err := s.ctx.DB().Create(user).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	err = helpers.SetPasswordEmail(
		s.ctx.ENV().String(consts.ENVSMTPHost),
		s.ctx.ENV().String(consts.ENVSMTPPort),
		s.ctx.ENV().String(consts.ENVSenderEmail),
		payload.Email,
		fmt.Sprintf("%s/verify/%s", s.ctx.ENV().String(consts.ENVWebPortalBaseURL), utils.GetString(user.VerifyToken)),
		s.ctx.ENV().String(consts.ENVWebPortalBaseURL))
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(user.ID)
}

func (s userService) Pagination(pageOptions *core.PageOptions) ([]models.User, *core.PageResponse, core.IError) {
	users := make([]models.User, 0)

	db := s.ctx.DB()
	if len(pageOptions.OrderBy) == 0 {
		pageOptions.OrderBy = []string{"created_at DESC"}

	}
	db = core.SetSearchSimple(db, pageOptions.Q, []string{"email", "first_name", "last_name", "role"})
	pageResponse, err := core.Paginate(db, &users, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return users, pageResponse, nil
}

func (s userService) Verify(verifyToken string, password string) (*models.User, core.IError) {
	user := &models.User{}
	err := s.ctx.DB().First(user, "verify_token = ?", verifyToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("token"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	user.Status = consts.UserStatusActive
	token := ""
	user.VerifyToken = &token
	user.SetPassword(password)
	user.UpdatedAt = utils.GetCurrentDateTime()
	err = s.ctx.DB().Updates(user).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(user.ID)
}

func (s userService) Find(id string) (*models.User, core.IError) {
	user := &models.User{}
	err := s.ctx.DB().
		Preload("Organization").
		First(user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("user"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return user, nil
}

func (s userService) Logout(token string) core.IError {
	err := s.ctx.DB().Delete(&models.UserAccessToken{}, "token = ?", token).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

type UserPasswordChangePayload struct {
	ID              string
	CurrentPassword string
	NewPassword     string
}

func (s userService) ChangePassword(payload *UserPasswordChangePayload) core.IError {
	user, ierr := s.Find(payload.ID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	isValidPassword := helpers.ComparePassword(utils.GetString(user.Password), payload.CurrentPassword)
	if !isValidPassword {
		return s.ctx.NewError(ierr, emsgs.AuthCurrentPasswordInvalid, user)
	}

	user.SetPassword(payload.NewPassword)
	user.UpdatedAt = utils.GetCurrentDateTime()

	err := s.ctx.DB().Updates(user).Error
	if err != nil {
		return s.ctx.NewError(ierr, errmsgs.DBError)
	}

	err = s.ctx.DB().Delete(&models.UserAccessToken{}, "user_id = ?", payload.ID).Error
	if err != nil {
		return s.ctx.NewError(ierr, errmsgs.DBError)
	}

	return nil
}

func (s userService) ResetPassword(id string) core.IError {
	user, ierr := s.Find(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	token := utils.NewSha256(user.ID + utils.GetCurrentDateTime().String() + utils.GetUUID() + consts.UserStatusActive)
	user.VerifyToken = &token
	err := s.ctx.DB().Save(user).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	err = helpers.SendResetPasswordEmail(
		s.ctx.ENV().String(consts.ENVSMTPHost),
		s.ctx.ENV().String(consts.ENVSMTPPort),
		s.ctx.ENV().String(consts.ENVSenderEmail),
		user.Email,
		fmt.Sprintf("%s/verify/%s", s.ctx.ENV().String(consts.ENVWebPortalBaseURL), utils.GetString(user.VerifyToken)),
		s.ctx.ENV().String(consts.ENVWebPortalBaseURL))
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

func (s userService) Delete(id string) core.IError {
	user, ierr := s.Find(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Delete(user).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	return nil
}

func (s userService) CheckVerifyToken(token string) core.IError {
	if token == "" {
		return s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFoundCustomError("token"))
	}

	user := &models.User{}
	err := s.ctx.DB().First(user, "verify_token = ?", token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.ctx.NewError(err, errmsgs.NotFoundCustomError("token"))
	}

	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}
