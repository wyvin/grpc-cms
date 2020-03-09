package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-content/pkg/handler/articleDb"
	"grpc-content/pkg/util"
	pb "grpc-content/proto"
)

type articleService struct {
	*pb.UnimplementedArticleServer
}

func NewArticleService() *articleService {
	return &articleService{}
}

func CheckAppIdAndGroupId(appId string, groupId uint32) error {
	if appId == "" {
		return status.Error(codes.InvalidArgument, "缺少app_id参数")
	}
	if groupId == 0 {
		return status.Error(codes.InvalidArgument, "缺少group_id参数")
	}
	return nil
}

// 文章
func (a articleService) MigrateArticle(ctx context.Context, req *pb.MigrateArticleRequest) (*pb.MigrateArticleResponse, error) {
	articleInstance := articleDb.Article{}
	err := articleInstance.Migrate()
	if err != nil {
		return nil, err
	}
	return &pb.MigrateArticleResponse{
		Message: "success",
	}, nil
}

func (a articleService) AddArticle(ctx context.Context, req *pb.AddArticleRequest) (*pb.AddArticleResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}

	articleInstance := articleDb.Article{
		AppId:          req.AppId,
		GroupId:        req.GroupId,
		TopicId:        req.TopicId,
		Cover:          req.Cover,
		Title:          req.Title,
		Author:         req.Author,
		Source:         req.Source,
		Recommendation: req.Recommendation,
		Content:        req.Content,
		State:          req.State,
	}

	articleId, err := articleInstance.Insert()
	if err != nil {
		return nil, err
	}

	return &pb.AddArticleResponse{
		ArticleId: articleId,
	}, nil
}

func (a articleService) GetArticleList(ctx context.Context, req *pb.GetArticleListRequest) (*pb.GetArticleListResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}

	util.FormatPageAndPerPage(&req.Page, &req.PerPage)

	articleInstance := articleDb.Article{
		AppId:   req.AppId,
		GroupId: req.GroupId,
		Page:    req.Page,
		PerPage: req.PerPage,
		Source:  req.Source,
		Title:   req.Title,
		State:   req.State,
	}
	record, total, err := articleInstance.GetArticleList()
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ArticleListElement, len(record))
	for i := 0; i < len(record); i++ {
		ele := record[i]
		list[i] = &pb.ArticleListElement{
			ArticleId:      uint32(ele.ID),
			TopicId:        ele.TopicId,
			Cover:          ele.Cover,
			Title:          ele.Title,
			Author:           ele.Author,
			Source:         ele.Source,
			Recommendation: ele.Recommendation,
			State:          ele.State,
			CreateAt:       util.FormatDateTime(ele.CreatedAt),
		}
	}

	return &pb.GetArticleListResponse{
		List:  list,
		Total: total,
	}, nil
}

// 专题
func (a articleService) MigrateArticleTopic(ctx context.Context, req *pb.MigrateArticleTopicRequest) (*pb.MigrateArticleTopicResponse, error) {
	topicInstance := articleDb.ArticleTopic{}
	err := topicInstance.Migrate()
	if err != nil {
		return nil, err
	}
	return &pb.MigrateArticleTopicResponse{
		Message: "success",
	}, nil
}

func (a articleService) AddArticleTopic(ctx context.Context, req *pb.AddArticleTopicRequest) (*pb.AddArticleTopicResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}

	articleInstance := articleDb.ArticleTopic{
		AppId:          req.AppId,
		GroupId:        req.GroupId,
		Name:          req.Name,
		Sort:         req.Sort,
		State:          req.State,
	}
	var err error
	exists, err := articleInstance.ExistByName()
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, status.Error(codes.AlreadyExists, "专题名已存在")
	}

	topicId, err := articleInstance.Insert()
	if err != nil {
		return nil, err
	}

	return &pb.AddArticleTopicResponse{
		TopicId: topicId,
	}, nil
}

func (a articleService) GetArticleTopicList(ctx context.Context, req *pb.GetArticleTopicListRequest) (*pb.GetArticleTopicListResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}

	util.FormatPageAndPerPage(&req.Page, &req.PerPage)

	topicInstance := articleDb.ArticleTopic{
		AppId:   req.AppId,
		GroupId: req.GroupId,
		Page:    req.Page,
		PerPage: req.PerPage,
		State:   req.State,
	}
	record, total, err := topicInstance.GetTopicList()
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ArticleTopicListElement, len(record))
	for i := 0; i < len(record); i++ {
		ele := record[i]
		list[i] = &pb.ArticleTopicListElement{
			TopicId:      uint32(ele.ID),
			Name:          ele.Name,
			Sort:           ele.Sort,
			State:          ele.State,
			CreateAt:       util.FormatDateTime(ele.CreatedAt),
		}
	}

	return &pb.GetArticleTopicListResponse{
		List:  list,
		Total: total,
	}, nil
}
