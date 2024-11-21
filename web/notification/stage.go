package notif

import "nbientry/web"

var _ = web.Migrate(Stage{})

type Stage struct {
	Id    int `gorm:"primaryKey"`
	Name  string
	Index int
}

func AllStages() []Stage {
	return web.GetSorted(Stage{}, "`index` ASC")
}

func (s Stage) Top() bool {
	return s.Index == 1
}

func (s Stage) Bottom() bool {
	return s.Index == len(AllStages())
}

func (s *Stage) Increment() {
	if s.Bottom() {
		return
	}

	below := web.GetFirst(Stage{Index: s.Index + 1})
	s.Index++
	below.Index--

	web.Save(s)
	web.Save(below)
}

func (s *Stage) Decrement() {
	if s.Top() {
		return
	}

	above := web.GetFirst(Stage{Index: s.Index - 1})
	s.Index--
	above.Index++

	web.Save(s)
	web.Save(above)
}

func (s Stage) New() Stage {
	stages := []Stage{}
	web.Db().Model(&Stage{}).Where("`index` > ?", s.Index).Find(&stages)

	for i := range stages {
		stages[i].Index++
		web.Save(&stages[i])
	}

	ret := Stage{Name: "", Index: s.Index + 1}
	web.Save(&ret)
	return ret
}