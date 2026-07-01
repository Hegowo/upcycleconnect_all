package database

import (
	"fmt"

	"gorm.io/gorm"
)

func PurgeUser(db *gorm.DB, userID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var prestationIDs []uint
		tx.Table("prestations").Where("provider_id = ?", userID).Pluck("id", &prestationIDs)
		if len(prestationIDs) > 0 {
			tx.Exec("DELETE FROM prestation_images WHERE prestation_id IN ?", prestationIDs)
		}

		var projectIDs []uint
		tx.Table("upcycling_projects").Where("user_id = ?", userID).Pluck("id", &projectIDs)
		if len(projectIDs) > 0 {
			tx.Exec("DELETE FROM project_step_images WHERE step_id IN (SELECT id FROM project_steps WHERE project_id IN ?)", projectIDs)
			tx.Exec("DELETE FROM project_images WHERE project_id IN ?", projectIDs)
			tx.Exec("DELETE FROM project_steps WHERE project_id IN ?", projectIDs)
		}

		var campaignIDs []uint
		tx.Table("campaigns").Where("provider_id = ?", userID).Pluck("id", &campaignIDs)
		if len(campaignIDs) > 0 {
			tx.Exec("DELETE FROM campaign_items WHERE campaign_id IN ?", campaignIDs)
		}

		deletes := []struct{ table, col string }{
			{"reservation_messages", "sender_id"},
			{"reservations", "user_id"},
			{"event_registrations", "user_id"},
			{"event_messages", "user_id"},
			{"deposit_requests", "user_id"},
			{"deposit_purchases", "provider_id"},
			{"deposit_favorites", "provider_id"},
			{"forum_reports", "reporter_id"},
			{"forum_replies", "user_id"},
			{"forum_threads", "user_id"},
			{"forum_bans", "user_id"},
			{"invoices", "user_id"},
			{"contracts", "user_id"},
			{"subscriptions", "user_id"},
			{"campaigns", "provider_id"},
			{"upcycling_projects", "user_id"},
			{"prestations", "provider_id"},
			{"notifications", "user_id"},
			{"passkeys", "user_id"},
			{"email_change_requests", "user_id"},
			{"email_verifications", "user_id"},
			{"user_known_ips", "user_id"},
			{"ticket_messages", "sender_id"},
			{"tickets", "user_id"},
			{"staff_shifts", "employee_id"},
			{"staff_event_members", "user_id"},
			{"provider_profiles", "user_id"},
			{"user_roles", "user_id"},
		}
		for _, d := range deletes {
			tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = ?", d.table, d.col), userID)
		}

		tx.Exec("UPDATE audit_logs SET user_id = NULL WHERE user_id = ?", userID)
		tx.Exec("UPDATE events SET created_by = NULL WHERE created_by = ?", userID)

		if err := tx.Exec("DELETE FROM users WHERE id = ?", userID).Error; err != nil {
			return err
		}
		return nil
	})
}
