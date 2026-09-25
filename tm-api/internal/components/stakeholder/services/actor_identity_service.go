package services

import (
	"context"
	"strings"

	coreModels "admin-panel-dashboard/internal/models"

	"gorm.io/gorm"
)

// loadOrganizationActors resolves organization-member identities in two
// bounded queries so list endpoints do not introduce N+1 database traffic.
func loadOrganizationActors(ctx context.Context, db *gorm.DB, memberIDs, organizationIDs []string) (map[string]coreModels.OrganizationMember, map[string]coreModels.Organization, error) {
	membersByID := make(map[string]coreModels.OrganizationMember)
	organizationsByID := make(map[string]coreModels.Organization)
	memberIDs = uniqueNonEmptyStrings(memberIDs)
	organizationIDs = uniqueNonEmptyStrings(organizationIDs)

	if len(memberIDs) > 0 {
		var members []coreModels.OrganizationMember
		if err := db.WithContext(ctx).Where("id IN ?", memberIDs).Find(&members).Error; err != nil {
			return nil, nil, err
		}
		for _, member := range members {
			membersByID[member.ID] = member
			organizationIDs = append(organizationIDs, member.OrganizationID)
		}
	}
	organizationIDs = uniqueNonEmptyStrings(organizationIDs)
	if len(organizationIDs) > 0 {
		var organizations []coreModels.Organization
		if err := db.WithContext(ctx).Where("id IN ?", organizationIDs).Find(&organizations).Error; err != nil {
			return nil, nil, err
		}
		for _, organization := range organizations {
			organizationsByID[organization.ID] = organization
		}
	}
	return membersByID, organizationsByID, nil
}

func uniqueNonEmptyStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func organizationMemberDisplayName(member coreModels.OrganizationMember) string {
	name := strings.TrimSpace(strings.TrimSpace(member.FirstName) + " " + strings.TrimSpace(member.LastName))
	if name == "" {
		return member.Email
	}
	return name
}
