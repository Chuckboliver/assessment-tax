package admin

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpdatePersonalDeduction(t *testing.T) {
	ctrl := gomock.NewController(t)
	adminRepo := NewMockAdminRepository(ctrl)
	adminService := NewAdminService(adminRepo)

	adminRepo.EXPECT().UpdatePersonalDeduction(gomock.Any(), 20000.0).Times(1).Return(20000.0, nil)

	updatePersonalDeduction, err := adminService.UpdatePersonalDeduction(context.Background(), 20000.0)
	require.NoError(t, err)
	require.Equal(t, 20000.0, updatePersonalDeduction)
}
