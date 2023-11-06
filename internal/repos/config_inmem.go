package repos

import (
	"github.com/c12s/star/internal/domain"
	"log"
)

type configInMemRepo struct {
	Groups map[string]*domain.ConfigGroup
}

func NewConfigInMemRepo() (domain.ConfigRepo, error) {
	return &configInMemRepo{
		Groups: make(map[string]*domain.ConfigGroup),
	}, nil
}

func (c configInMemRepo) Put(group domain.ConfigGroup) error {
	c.Groups[group.Id] = &group
	return nil
}

func (c configInMemRepo) Get(groupId string) (*domain.ConfigGroup, error) {
	log.Println("Printing all groups")
	log.Printf("%+v\n", c.Groups)
	group, ok := c.Groups[groupId]
	if !ok {
		return nil, domain.ErrNotFound()
	}
	return group, nil
}
