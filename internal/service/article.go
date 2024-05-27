package service

import (
	"context"
	"time"
	"xiaoweishu/webook/internal/domain"
	"xiaoweishu/webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, id int64) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
	ListPub(ctx context.Context,
		start time.Time, offset, limit int) ([]domain.Article, error)
}
type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	//只是编辑文章，还没有到发表，所以状态设置成未发表
	//id>0,说明这是一篇老文章
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	} //新文章，直接创建，此时的创建只是存在于数据库中
	return a.repo.Create(ctx, art)

}

// Publish 也就是同步的意思，将制作库的东西同步到线上库中
func (a articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)

}

func (a articleService) Withdraw(ctx context.Context, uid int64, id int64) error {
	//隐藏文章，直接状态改成不可见或私人即可
	return a.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)

}

func (a articleService) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.GetByAuthor(ctx, uid, offset, limit)
}

func (a articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetById(ctx, id)
}

func (a articleService) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := a.repo.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return art, nil
}
func (a *articleService) ListPub(ctx context.Context,
	start time.Time, offset, limit int) ([]domain.Article, error) {
	return a.repo.ListPub(ctx, start, offset, limit)
}
