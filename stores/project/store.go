package project

import (
	"github.com/eolinker/ap-account/store/member"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/store"
)
import "reflect"

type IProjectStore interface {
	store.ISearchStore[Project]
}
type imlProjectStore struct {
	store.SearchStore[Project]
}

type IMemberStore member.IMemberStore

//type IMemberRoleStore interface {
//	store.IBaseStore[MemberRole]
//}
//type imlMemberRoleStore struct {
//	store.Store[MemberRole]
//}

type IAuthorizationStore interface {
	store.ISearchStore[Authorization]
}

type imlAuthorizationStore struct {
	store.SearchStore[Authorization]
}

func init() {
	autowire.Auto[IProjectStore](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectStore))
	})
	autowire.Auto[IMemberStore](func() reflect.Value {
		return reflect.ValueOf(member.NewMemberStore("project"))
	})
	//autowire.Auto[IMemberRoleStore](func() reflect.Value {
	//	return reflect.ValueOf(new(imlMemberRoleStore))
	//})
	autowire.Auto[IAuthorizationStore](func() reflect.Value {
		return reflect.ValueOf(new(imlAuthorizationStore))
	})
}
