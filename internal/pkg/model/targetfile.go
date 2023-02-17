package model

import "fmt"

type TargetFile struct {
	Path string
	Name string
	Tpe  string
}

func (target TargetFile) CanonicalName() string {
	return fmt.Sprintf("%s%s.%s", target.Path, target.Name, target.Tpe)
}
