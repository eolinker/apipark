package project_member

import (
	"reflect"

	"github.com/eolinker/ap-account/service/member"
	"github.com/eolinker/apipark/stores/project"
	"github.com/eolinker/go-common/autowire"
)

type IMemberService member.IMemberService

func init() {
	autowire.Auto[IMemberService](func() reflect.Value {
		return reflect.ValueOf(new(member.Service[project.IMemberStore]))
	})
}
