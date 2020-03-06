package adDb

import "grpc-content/models"

type Ad struct {
	Appid       string
	Groupid     uint32
	Adid        uint32
	AdidList    []uint32
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

func (a *Ad) ExistByName() (bool, error) {
	return models.ExistAdByName(a.Name)
}

func (a *Ad) ExistByID() (bool, error) {
	return models.ExistAdByID(a.Appid, a.Groupid, a.Adid)
}

func (a *Ad) Insert() (uint32, error) {
	return models.InsertAd(a.Appid, a.Groupid, a.Name, a.Title, a.Description, a.Remark, a.Cover, a.Url, a.Priority, a.Display, a.State)
}

func (a *Ad) GetAdList() ([]models.Ad, uint32, error) {
	return models.GetAdList(a.Page, a.PerPage, a.Name, a.Title, a.getMaps())
}

func (a *Ad) Edit() (uint32, error) {
	return models.EditAd(a.Appid, a.Groupid, a.Adid, a.editMaps())
}

func (a *Ad) Delete() (uint32, error) {
	return models.DeleteAd(a.Appid, a.Groupid, a.AdidList)
}

func (a *Ad) GetAdPlacementList() ([]models.Ad, uint32, error) {
	return models.GetAdPlacementList(a.Page, a.PerPage, a.getMaps())
}

func (a *Ad) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["appid"] = a.Appid
	maps["groupid"] = a.Groupid
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
	maps["priority"] = a.Priority
	maps["display"] = a.Display
	maps["state"] = a.State
	return maps
}


