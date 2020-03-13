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
		var releasedAt string
		if ele.ReleasedAt.IsZero() {
			releasedAt = ""
		} else {
			releasedAt = util.FormatDateTime(ele.ReleasedAt)
		}
		list[i] = &pb.ArticleListElement{
			ArticleId:      uint32(ele.ID),
			TopicId:        ele.TopicId,
			Cover:          ele.Cover,
			Title:          ele.Title,
			Author:         ele.Author,
			Source:         ele.Source,
			Recommendation: ele.Recommendation,
			State:          ele.State,
			ReleasedAt:     releasedAt,
			CreatedAt:      util.FormatDateTime(ele.CreatedAt),
		}
	}

	return &pb.GetArticleListResponse{
		List:  list,
		Total: total,
	}, nil
}

func (a articleService) GetArticlePlacementList(ctx context.Context, req *pb.GetArticlePlacementListRequest) (*pb.GetArticlePlacementListResponse, error) {
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
	}
	record, total, err := articleInstance.GetArticlePlacementList()
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ArticlePlacementListElement, len(record))
	for i := 0; i < len(record); i++ {
		ele := record[i]
		list[i] = &pb.ArticlePlacementListElement{
			ArticleId:      uint32(ele.ID),
			TopicId:        ele.TopicId,
			Cover:          ele.Cover,
			Title:          ele.Title,
			Author:         ele.Author,
			Source:         ele.Source,
			Recommendation: ele.Recommendation,
			ReleasedAt:     util.FormatDateTime(ele.ReleasedAt),
		}
	}

	return &pb.GetArticlePlacementListResponse{
		List:  list,
		Total: total,
	}, nil
}

func (a articleService) EditArticle(ctx context.Context, req *pb.EditArticleRequest) (*pb.RowsAffectedResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}
	if req.ArticleId == 0 {
		return nil, status.Error(codes.InvalidArgument, "缺少article_id参数")
	}
	articleInstance := articleDb.Article{
		AppId:          req.AppId,
		GroupId:        req.GroupId,
		AritcleId:      req.ArticleId,
		TopicId:        req.TopicId,
		Cover:          req.Cover,
		Title:          req.Title,
		Author:         req.Author,
		Source:         req.Source,
		Recommendation: req.Recommendation,
		Content:        req.Content,
		State:          req.State,
	}
	exists, err := articleInstance.ExistByID()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "记录不存在")
	}

	var rowsAffected uint32
	rowsAffected, err = articleInstance.Edit()
	if err != nil {
		return nil, err
	}

	return &pb.RowsAffectedResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (a articleService) GetArticleDetail(ctx context.Context, req *pb.GetArticleDetailRequest) (*pb.GetArticleDetailResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}
	articleInstance := articleDb.Article{
		AppId:     req.AppId,
		GroupId:   req.GroupId,
		AritcleId: req.ArticleId,
	}

	article, err := articleInstance.GetArticleDetail()
	if err != nil {
		return nil, err
	}
	if article.ID == 0 {
		return nil, status.Error(codes.NotFound, "记录不存在")
	}

	var releasedAt string
	if article.ReleasedAt.IsZero() {
		releasedAt = ""
	} else {
		releasedAt = util.FormatDateTime(article.ReleasedAt)
	}

	return &pb.GetArticleDetailResponse{
		ArticleId:      uint32(article.ID),
		TopicId:        article.TopicId,
		Cover:          article.Cover,
		Title:          article.Title,
		Author:         article.Author,
		Source:         article.Source,
		Recommendation: article.Recommendation,
		Content:        article.Content,
		ReleasedAt:     releasedAt,
		State:          article.State,
		UpdatedAt:      util.FormatDateTime(article.UpdatedAt),
		CreatedAt:      util.FormatDateTime(article.CreatedAt),
	}, nil
}

func (a articleService) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.RowsAffectedResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}
	articleInstance := articleDb.Article{
		AppId:         req.AppId,
		GroupId:       req.GroupId,
		ArticleIdList: req.ArticleIdList,
	}
	rowsAffected, err := articleInstance.Delete()
	if err != nil {
		return nil, err
	}

	return &pb.RowsAffectedResponse{
		RowsAffected: rowsAffected,
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
		AppId:   req.AppId,
		GroupId: req.GroupId,
		Name:    req.Name,
		Sort:    req.Sort,
		State:   req.State,
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
			TopicId:   uint32(ele.ID),
			Name:      ele.Name,
			Sort:      ele.Sort,
			State:     ele.State,
			CreatedAt: util.FormatDateTime(ele.CreatedAt),
		}
	}

	return &pb.GetArticleTopicListResponse{
		List:  list,
		Total: total,
	}, nil
}

func (a articleService) EditArticleTopic(ctx context.Context, req *pb.EditArticleTopicRequest) (*pb.RowsAffectedResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}
	if req.TopicId == 0 {
		return nil, status.Error(codes.InvalidArgument, "缺少topic_id参数")
	}
	topicInstance := articleDb.ArticleTopic{
		AppId:         req.AppId,
		GroupId:       req.GroupId,
		TopicId: req.TopicId,
		Name: req.Name,
		Sort: req.Sort,
		State: req.State,
	}
	exists, err := topicInstance.ExistByID()
	if err != nil {
		return  nil, err
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "记录不存在")
	}

	var rowsAffected uint32
	rowsAffected, err = topicInstance.Edit()
	if err != nil {
		return nil, err
	}

	return &pb.RowsAffectedResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (a articleService) DeleteArticleTopic(ctx context.Context, req *pb.DeleteArticleTopicRequest) (*pb.RowsAffectedResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}
	topicInstance := articleDb.ArticleTopic{
		AppId:         req.AppId,
		GroupId:       req.GroupId,
		TopicIdList: req.TopicIdList,
	}
	rowsAffected, err := topicInstance.Delete()
	if err != nil {
		return nil, err
	}

	return &pb.RowsAffectedResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (a articleService) GetArticleTopicPlacementList(ctx context.Context, req *pb.GetArticleTopicPlacementListRequest) (*pb.GetArticleTopicPlacementListResponse, error) {
	if err := CheckAppIdAndGroupId(req.AppId, req.GroupId); err != nil {
		return nil, err
	}

	util.FormatPageAndPerPage(&req.Page, &req.PerPage)

	topicInstance := articleDb.ArticleTopic{
		AppId:   req.AppId,
		GroupId: req.GroupId,
		Page:    req.Page,
		PerPage: req.PerPage,
	}
	record, total, err := topicInstance.GetTopicPlacementList()
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ArticleTopicPlacementListElement, len(record))
	for i := 0; i < len(record); i++ {
		ele := record[i]
		list[i] = &pb.ArticleTopicPlacementListElement{
			TopicId:   uint32(ele.ID),
			Name:      ele.Name,
			Sort:      ele.Sort,
		}
	}

	return &pb.GetArticleTopicPlacementListResponse{
		List:  list,
		Total: total,
	}, nil
}

