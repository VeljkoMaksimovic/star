package services

import (
	"context"
	oort "github.com/c12s/oort/pkg/api"
	"github.com/c12s/star/internal/domain"
)

type ConfigService struct {
	Repo      domain.ConfigRepo
	evaluator oort.OortEvaluatorClient
}

func NewConfigService(repo domain.ConfigRepo, evaluator oort.OortEvaluatorClient) (*ConfigService, error) {
	return &ConfigService{
		Repo:      repo,
		evaluator: evaluator,
	}, nil
}

func (c *ConfigService) Put(req domain.PutConfigGroupReq) (*domain.PutConfigGroupResp, error) {
	err := c.Repo.Put(req.Group)
	if err != nil {
		return nil, err
	}
	return &domain.PutConfigGroupResp{}, nil
}

func (c *ConfigService) Get(req domain.GetConfigGroupReq) (*domain.GetConfigGroupResp, error) {
	resp, err := c.evaluator.Authorize(context.TODO(), &oort.AuthorizationReq{
		Subject: &oort.Resource{
			Id:   req.SubId,
			Kind: req.SubKind,
		},
		Object: &oort.Resource{
			Id:   req.GroupId,
			Kind: "config",
		},
		PermissionName: "config.get",
	})
	if err != nil {
		return nil, err
	}
	if !resp.Authorized {
		return nil, domain.ErrUnauthorized()
	}

	cg, err := c.Repo.Get(req.GroupId)
	if err != nil {
		return nil, err
	}
	return &domain.GetConfigGroupResp{
		Group: *cg,
	}, nil
}
