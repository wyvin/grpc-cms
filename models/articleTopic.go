package models

import (
	"github.com/jinzhu/gorm"
)

type ArticleTopic struct {
	gorm.Model
	AppId    string `gorm:"type:varchar(64);NOT NULL;default:'';comment:'应用ID'"`
	GroupId  uint32 `gorm:"NOT NULL;default:0;comment:'用户组ID'"`
	Name     string `gorm:"type:varchar(64);NOT NULL;default:'';comment:'专题名称'"`
	Sort     uint32 `gorm:"NOT NULL;default:0;comment:'排序'"`
	State    uint32 `gorm:"NOT NULL;default:1;comment:'状态 1已发布 2停用'"`
	Articles []Article
}

const (
	ArticleTopicStateReleased = 1 // 已发布
	ArticleTopicStateCease    = 2 // 停用
)

func ExistArticleTopicByName(appId string, groupId uint32, name string) (bool, error) {
	var topic ArticleTopic
	err := DB.Select("id").Where("app_id = ? AND group_id = ? AND name = ?", appId, groupId, name).First(&topic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if topic.ID > 0 {
		return true, nil
	}
	return false, nil
}

func MigrateArticleTopic() error {
	if err := DB.AutoMigrate(&ArticleTopic{}).Error; err != nil {
		return err
	}
	return nil
}

// 插入文章专题
func InsertArticleTopic(appId string, groupId uint32, name string, sort, state uint32) (uint32, error) {
	articleTopic := ArticleTopic{
		AppId:   appId,
		GroupId: groupId,
		Name:    name,
		Sort:    sort,
		State:   state,
	}
	if err := DB.Create(&articleTopic).Error; err != nil {
		return 0, err
	}

	return uint32(articleTopic.ID), nil
}

// 获取专题列表
func GetArticleTopicList(page, perPage uint32, maps interface{}) ([]ArticleTopic, uint32, error) {
	var (
		topicList []ArticleTopic
		total     uint32
		err       error
	)

	record := DB.Model(&ArticleTopic{}).Where(maps)
	err = record.Offset(page).Limit(perPage).Find(&topicList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	record.Count(&total)

	return topicList, total, nil
}
