package model

import (
	larkpb "feishu/proto/go_proto"
)

type ArgocdBody struct {
	State     string `json:"state"`
	Context   string `json:"context"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	RepoUrl   string `json:"repo_url"`
	Revision  string `json:"revision"`
	Author    string `json:"author"`
	Message   string `json:"message"`
}

type Jarvis struct {
	ContainerName string `json:"container_name,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	Author        string `json:"author,omitempty"`
	Image         string `json:"image,omitempty"`
}

func PbJarvisToModel(request *larkpb.AppImageTagUpdateRequest) Jarvis {
	return Jarvis{
		ContainerName: request.ContainerName,
		CommitMessage: request.CommitMessage,
		Author:        request.Author,
		Image:         request.Image,
	}
}
