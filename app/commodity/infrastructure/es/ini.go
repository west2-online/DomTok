package es

import (
	elastic "github.com/elastic/go-elasticsearch"
	"github.com/west2-online/DomTok/app/commodity/domain/repository"
)

type CommodityElastic struct {
	client *elastic.Client
}

func NewCommodityElastic(client *elastic.Client) repository.CommodityElastic {
	return &CommodityElastic{client: client}
}
