package web

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/middlewares"
	core "ssi-gitlab.teda.th/ssi/core"
)

func NewHomeHTTPHandler(r *echo.Echo) {
	home := &HomeController{}
	auth := &AuthController{}
	user := &UserController{}
	me := &MeController{}
	//org := &OrganizationController{}
	vc := &VCController{}
	schema := &SchemaController{}
	config := &ConfigController{}
	mobileUser := &MobileUserController{}
	mobileGroup := &MobileGroupController{}
	key := &KeyController{}
	vp := &VPController{}
	r.GET("/", core.WithHTTPContext(home.Get))

	r.GET("/web/configs/wallets", core.WithHTTPContext(config.WalletGet), middlewares.IsUser)
	r.POST("/web/configs/wallets", core.WithHTTPContext(config.WalletSetting), middlewares.IsUser)
	r.DELETE("/web/configs/wallets/:id", core.WithHTTPContext(config.WalletDelete), middlewares.IsUser)

	r.GET("/web/configs/schemas", core.WithHTTPContext(config.SchemaRepositoryPagination), middlewares.IsUser)
	r.POST("/web/configs/schemas", core.WithHTTPContext(config.SchemaRepositoryCreate), middlewares.IsUser)
	r.GET("/web/configs/schemas/:id", core.WithHTTPContext(config.SchemaRepositoryFind), middlewares.IsUser)
	r.PUT("/web/configs/schemas/:id", core.WithHTTPContext(config.SchemaRepositoryUpdate), middlewares.IsUser)
	r.DELETE("/web/configs/schemas/:id", core.WithHTTPContext(config.SchemaRepositoryDelete), middlewares.IsUser)

	r.GET("/web/schemas/:repository_id", core.WithHTTPContext(schema.Pagination), middlewares.IsUser)
	r.POST("/web/schemas/:repository_id", core.WithHTTPContext(schema.Create), middlewares.IsUser)
	r.GET("/web/schemas/:repository_id/types", core.WithHTTPContext(schema.Types), middlewares.IsUser)
	r.GET("/web/schemas/:repository_id/tokens", core.WithHTTPContext(schema.TokenPagination), middlewares.IsUser)
	r.POST("/web/schemas/:repository_id/tokens", core.WithHTTPContext(schema.TokenCreate), middlewares.IsUser)
	r.GET("/web/schemas/:repository_id/tokens/:token_id", core.WithHTTPContext(schema.TokenFind), middlewares.IsUser)
	r.PUT("/web/schemas/:repository_id/tokens/:token_id", core.WithHTTPContext(schema.TokenUpdate), middlewares.IsUser)
	r.DELETE("/web/schemas/:repository_id/tokens/:token_id", core.WithHTTPContext(schema.TokenDelete), middlewares.IsUser)
	r.POST("/web/schemas/:repository_id/tokens", core.WithHTTPContext(schema.TokenCreate), middlewares.IsUser)
	r.POST("/web/schemas/:repository_id/upload", core.WithHTTPContext(schema.CreateByUpload), middlewares.IsUser)
	r.PUT("/web/schemas/:repository_id/:schema_id", core.WithHTTPContext(schema.Update), middlewares.IsUser)
	r.POST("/web/schemas/:repository_id/:schema_id/upload", core.WithHTTPContext(schema.UploadByUpload), middlewares.IsUser)

	r.GET("/web/schemas/:repository_id/:schema_id", core.WithHTTPContext(schema.Find), middlewares.IsUser)
	r.PUT("/web/schemas/:repository_id/:schema_id", core.WithHTTPContext(schema.Update), middlewares.IsUser)
	r.GET("/web/schemas/:repository_id/:schema_id/history", core.WithHTTPContext(schema.FindHistory), middlewares.IsUser)
	r.GET("/web/schemas/:repository_id/:schema_id/:version", core.WithHTTPContext(schema.FindByVersion), middlewares.IsUser)
	r.GET("/web/schemas/:repository_id/:schema_id/:version/schema", core.WithHTTPContext(schema.FindSchemaInstance), middlewares.IsUser)
	r.GET("/web/schemas/:repository_id/:schema_id/:version/:reference", core.WithHTTPContext(schema.FindSchemaReference), middlewares.IsUser)

	r.POST("/web/login", core.WithHTTPContext(auth.Login))
	r.POST("/web/logout", core.WithHTTPContext(auth.Logout), middlewares.IsUser)

	r.POST("/web/key/generate", core.WithHTTPContext(key.Generate), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/key/generate/rsa", core.WithHTTPContext(key.GenerateRSA), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/key/upload", core.WithHTTPContext(key.Upload), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/key/upload/x509", core.WithHTTPContext(key.UploadX509), middlewares.IsUser, middlewares.IsOrganizationMember)

	r.GET("/web/me", core.WithHTTPContext(me.Profile), middlewares.IsUser)
	r.POST("/web/me/change-password", core.WithHTTPContext(me.ChangePassword), middlewares.IsUser)

	r.GET("/web/mobile/groups", core.WithHTTPContext(mobileGroup.Get), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/mobile/groups", core.WithHTTPContext(mobileGroup.GroupCreate), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/mobile/groups/users", core.WithHTTPContext(mobileGroup.AddGroupUser), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.PUT("/web/mobile/groups/:id", core.WithHTTPContext(mobileGroup.GroupUpdate), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.DELETE("/web/mobile/groups/:id", core.WithHTTPContext(mobileGroup.GroupDelete), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.GET("/web/mobile/groups/:id/users", core.WithHTTPContext(mobileGroup.Pagination), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/mobile/groups/:id/users/remove", core.WithHTTPContext(mobileGroup.RemoveGroupUser), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.GET("/web/mobile/users", core.WithHTTPContext(mobileUser.Pagination), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.GET("/web/mobile/users/:id", core.WithHTTPContext(mobileUser.Find), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.PUT("/web/mobile/users/:id", core.WithHTTPContext(mobileUser.Update), middlewares.IsUser, middlewares.IsOrganizationMember)

	r.POST("/web/vcs", core.WithHTTPContext(vc.Create), middlewares.IsUser)
	r.GET("/web/vcs", core.WithHTTPContext(vc.Pagination), middlewares.IsUser)
	r.GET("/web/vcs/:id", core.WithHTTPContext(vc.Find), middlewares.IsUser)
	r.GET("/web/vcs/did/:did", core.WithHTTPContext(vc.FindByDID), middlewares.IsUser)
	r.PUT("/web/vcs/:id", core.WithHTTPContext(vc.Update), middlewares.IsUser)
	r.POST("/web/vcs/:id/approve", core.WithHTTPContext(vc.Approve), middlewares.ValidateMessageMiddleware)
	r.POST("/web/vcs/:id/reject", core.WithHTTPContext(vc.Reject), middlewares.ValidateMessageMiddleware)
	r.POST("/web/vcs/:id/revoke", core.WithHTTPContext(vc.Revoke), middlewares.IsUser)
	r.POST("/web/vcs/qr", core.WithHTTPContext(vc.CreateVCQR), middlewares.IsUser)
	r.GET("/web/vcs/qr/:token_id", core.WithHTTPContext(vc.VerifyVCQR), middlewares.IsQRVerify, middlewares.ValidateSignatureMiddleware)
	r.POST("/web/vcs/:id/signing", core.WithHTTPContext(vc.VerifyVCQR), middlewares.IsQRVerify, middlewares.ValidateSignatureMiddleware)
	r.POST("/web/vcs/verify", core.WithHTTPContext(vc.VerifyVC), middlewares.IsUser)

	r.GET("/web/requested-vps", core.WithHTTPContext(vp.RequestedPagination), middlewares.IsUser)
	r.POST("/web/requested-vps", core.WithHTTPContext(vp.RequestedCreate), middlewares.IsUser)
	r.GET("/web/requested-vps/qr/:id", core.WithHTTPContext(vp.RequestedFindQR))
	r.GET("/web/requested-vps/:id", core.WithHTTPContext(vp.RequestedFind), middlewares.IsUser)
	r.PUT("/web/requested-vps/:id", core.WithHTTPContext(vp.RequestedUpdate), middlewares.IsUser)
	r.POST("/web/requested-vps/cancel", core.WithHTTPContext(vp.RequestUpdateList), middlewares.IsUser)
	r.POST("/web/requested-vps/:id/submit", core.WithHTTPContext(vp.SubmitVP), middlewares.ValidateJWTMessageMiddleware)
	r.GET("/web/requested-vps/:id/submitted-vps", core.WithHTTPContext(vp.SubmittedPagnination), middlewares.IsUser)
	r.POST("/web/requested-vps/:id/update-qrcode", core.WithHTTPContext(vp.RequestedUpdateQRCode), middlewares.IsUser)

	r.GET("/web/submitted-vps/:id/vcs", core.WithHTTPContext(vp.SubmittedVCList), middlewares.IsUser)
	r.GET("/web/submitted-vps/:id/vcs/:vc_id", core.WithHTTPContext(vp.SubmittedVCFind), middlewares.IsUser)
	r.GET("/web/submitted-vps/:id", core.WithHTTPContext(vp.SubmittedFind), middlewares.IsUser)
	r.POST("/web/submitted-vps/:id/tag", core.WithHTTPContext(vp.TagStatus), middlewares.IsUser)

	r.POST("/web/vps/verify", core.WithHTTPContext(vc.VerifyVP), middlewares.IsUser)

	r.GET("/web/users", core.WithHTTPContext(user.Pagination), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/users", core.WithHTTPContext(user.Register), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.GET("/web/users/:user_id", core.WithHTTPContext(user.Find), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.PUT("/web/users/:user_id", core.WithHTTPContext(user.Update), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.DELETE("/web/users/:user_id", core.WithHTTPContext(user.Delete), middlewares.IsUser, middlewares.IsOrganizationMember)
	r.POST("/web/users/:user_id/reset-password", core.WithHTTPContext(user.ResetPassword), middlewares.IsUser)
	r.GET("/web/user-verify", core.WithHTTPContext(user.CheckVerifyToken))
	r.POST("/web/user-verify", core.WithHTTPContext(user.Verify))

}
