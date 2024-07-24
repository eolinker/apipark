package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/eolinker/apipark/service/upstream"

	"gorm.io/gorm"

	"github.com/eolinker/apipark/service/team"

	"github.com/google/uuid"

	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/service/project"

	"github.com/eolinker/go-common/store"

	"github.com/eolinker/apipark/service/api"

	api_dto "github.com/eolinker/apipark/module/api/dto"
)

var _ IApiModule = (*imlApiModule)(nil)
var (
	projectMustAsServer = map[string]bool{
		"as_server": true,
	}
)

type imlApiModule struct {
	teamService     team.ITeamService         `autowired:""`
	projectService  project.IProjectService   `autowired:""`
	apiService      api.IAPIService           `autowired:""`
	upstreamService upstream.IUpstreamService `autowired:""`
	transaction     store.ITransaction        `autowired:""`
}

func (i *imlApiModule) SimpleList(ctx context.Context, input *api_dto.ListInput) ([]*api_dto.ApiSimpleItem, error) {
	projectIds := input.Projects
	w := make(map[string]interface{})
	if len(projectIds) > 0 {
		w["project"] = projectIds
	}
	list, err := i.apiService.Search(ctx, "", w)
	apiInfos, err := i.apiService.ListInfo(ctx, utils.SliceToSlice(list, func(s *api.API) string {
		return s.UUID
	})...)
	if err != nil {
		return nil, err
	}

	out := utils.SliceToSlice(apiInfos, func(item *api.APIInfo) *api_dto.ApiSimpleItem {
		return &api_dto.ApiSimpleItem{
			Id:     item.UUID,
			Name:   item.Name,
			Method: item.Method,
			Path:   item.Path,
		}
	})
	return out, nil
}

func (i *imlApiModule) Detail(ctx context.Context, pid string, aid string) (*api_dto.ApiDetail, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}

	detail, err := i.apiService.GetInfo(ctx, aid)
	if err != nil {
		return nil, err
	}

	apiDetail := &api_dto.ApiDetail{
		ApiSimpleDetail: *api_dto.GenApiSimpleDetail(detail),
	}
	proxy, err := i.apiService.LatestProxy(ctx, aid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if proxy != nil {

		apiDetail.Proxy = api_dto.FromServiceProxy(proxy.Data)
	}

	document, err := i.apiService.LatestDocument(ctx, aid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if document != nil {
		doc := make(map[string]interface{})
		err = json.Unmarshal([]byte(document.Data.Content), &doc)
		if err != nil {
			return nil, err
		}
		apiDetail.Doc = doc
	}

	return apiDetail, nil
}

func (i *imlApiModule) SimpleDetail(ctx context.Context, pid string, aid string) (*api_dto.ApiSimpleDetail, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}

	detail, err := i.apiService.GetInfo(ctx, aid)
	if err != nil {
		return nil, err
	}

	return api_dto.GenApiSimpleDetail(detail), nil
}

func (i *imlApiModule) Search(ctx context.Context, keyword string, pid string) ([]*api_dto.ApiItem, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}

	list, err := i.apiService.Search(ctx, keyword, map[string]interface{}{
		"project": pid,
	})
	if err != nil {
		return nil, err
	}
	apiInfos, err := i.apiService.ListInfo(ctx, utils.SliceToSlice(list, func(s *api.API) string {
		return s.UUID
	})...)
	if err != nil {
		return nil, err
	}
	utils.Sort(apiInfos, func(a, b *api.APIInfo) bool {
		return a.UpdateAt.After(b.UpdateAt)
	})
	out := utils.SliceToSlice(apiInfos, func(item *api.APIInfo) *api_dto.ApiItem {
		return &api_dto.ApiItem{
			Id:         item.UUID,
			Name:       item.Name,
			Method:     item.Method,
			Path:       item.Path,
			Creator:    auto.UUID(item.Creator),
			Updater:    auto.UUID(item.Updater),
			CreateTime: auto.TimeLabel(item.CreateAt),
			UpdateTime: auto.TimeLabel(item.UpdateAt),
			CanDelete:  true,
		}
	})

	return out, nil
}

