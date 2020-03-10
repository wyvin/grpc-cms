package models

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Article struct {
	gorm.Model
	AppId          string    `gorm:"type:varchar(64);NOT NULL;default:'';comment:'应用ID'"`
	GroupId        uint32    `gorm:"NOT NULL;default:0;comment:'用户组ID'"`
	TopicId        uint32    `gorm:"NOT NULL;default:0;comment:'专题ID'"`
	Cover          string    `gorm:"type:varchar(256);NOT NULL;default:'';comment:'封面'"`
	Title          string    `gorm:"type:varchar(100);NOT NULL;default:'';comment:'标题'"`
	Author         string    `gorm:"type:varchar(64);NOT NULL;default:'';comment:'作者'"`
	Source         string    `gorm:"type:varchar(128);NOT NULL;default:'';comment:'来源'"`
	Recommendation uint32    `gorm:"NOT NULL;default:0;comment:'推荐度'"`
	Content        string    `gorm:"type:text;NOT NULL;comment:'内容/正文'"`
	ReleasedAt     time.Time `gorm:"default:NULL;comment:'发布时间'"`
	State          uint32    `gorm:"NOT NULL;default:1;comment:'状态 1草稿 2已发布 3下架'"`
}

const (
	ArticleStateDraft    = 1 // 草稿
	ArticleStateReleased = 2 // 已发布
	ArticleStateCeased   = 3 // 下架
)

func MigrateArticle() error {
	if err := DB.AutoMigrate(&Article{}).Error; err != nil {
		return err
	}
	return nil
}

func ExistArticleByID(appId string, groupId, articleId uint32) (bool, error) {
	var article Article
	err := DB.Select("id").Where("app_id = ? AND group_id = ?", appId, groupId).First(&article, articleId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 插入文章
func InsertArticle(appId string, groupId, topicId uint32, cover, title, author string, source string, recommendation uint32, content string, state uint32) (uint32, error) {
	var releasedAt time.Time
	if state == ArticleStateReleased {
		releasedAt = time.Now()
	}
	article := Article{
		AppId:          appId,
		GroupId:        groupId,
		TopicId:        topicId,
		Cover:          cover,
		Title:          title,
		Author:         author,
		Source:         source,
		Recommendation: recommendation,
		Content:        content,
		ReleasedAt:     releasedAt,
		State:          state,
	}
	if err := DB.Create(&article).Error; err != nil {
		return 0, err
	}

	return uint32(article.ID), nil
}

// 获取文章列表
func GetArticleList(page, perPage uint32, title, source string, maps interface{}) ([]Article, uint32, error) {
	var (
		articleList []Article
		total       uint32
		err         error
	)

	record := DB.Model(&Article{}).Where(maps)
	if title != "" {
		record = record.Where("title LIKE ?", "%"+title+"%")
	}
	if source != "" {
		record = record.Where("source LIKE ?", "%"+source+"%")
	}
	err = record.Offset(page).Limit(perPage).Find(&articleList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	record.Count(&total)

	return articleList, total, nil
}

// 获取文章详情
func GetArticleDetail(articleId uint32, maps interface{}) (*Article, error) {
	var (
		article Article
		err     error
	)

	record := DB.Where(maps)
	record.Joins("JOIN article_topic ON article.topic_id = article_topic.id").First(&article, articleId)
	err = record.Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
    log.Println(&article)
	return &article, nil
}

// 编辑文章
func EditArticle(appId string, groupId, articleId uint32, data interface{}) (uint32, error) {
	result := DB.Model(&Article{}).Where("app_id = ? AND group_id = ? AND id = ?", appId, groupId, articleId).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return uint32(result.RowsAffected), nil
}