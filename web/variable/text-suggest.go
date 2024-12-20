package variable

import (
	"nbientry/web"
	"strings"
)

var _ = web.Migrate(&Autocomplete{})

type Autocomplete struct {
	VarId int
	Value string
}

type var_str_suggest struct {
	Real  VariableEntry
	T     Variable
	Value string
}

func (v var_str_suggest) GetNbiId() int {
	return v.Real.NbiId
}

func (v var_str_suggest) GetVarId() int {
	return v.Real.Id
}

func (v var_str_suggest) GetIndex() int {
	return v.Real.Index
}

func (v *var_str_suggest) SetIndex(idx int) {
	v.Real.Index = idx
}

func (v var_str_suggest) Type() Variable {
	if v.T.Id == 0 {
		v.T = web.GetFirst(Variable{Id: v.Real.Id})
	}
	return v.T
}

func (v *var_str_suggest) Set(value string) error {
	v.Value = value
	return nil
}

func (v var_str_suggest) Get() string {
	return v.Value
}

func (v *var_str_suggest) Reset() {
	v.Value = ""
}

func (v var_str_suggest) ToSqlEntry() VariableEntry {
	return v.Real
}

func (v var_str_suggest) ToTypedEntry() (VarEntry, error) {
	return &v, nil
}

func (v var_str_suggest) exists() bool {
	autocompletes := web.GetSorted(Autocomplete{VarId: v.Real.Id}, "Value ASC")
	lower := strings.ToLower(v.Value)
	for _, a := range autocompletes {
		if strings.ToLower(a.Value) == lower {
			return true
		}
	}
	return false
}

func (v var_str_suggest) list() []string {
	autocompletes := web.GetSorted(Autocomplete{VarId: v.Real.Id}, "Value ASC")
	var ret []string
	for _, a := range autocompletes {
		ret = append(ret, a.Value)
	}
	return ret
}
