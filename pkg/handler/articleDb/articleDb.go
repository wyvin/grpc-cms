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
