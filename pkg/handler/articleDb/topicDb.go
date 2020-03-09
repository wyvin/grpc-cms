package articleDb

import (
	"grpc-content/models"
)

type ArticleTopic struct {
	AppId   string
	GroupId uint32
	TopicId uint32
	Name    string
	Sort    uint32
	State   uint32

	Page    uint32
	PerPage uint32
}

func (a *ArticleTopic) ExistByName() (bool, error) {
	return models.ExistArticleTopicByName(a.AppId, a.GroupId, a.Name)
}

func (a *ArticleTopic) Migrate() error {
	return models.MigrateArticleTopic()
}

func (a *ArticleTopic) Insert() (uint32, error) {
	return models.InsertArticleTopic(a.AppId, a.GroupId, a.Name, a.Sort, a.State)
}

func (a *ArticleTopic) GetTopicList() ([]models.ArticleTopic, uint32, error) {
	return models.GetArticleTopicList(a.Page, a.PerPage, a.getMaps())
}

func (a *ArticleTopic) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["app_id"] = a.AppId
	maps["group_id"] = a.GroupId
	if a.State > 0 {
		maps["state"] = a.State
	}
	return maps
}
