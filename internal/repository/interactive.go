package repository

import (
	"context"
	"errors"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/pkg/logger"
	"xiaoweishu/webook/internal/repository/cache"
	"xiaoweishu/webook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}
type CachedInteractiveRepository struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.LoggerV1
}

func (c CachedInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	//这时候需要设置缓存，否则会导致数据不一致问题
	err = c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return nil
}

func (c CachedInteractiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	//一个人只能点赞一次，所以点赞的时候直接插入点赞的表格即可
	//再次点赞就更新UTime即可
	err := c.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return nil
}

func (c CachedInteractiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {
	return c.dao.DeleteLikeInfo(ctx, biz, id, uid)
}

func (c CachedInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	var res = dao.UserCollectionBiz{
		Uid:   uid,
		BizId: id,
		Biz:   biz,
		Cid:   cid,
	}
	err := c.dao.InsertCollectionBiz(ctx, res)
	if err != nil {
		return err
	}
	return nil
}

func (c *CachedInteractiveRepository) Get(ctx context.Context, biz string, id int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, id)
	if err == nil {
		return intr, err
	}
	ie, err := c.dao.Get(ctx, biz, id)
	if err != nil {
		return domain.Interactive{}, err
	}
	res := c.ToDomain(ie)
	//回写缓存
	err = c.cache.Set(ctx, biz, id, res)
	if err != nil {
		c.l.Error("回写缓存失败", logger.String("biz", biz),
			logger.Int64("bizId", id),
			logger.Error(err))
	}
	return res, nil

}

// 关于用户喜欢的逻辑，这里定义成，若用户喜欢，那么就会生成喜欢的表，取消喜欢赞时，就会把对应的表格删除
// 所以只要dao层能找到该表，那就表明了用户点赞，否则就是没有点赞
func (c CachedInteractiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, dao.ErrRecordNotFound):
		return false, err
	default:
		return false, err
	}
}

func (c CachedInteractiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, err
	default:
		return false, err
	}
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO,
	cache cache.InteractiveCache,
	l logger.LoggerV1) InteractiveRepository {
	return &CachedInteractiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}

}
func (c *CachedInteractiveRepository) ToDomain(ie dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    ie.ReadCnt,
		CollectCnt: ie.CollectCnt,
		LikeCnt:    ie.LikeCnt,
	}
}