func (i *imlApiModule) SimpleSearch(ctx context.Context, keyword string, pid string) ([]*api_dto.ApiSimpleItem, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}

	list, err := i.apiService.Search(ctx, keyword, map[string]interface{}{
		"project": pid,
	})
	apiInfos, err := i.apiService.ListInfo(ctx, utils.SliceToSlice(list, func(s *api.API) string {
		return s.UUID
	})...)
	if err != nil {
		return nil, err
	}
	out := utils.SliceToSlice(apiInfos, func(item *api.APIInfo) *api_dto.ApiSimpleItem {
		return &api_dto.ApiSimpleItem{
			Id:     item.UUID,
			Name:   item.Name,
			Method: item.Method,
			Path:   item.Path,
		}
	})
	return out, nil
}

func (i *imlApiModule) Create(ctx context.Context, pid string, dto *api_dto.CreateApi) (*api_dto.ApiSimpleDetail, error) {
	info, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}
	prefix, err := i.Prefix(ctx, pid)
	if err != nil {
		return nil, err
	}
	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		if dto.Id == "" {
			dto.Id = uuid.New().String()
		}
		err = dto.Validate()
		if err != nil {
			return err
		}

		path := fmt.Sprintf("%s%s", prefix, dto.Path)
		err = i.apiService.Exist(ctx, "", &api.ExistAPI{Path: dto.Path, Method: dto.Method})
		if err != nil {
			return fmt.Errorf("api path %s,method: %s already exist", dto.Path, dto.Method)
		}
		proxy := api_dto.ToServiceProxy(dto.Proxy)
		err = i.apiService.SaveProxy(ctx, dto.Id, proxy)
		if err != nil {
			return err
		}
		err = i.apiService.SaveDocument(ctx, dto.Id, api_dto.ToServiceDocument(nil))
		if err != nil {
			return err
		}

		match, _ := json.Marshal(dto.MatchRules)
		return i.apiService.Create(ctx, &api.CreateAPI{
			UUID:        dto.Id,
			Name:        dto.Name,
			Description: dto.Description,
			Project:     pid,
			Team:        info.Team,
			Method:      dto.Method,
			Path:        path,
			Match:       string(match),
			//Upstream:    proxy.Upstream,
		})
	})
	if err != nil {
		return nil, err
	}
	return i.SimpleDetail(ctx, pid, dto.Id)
}

func (i *imlApiModule) Edit(ctx context.Context, pid string, aid string, dto *api_dto.EditApi) (*api_dto.ApiSimpleDetail, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}

	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		var up *string
		if dto.Proxy != nil {
			err = i.apiService.SaveProxy(ctx, aid, api_dto.ToServiceProxy(dto.Proxy))
			if err != nil {
				return err
			}
			//if dto.Proxy.Upstream != "" {
			//	up = &dto.Proxy.Upstream
			//}
		}
		err = i.apiService.Save(ctx, aid, &api.EditAPI{
			Name:        dto.Info.Name,
			Description: dto.Info.Description,
			Upstream:    up,
		})
		if err != nil {
			return err
		}

		if dto.Doc != nil {
			err = i.apiService.SaveDocument(ctx, aid, api_dto.ToServiceDocument(*dto.Doc))
			if err != nil {
				return err
			}
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return i.SimpleDetail(ctx, pid, aid)
}

func (i *imlApiModule) Delete(ctx context.Context, pid string, aid string) error {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return err
	}
	return i.apiService.Delete(ctx, aid)
}

