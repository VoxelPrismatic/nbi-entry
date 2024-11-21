package notif

import "nbientry/web"

var _ = web.Migrate(BusinessUnit{})

type BusinessUnit struct {
	Id    int `gorm:"primaryKey"`
	Name  string
	Color string
}

func AllBusinessUnits() []BusinessUnit {
	return web.GetSorted(BusinessUnit{}, "name ASC")
}
