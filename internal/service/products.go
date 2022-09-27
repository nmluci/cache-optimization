package service

import (
	"context"

	"github.com/nmluci/cache-optimization/internal/model"
)

var (
	logTagNewProduct        = "[ProductService-New]"
	logTagFindByIDProduct   = "[ProductService-GetByID]"
	logTagFindByIDProductNC = "[NC-ProductService-GetByID]"
	logTagFindAllProduct    = "[ProductService-All]"
	logTagFindAllProductNC  = "[NC-ProductService-All]"
	logTagEditProduct       = "[ProductService-Edit]"
	logTagDeleteProduct     = "[ProductService-Delete]"
)

func (s *service) FindProductByID(ctx context.Context, id uint64) (res *model.Product, err error) {
	res, err = s.repository.FindProductByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to fetch product data: %+v", logTagFindByIDProduct, err)
		return
	}

	if res == nil {
		s.logger.Errorf("%s productID: %d not found", logTagFindByIDProduct, id)
		return nil, nil
	}

	return res, nil
}

func (s *service) FindProducts(ctx context.Context) (res []*model.Product, err error) {
	res, err = s.repository.FindProducts(ctx)
	if err != nil {
		s.logger.Errorf("%s failed to fetch product data: %+v", logTagFindAllProduct, err)
		return
	}

	if res == nil {
		s.logger.Errorf("%s not found", logTagFindAllProduct)
		return nil, nil
	}

	return res, nil
}

func (s *service) ForceFindProductByID(ctx context.Context, id uint64) (res *model.Product, err error) {
	res, err = s.repository.ForceFindProductByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to fetch product data: %+v", logTagFindByIDProductNC, err)
		return
	}

	if res == nil {
		s.logger.Errorf("%s productID: %d not found", logTagFindByIDProductNC, id)
		return nil, nil
	}

	return res, nil
}

func (s *service) ForceFindProducts(ctx context.Context) (res []*model.Product, err error) {
	res, err = s.repository.ForceFindProducts(ctx)
	if err != nil {
		s.logger.Errorf("%s failed to fetch product data: %+v", logTagFindAllProductNC, err)
		return
	}

	if res == nil {
		s.logger.Errorf("%s not found", logTagFindAllProductNC)
		return nil, nil
	}

	return res, nil
}

func (s *service) InsertProduct(ctx context.Context, payload *model.Product) (err error) {
	err = s.repository.InsertNewProduct(ctx, payload)
	if err != nil {
		s.logger.Errorf("%s failed to insert product data: %+v", logTagNewProduct, err)
		return
	}

	return nil
}

func (s *service) UpdateProduct(ctx context.Context, id uint64, payload *model.Product) (err error) {
	err = s.repository.UpdateProduct(ctx, id, payload)
	if err != nil {
		s.logger.Errorf("%s failed to insert product data: %+v", logTagEditProduct, err)
		return
	}

	return
}

func (s *service) DeleteProduct(ctx context.Context, id uint64) (err error) {
	err = s.repository.DeleteProductByID(ctx, id)
	if err != nil {
		s.logger.Errorf("%s failed to insert product data: %+v", logTagDeleteProduct, err)
		return
	}

	return
}
