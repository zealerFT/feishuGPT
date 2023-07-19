package lark

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkext "github.com/larksuite/oapi-sdk-go/v3/service/ext"
)

func (l *Larkbiz) GetAppAccessTokenBySelfBuiltApp(ctx context.Context) {
	var resp, err = l.Lark.GetAppAccessTokenBySelfBuiltApp(ctx, &larkcore.SelfBuiltAppAccessTokenReq{
		AppID:     l.Cfg.LarkSetting.AppID,
		AppSecret: l.Cfg.LarkSetting.AppSecret,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func (l *Larkbiz) GetAppAccessTokenByMarketApp(ctx context.Context, appTicket string) {
	client := lark.NewClient(l.Cfg.LarkSetting.AppID, l.Cfg.LarkSetting.AppSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug), lark.WithOpenBaseUrl("https://open.larksuite-boe.com"))

	var resp, err = client.GetAppAccessTokenByMarketplaceApp(ctx, &larkcore.MarketplaceAppAccessTokenReq{
		AppID:     l.Cfg.LarkSetting.AppID,
		AppSecret: l.Cfg.LarkSetting.AppSecret,
		AppTicket: appTicket,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func (l *Larkbiz) GetTenantAccessTokenBySelfBuiltApp(ctx context.Context) {
	var resp, err = l.Lark.GetTenantAccessTokenBySelfBuiltApp(ctx, &larkcore.SelfBuiltTenantAccessTokenReq{
		AppID:     l.Cfg.LarkSetting.AppID,
		AppSecret: l.Cfg.LarkSetting.AppSecret,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func (l *Larkbiz) GetTenantAccessTokenByMarketApp(ctx context.Context, appAccessToken, tenantKey string) {
	var resp, err = l.Lark.GetTenantAccessTokenByMarketplaceApp(ctx, &larkcore.MarketplaceTenantAccessTokenReq{
		AppAccessToken: appAccessToken,
		TenantKey:      tenantKey,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func (l *Larkbiz) ResendAppTicket(ctx context.Context) {
	var resp, err = l.Lark.ResendAppTicket(ctx, &larkcore.ResendAppTicketReq{
		AppID:     l.Cfg.LarkSetting.AppID,
		AppSecret: l.Cfg.LarkSetting.AppSecret,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func (l *Larkbiz) GetAuthenAccessToken(ctx context.Context, code string) {
	var resp, err = l.Lark.Ext.Authen.AuthenAccessToken(ctx,
		larkext.NewAuthenAccessTokenReqBuilder().
			Body(larkext.NewAuthenAccessTokenReqBodyBuilder().
				GrantType(larkext.GrantTypeAuthorizationCode).
				Code(code).
				Build()).
			Build())
	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp.Data))
}

func (l *Larkbiz) RefreshAuthenAccessToken(ctx context.Context, refreshToken string) {
	var resp, err = l.Lark.Ext.Authen.RefreshAuthenAccessToken(ctx,
		larkext.NewRefreshAuthenAccessTokenReqBuilder().
			Body(larkext.NewRefreshAuthenAccessTokenReqBodyBuilder().
				GrantType(larkext.GrantTypeRefreshCode).
				RefreshToken(refreshToken).
				Build()).
			Build())
	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(resp.Data.RefreshToken)

	fmt.Println(larkcore.Prettify(resp))
}

func (l *Larkbiz) AuthenUserInfo(ctx context.Context, userAccessToken string) {
	var resp, err = l.Lark.Ext.Authen.AuthenUserInfo(ctx, larkcore.WithUserAccessToken(userAccessToken))
	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}
