package lark

import (
	"context"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

func (l *Larkbiz) GetUserInfo(ctx context.Context, userId, userIdType string) {
	if userIdType == "" {
		userIdType = larkcontact.UserIdTypeOpenId
	}
	resp, err := l.Lark.Contact.User.Get(ctx, larkcontact.NewGetUserReqBuilder().
		UserIdType(userIdType).
		UserId(userId).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(resp.Data.User)
	} else {
		fmt.Println(resp.Msg, resp.Code, resp.RequestId())
	}

}

func (l *Larkbiz) PatchUser(ctx context.Context, userId, userIdType string) {
	if userIdType == "" {
		userIdType = larkcontact.UserIdTypeOpenId
	}
	user := larkcontact.NewUserBuilder().Build()
	resp, err := l.Lark.Contact.User.Patch(ctx,
		larkcontact.NewPatchUserReqBuilder().
			UserId(userId).
			UserIdType(userIdType).
			User(user).
			Build(), larkcore.WithUserAccessToken("ssss"))

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(resp.Data.User)
	} else {
		fmt.Println(resp.Msg, resp.Code, resp.RequestId())
	}
}

func (l *Larkbiz) CreateUser(ctx context.Context, userIdType string) {
	if userIdType == "" {
		userIdType = larkcontact.UserIdTypeOpenId
	}
	resp, err := l.Lark.Contact.User.Create(ctx,
		larkcontact.NewCreateUserReqBuilder().UserIdType(userIdType).User(larkcontact.NewUserBuilder().Build()).Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(resp.Data.User)
	} else {
		fmt.Println(resp.Msg, resp.Code, resp.RequestId())
	}
}

func (l *Larkbiz) BatchGetId(ctx context.Context, userIdType string, phones []string) {
	if userIdType == "" {
		userIdType = larkcontact.UserIdTypeOpenId
	}
	resp, err := l.Lark.Contact.User.BatchGetId(ctx,
		larkcontact.NewBatchGetIdUserReqBuilder().
			UserIdType(userIdType).
			Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
				Mobiles(phones).
				Build()).
			Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(resp.Data)
	} else {
		fmt.Println(resp.Msg, resp.Code, resp.RequestId())
	}
}
