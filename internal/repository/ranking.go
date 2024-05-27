package repository

import (
	"context"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}
type CacheRankingRepository struct {
	cache cache.RankingCache //这种是比较抽象的写法，方便测试

	//V1写法，可读性更好，但是不算是面向接口编程，测试性一般
	redisCache *cache.RankingRedisCache
	localCache *cache.RankingLocalCache
}
