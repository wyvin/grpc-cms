package models

import (
	"github.com/jinzhu/gorm"
)

type ArticleTopic struct {
	gorm.Model
	AppId   string `gorm:"type:varchar(64);NOT NULL;default:'';comment:'应用ID'"`
	GroupId uint32 `gorm:"NOT NULL;default:0;comment:'用户组ID'"`
	Name    string `gorm:"type:varchar(64);NOT NULL;default:'';comment:'专题名称'"`
	Sort    uint32 `gorm:"NOT NULL;default:1;comment:'排序'"`
	State   uint32 `gorm:"NOT NULL;default:1;comment:'状态 1已发布 2停用'"`
}

const (
	ArticleTopicStateReleased = 1 // 已发布
	ArticleTopicStateCease    = 2 // 停用
)

func ExistArticleTopicByID(appId string, groupId, topicId uint32) (bool, error) {
	var topic ArticleTopic
	err := DB.Select("id").Where("app_id = ? AND group_id = ?", appId, groupId).First(&topic, topicId).Error
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

// 编辑专题
func EditArticleTopic(appId string, groupId, topicId uint32, data interface{}) (uint32, error) {
	result := DB.Model(&ArticleTopic{}).Where("app_id = ? AND group_id = ? AND id = ?", appId, groupId, topicId).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return uint32(result.RowsAffected), nil
}

// 批量删除专题
func DeleteArticleTopic(appId string, groupId uint32, topicIdList []uint32) (uint32, error) {
	tx := DB.Begin()
	result := tx.Where("app_id = ? AND group_id = ? AND id IN (?)", appId, groupId, topicIdList).Delete(&ArticleTopic{})
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}
	// 将文章中的专题ID至为0
	err := tx.Model(&Article{}).Where("app_id = ? AND group_id = ? AND topic_id IN (?)", appId, groupId, topicIdList).Update("topic_id", 0).Error
	if err != nil {
		tx.Rollback()
		return 0, nil
	}
	tx.Commit()
	return uint32(result.RowsAffected), nil
}
