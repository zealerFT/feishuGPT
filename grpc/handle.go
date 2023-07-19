package grpc

import (
	"context"

	"feishu/model"

	larkpb "feishu/proto/go_proto"
)

func (s *Servlet) AppImageTagUpdate(ctx context.Context, request *larkpb.AppImageTagUpdateRequest) (*larkpb.AppImageTagUpdateResponse, error) {
	resp := &larkpb.AppImageTagUpdateResponse{}
	err := s.Dep.Lark.SendJarvisMsg(ctx, model.PbJarvisToModel(request), s.Dep.Cfg.LarkSetting.ArgocdChatId, "chat_id")
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *Servlet) ImagesSyncDone(ctx context.Context, request *larkpb.ImagesSyncDoneRequest) (*larkpb.ImagesSyncDoneResponse, error) {
	resp := &larkpb.ImagesSyncDoneResponse{}
	err := s.Dep.Lark.SendImagesSyncMsg(ctx, request.Image, s.Dep.Cfg.LarkSetting.ArgocdChatId, "chat_id")
	if err != nil {
		return resp, err
	}
	return resp, nil
}
