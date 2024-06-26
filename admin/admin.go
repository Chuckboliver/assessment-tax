package admin

import (
	"context"
)

type AdminRepository interface {
	UpdatePersonalDeduction(ctx context.Context, personalDeduction float64) (float64, error)
	UpdateKReceiptDeduction(ctx context.Context, kReceiptDeduction float64) (float64, error)
}

type AdminService interface {
	UpdatePersonalDeduction(ctx context.Context, personalDeduction float64) (float64, error)
	UpdateKReceiptDeduction(ctx context.Context, kReceiptDeduction float64) (float64, error)
}

var _ AdminService = (*adminService)(nil)

type adminService struct {
	adminRepository AdminRepository
}

func NewAdminService(adminRepository AdminRepository) AdminService {
	return &adminService{
		adminRepository: adminRepository,
	}
}

func (a *adminService) UpdatePersonalDeduction(ctx context.Context, personalDeduction float64) (float64, error) {
	return a.adminRepository.UpdatePersonalDeduction(ctx, personalDeduction)
}

func (a *adminService) UpdateKReceiptDeduction(ctx context.Context, kReceiptDeduction float64) (float64, error) {
	return a.adminRepository.UpdateKReceiptDeduction(ctx, kReceiptDeduction)
}
