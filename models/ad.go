package models

import (
	"github.com/jinzhu/gorm"
)

type Ad struct {
	gorm.Model
	Appid       string `gorm:"type:varchar(64);NOT NULL;default:'';comment:'应用ID'"`
	Groupid     uint32 `gorm:"NOT NULL;default:0;comment:'用户组ID'"`
	Name        string `gorm:"type:varchar(64);NOT NULL;default:'';comment:'广告名'"`
	Title       string `gorm:"type:varchar(64);NOT NULL;default:'';comment:'标题'"`
	Description string `gorm:"type:varchar(128);NOT NULL;default:'';comment:'描述'"`
	Remark      string `gorm:"type:varchar(256);NOT NULL;default:'';comment:'备注'"`
	Cover       string `gorm:"type:varchar(256);NOT NULL;default:'';comment:'封面'"`
	Url         string `gorm:"type:varchar(512);NOT NULL;default:'';comment:'链接'"`
	Priority    uint32 `gorm:"NOT NULL;default:1;comment:'优先级'"`
	Display     uint32 `gorm:"NOT NULL;default:1;comment:'展示方式 1首页banner'"`
	State       uint32 `gorm:"NOT NULL;default:1;comment:'状态 1未发布 2已发布 3停止'"`
}

const (
	AdStateUnreleased = 1
	AdStateReleased   = 2
	AdStateCeased     = 3
)

func MigrateAd() error {
	if err := db.AutoMigrate(&Ad{}).Error; err != nil {
		return err
	}
	return nil
}

func ExistAdByName(name string) (bool, error) {
	var ad Ad
	err := db.Select("id").Where("name = ?", name).First(&ad).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if ad.ID > 0 {
		return true, nil
	}

	return false, nil
}

func ExistAdByID(appid string, groupid, adid uint32) (bool, error) {
	var ad Ad
	err := db.Select("id").Where("appid = ? AND groupid = ? AND id = ?", appid, groupid, adid).First(&ad).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if ad.ID > 0 {
		return true, nil
	}
	return false, nil
}

// 插入广告
func InsertAd(appid string, groupid uint32, name, title, description, remark, cover, url string, priority, display, state uint32) (uint32, error) {
	ad := Ad{
		Appid:       appid,
		Groupid:     groupid,
		Name:        name,
		Title:       title,
		Description: description,
		Remark:      remark,
		Cover:       cover,
		Url:         url,
		Priority:    priority,
		Display:     display,
		State:       state,
	}
	if err := db.Create(&ad).Error; err != nil {
		return 0, err
	}

	return uint32(ad.ID), nil
}

// 获取广告列表
func GetAdList(page, perPage uint32, name, title string, maps interface{}) ([]Ad, uint32, error) {
	var (
		adList []Ad
		total  uint32
		err    error
	)

	record := db.Model(&Ad{}).Where(maps)
	if name != "" {
		record = record.Where("name LIKE ?", "%" + name + "%")
	}
	if title != "" {
		record = record.Where("title LIKE ?", "%" + title + "%")
	}
	err = record.Offset(page).Limit(perPage).Find(&adList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	record.Count(&total)

	return adList, total, nil
}

// 编辑广告
func EditAd(appid string, groupid, adid uint32, data interface{}) (uint32, error) {
	result := db.Model(&Ad{}).Where("appid = ? AND groupid = ? AND id = ?", appid, groupid, adid).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return  uint32(result.RowsAffected), nil
}

// 批量删除广告
func DeleteAd(appid string, groupid uint32, adidList []uint32) (uint32, error) {
	result := db.Where("appid = ? AND groupid = ? AND id IN (?)", appid, groupid, adidList).Delete(&Ad{})
	if result.Error != nil {
		return 0, result.Error
	}
	return  uint32(result.RowsAffected), nil
}

// 获取广告投放列表
func GetAdPlacementList(page, perPage uint32, maps interface{}) ([]Ad, uint32, error) {
	var (
		adPlacementList []Ad
		total           uint32
		err             error
	)
	record := db.Model(&Ad{}).Where(maps).Where("state = ?", AdStateReleased)
	err = record.Offset(page).Limit(perPage).Find(&adPlacementList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	if err = record.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return adPlacementList, total, nil
}
