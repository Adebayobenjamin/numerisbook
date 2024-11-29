package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Adebayobenjamin/numerisbook/pkg/helper"
	"github.com/Adebayobenjamin/numerisbook/pkg/models"
	repository_mocks "github.com/Adebayobenjamin/numerisbook/pkg/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupReminderTest(t *testing.T) (*repository_mocks.MockReminderRepository, *reminderService) {
	ctrl := gomock.NewController(t)
	mockRepo := repository_mocks.NewMockReminderRepository(ctrl)
	service := NewReminderService(mockRepo).(*reminderService)
	return mockRepo, service
}

func TestGetReminderDateFromSchedule(t *testing.T) {
	dueDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		schedule models.InvoiceReminderSchedule
		want     time.Time
	}{
		{
			name:     "14 days before due",
			schedule: models.InvoiceReminderSchedule14DaysBeforeDue,
			want:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "7 days before due",
			schedule: models.InvoiceReminderSchedule7DaysBeforeDue,
			want:     time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "3 days before due",
			schedule: models.InvoiceReminderSchedule3DaysBeforeDue,
			want:     time.Date(2024, 3, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "1 day before due",
			schedule: models.InvoiceReminderSchedule1DayBeforeDue,
			want:     time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "unknown schedule",
			schedule: "UNKNOWN_SCHEDULE",
			want:     dueDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getReminderDateFromSchedule(tt.schedule, dueDate)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetInvoiceReminders(t *testing.T) {
	mockRepo, service := setupReminderTest(t)
	ctx := context.Background()
	dueDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		invoice    *models.Invoice
		customerID uint
		schedules  map[models.InvoiceReminderSchedule]bool
		mockSetup  func([]models.InvoiceReminder)
		wantErr    bool
		errMsg     string
	}{
		{
			name: "successful reminders creation",
			invoice: &models.Invoice{
				ID:      1,
				DueDate: dueDate,
			},
			customerID: 1,
			schedules: map[models.InvoiceReminderSchedule]bool{
				models.InvoiceReminderSchedule14DaysBeforeDue: true,
				models.InvoiceReminderSchedule7DaysBeforeDue:  true,
			},
			mockSetup: func(expectedReminders []models.InvoiceReminder) {
				mockRepo.EXPECT().
					UpsertReminders(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, reminders []models.InvoiceReminder) error {
						assert.Len(t, reminders, 2)
						for _, reminder := range reminders {
							assert.Equal(t, uint(1), reminder.InvoiceID)
							assert.Equal(t, uint(1), reminder.CustomerID)
							assert.Nil(t, reminder.DeletedAt)
						}
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "reminders with some disabled",
			invoice: &models.Invoice{
				ID:      1,
				DueDate: dueDate,
			},
			customerID: 1,
			schedules: map[models.InvoiceReminderSchedule]bool{
				models.InvoiceReminderSchedule14DaysBeforeDue: true,
				models.InvoiceReminderSchedule7DaysBeforeDue:  false,
			},
			mockSetup: func(expectedReminders []models.InvoiceReminder) {
				mockRepo.EXPECT().
					UpsertReminders(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, reminders []models.InvoiceReminder) error {
						assert.Len(t, reminders, 2)
						for _, reminder := range reminders {
							if reminder.Schedule == models.InvoiceReminderSchedule7DaysBeforeDue {
								assert.NotNil(t, reminder.DeletedAt)
							} else {
								assert.Nil(t, reminder.DeletedAt)
							}
						}
						return nil
					})
			},
			wantErr: false,
		},
		{
			name: "repository error",
			invoice: &models.Invoice{
				ID:      1,
				DueDate: dueDate,
			},
			customerID: 1,
			schedules: map[models.InvoiceReminderSchedule]bool{
				models.InvoiceReminderSchedule14DaysBeforeDue: true,
			},
			mockSetup: func(expectedReminders []models.InvoiceReminder) {
				mockRepo.EXPECT().
					UpsertReminders(ctx, gomock.Any()).
					Return(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedReminders []models.InvoiceReminder
			for schedule, enabled := range tt.schedules {
				reminder := models.InvoiceReminder{
					InvoiceID:    tt.invoice.ID,
					CustomerID:   tt.customerID,
					ReminderDate: getReminderDateFromSchedule(schedule, tt.invoice.DueDate),
					Schedule:     schedule,
				}
				if !enabled {
					reminder.DeletedAt = helper.ReturnPointer(time.Now())
				}
				expectedReminders = append(expectedReminders, reminder)
			}

			tt.mockSetup(expectedReminders)

			err := service.SetInvoiceReminders(ctx, tt.invoice, tt.customerID, tt.schedules)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewReminderService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository_mocks.NewMockReminderRepository(ctrl)

	service := NewReminderService(mockRepo)

	assert.NotNil(t, service)

	// Type assertion to verify the concrete type
	reminderSvc, ok := service.(*reminderService)
	assert.True(t, ok, "service should be of type *reminderService")
	assert.Equal(t, mockRepo, reminderSvc.reminderRepository)
}
