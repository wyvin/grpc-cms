package articleDb

import (
	"grpc-content/models"
	"time"
)

type Article struct {
	AppId          string
	GroupId        uint32
	AritcleId      uint32
	TopicId        uint32
	Cover          string
	Title          string
	Author         string
	Source         string
	Recommendation uint32
	Content        string
	State          uint32

	ArticleIdList []uint32
	StartReleased time.Time
	EndReleased   time.Time
	Page          uint32
	PerPage       uint32
}

func (a *Article) Migrate() error {
	return models.MigrateArticle()
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.AppId, a.GroupId, a.AritcleId)
}

func (a *Article) Insert() (uint32, error) {
	return models.InsertArticle(a.AppId, a.GroupId, a.TopicId, a.Cover, a.Title, a.Author, a.Source, a.Recommendation, a.Content, a.State)
}

func (a *Article) GetArticleList() ([]models.Article, uint32, error) {
	return models.GetArticleList(a.Page, a.PerPage, a.Title, a.Source, a.getMaps())
}

func (a *Article) GetArticlePlacementList() ([]models.Article, uint32, error) {
	return models.GetArticlePlacementList(a.Page, a.PerPage, a.Title, a.Source, a.getMaps())
}

func (a *Article) GetArticleDetail() (*models.Article, error) {
	return models.GetArticleDetail(a.AritcleId, a.getMaps())
}

func (a *Article) Edit() (uint32, error) {
	return models.EditArticle(a.AppId, a.GroupId, a.AritcleId, a.editMaps())
}

func (a *Article) Delete() (uint32, error) {
	return models.DeleteArticle(a.AppId, a.GroupId, a.ArticleIdList)
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["app_id"] = a.AppId
	maps["group_id"] = a.GroupId
	if a.TopicId > 0 {
		maps["topic_id"] = a.TopicId
	}
	if a.State > 0 {
		maps["state"] = a.State
	}
	return maps
}

func (a *Article) editMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["topic_id"] = a.TopicId
	maps["cover"] = a.Cover
	maps["title"] = a.Title
	maps["author"] = a.Author
	maps["source"] = a.Source
	maps["recommendation"] = a.Recommendation
	maps["content"] = a.Content
	if a.State > 0 {
		maps["state"] = a.State
	}
	if a.State == models.ArticleStateReleased {
		maps["released_at"] = time.Now()
	}
	return maps
}