func (i *imlApiModule) Copy(ctx context.Context, pid string, aid string, dto *api_dto.CreateApi) (*api_dto.ApiSimpleDetail, error) {
	info, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}
	oldApi, err := i.apiService.Get(ctx, aid)
	if err != nil {
		return nil, err
	}
	prefix, err := i.Prefix(ctx, pid)
	if err != nil {
		return nil, err
	}
	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		if dto.Id == "" {
			dto.Id = uuid.New().String()
		}
		err = dto.Validate()
		if err != nil {
			return err
		}

		path := fmt.Sprintf("%s/%s", strings.TrimSuffix(prefix, "/"), strings.TrimPrefix(dto.Path, "/"))
		err = i.apiService.Exist(ctx, pid, &api.ExistAPI{Path: path, Method: dto.Method})
		if err != nil {
			return err
		}

		proxy, err := i.apiService.LatestProxy(ctx, oldApi.UUID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		//upstreamId := ""
		if proxy != nil {
			err = i.apiService.SaveProxy(ctx, dto.Id, proxy.Data)
			if err != nil {
				return err
			}
			//upstreamId = proxy.Data.Upstream
		}

		doc, err := i.apiService.LatestDocument(ctx, oldApi.UUID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		if doc != nil {
			err = i.apiService.SaveDocument(ctx, dto.Id, doc.Data)
			if err != nil {
				return err
			}
		}
		match, _ := json.Marshal(dto.MatchRules)
		return i.apiService.Create(ctx, &api.CreateAPI{
			UUID:    dto.Id,
			Name:    dto.Name,
			Project: pid,
			Team:    info.Team,
			Method:  dto.Method,
			Path:    path,
			Match:   string(match),
			//Upstream: upstreamId,
		})

	})
	if err != nil {
		return nil, err
	}
	return i.SimpleDetail(ctx, pid, dto.Id)
}

func (i *imlApiModule) ApiDocDetail(ctx context.Context, pid string, aid string) (*api_dto.ApiDocDetail, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}

	apiBase, err := i.apiService.Get(ctx, aid)
	if err != nil {
		return nil, err
	}
	if apiBase.IsDelete {
		return nil, errors.New("api is delete")
	}

	detail, err := i.apiService.GetInfo(ctx, apiBase.UUID)
	if err != nil {
		return nil, err
	}
	document, err := i.apiService.LatestDocument(ctx, aid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	var doc map[string]interface{}
	if document != nil {
		doc = make(map[string]interface{})
		err = json.Unmarshal([]byte(document.Data.Content), &doc)
		if err != nil {
			return nil, err
		}
	}
	return &api_dto.ApiDocDetail{
		ApiSimpleDetail: *api_dto.GenApiSimpleDetail(detail),
		Doc:             doc,
	}, nil
}

func (i *imlApiModule) ApiProxyDetail(ctx context.Context, pid string, aid string) (*api_dto.ApiProxyDetail, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return nil, err
	}
	apiBase, err := i.apiService.Get(ctx, aid)
	if err != nil {
		return nil, err
	}
	if apiBase.IsDelete {
		return nil, errors.New("api is delete")
	}
	if apiBase.Project != pid {
		return nil, errors.New("api is not in project")
	}

	detail, err := i.apiService.GetInfo(ctx, aid)
	if err != nil {
		return nil, err
	}

	apiDetail := &api_dto.ApiProxyDetail{
		ApiSimpleDetail: *api_dto.GenApiSimpleDetail(detail),
	}
	proxy, err := i.apiService.LatestProxy(ctx, aid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	if proxy != nil {
		apiDetail.Proxy = api_dto.FromServiceProxy(proxy.Data)
	}
	return apiDetail, nil

}

func (i *imlApiModule) Prefix(ctx context.Context, pid string) (string, error) {
	pInfo, err := i.projectService.CheckProject(ctx, pid, projectMustAsServer)
	if err != nil {
		return "", err
	}

	if pInfo.Prefix != "" {
		if pInfo.Prefix[0] != '/' {
			pInfo.Prefix = fmt.Sprintf("/%s", strings.TrimSuffix(pInfo.Prefix, "/"))
		}
	}
	return strings.TrimSuffix(pInfo.Prefix, "/"), nil
}
