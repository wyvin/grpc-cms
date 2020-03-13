package adDb

import (
	"grpc-cms/models"
)

type Ad struct {
	AppId       string
	GroupId     uint32
	AdId        uint32
	AdIdList    []uint32
	Name        string
	Title       string
	Description string
	Remark      string
	Cover       string
	Url         string
	Priority    uint32
	Display     uint32
	State       uint32

	Page    uint32
	PerPage uint32
}

func (a *Ad) Migrate() error {
	return models.MigrateAd()
}

func (a *Ad) ExistByID() (bool, error) {
	return models.ExistAdByID(a.AppId, a.GroupId, a.AdId)
}

func (a *Ad) Insert() (uint32, error) {
	return models.InsertAd(a.AppId, a.GroupId, a.Name, a.Title, a.Description, a.Remark, a.Cover, a.Url, a.Priority, a.Display, a.State)
}

func (a *Ad) GetAdList() ([]models.Ad, uint32, error) {
	return models.GetAdList(a.Page, a.PerPage, a.Name, a.Title, a.getMaps())
}

func (a *Ad) Edit() (uint32, error) {
	return models.EditAd(a.AppId, a.GroupId, a.AdId, a.editMaps())
}

func (a *Ad) Delete() (uint32, error) {
	return models.DeleteAd(a.AppId, a.GroupId, a.AdIdList)
}

func (a *Ad) GetAdPlacementList() ([]models.Ad, uint32, error) {
	return models.GetAdPlacementList(a.Page, a.PerPage, a.getMaps())
}

func (a *Ad) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["app_id"] = a.AppId
	maps["group_id"] = a.GroupId
	if a.Display > 0 {
		maps["display"] = a.Display
	}
	if a.State > 0 {
		maps["state"] = a.State
	}
	return maps
}

func (a *Ad) editMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["name"] = a.Name
	maps["title"] = a.Title
	maps["description"] = a.Description
	maps["remark"] = a.Remark
	maps["cover"] = a.Cover
	maps["url"] = a.Url
	if a.Priority > 0 {
		maps["priority"] = a.Priority
	}
	if a.Display > 0 {
		maps["display"] = a.Display
	}
	if a.State > 0 {
		maps["state"] = a.State
	}
	return maps
}
