package models

import "testing"

func TestDashboardRoleForStakeholderType(t *testing.T) {
	tests := []struct {
		name            string
		stakeholderType string
		wantRole        string
		wantOK          bool
	}{
		{name: "trustee", stakeholderType: " trustees ", wantRole: DashboardRoleTrustee, wantOK: true},
		{name: "custodian", stakeholderType: "approved_asset_custodian", wantRole: DashboardRoleAssetCustodian, wantOK: true},
		{name: "manager", stakeholderType: "ASSET_MANAGER", wantRole: DashboardRoleAssetManager, wantOK: true},
		{name: "unsupported org type does not grant access", stakeholderType: "ASSET_CUSTODIAN", wantOK: false},
		{name: "empty", stakeholderType: "", wantOK: false},
		// These 4 types now have portal dashboard access (view assigned assets).
		{name: "legal adviser has portal access", stakeholderType: "legal_adviser", wantRole: DashboardRoleLegalAdviser, wantOK: true},
		{name: "financial adviser has portal access", stakeholderType: "FINANCIAL_ADVISER", wantRole: DashboardRoleFinancialAdviser, wantOK: true},
		{name: "issuing house has portal access", stakeholderType: "asset_issuing_house", wantRole: DashboardRoleIssuingHouse, wantOK: true},
		{name: "rating agency has portal access", stakeholderType: " rating_agency ", wantRole: DashboardRoleRatingAgency, wantOK: true},
		// legal_and_professionals is a DISTINCT type (not an adviser) and remains
		// portal-excluded in this scope.
		{name: "legal and professionals is portal-excluded", stakeholderType: "legal_and_professionals", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRole, gotOK := DashboardRoleForStakeholderType(tt.stakeholderType)
			if gotOK != tt.wantOK {
				t.Fatalf("ok = %v, want %v", gotOK, tt.wantOK)
			}
			if gotRole != tt.wantRole {
				t.Fatalf("role = %q, want %q", gotRole, tt.wantRole)
			}
		})
	}
}

func TestRoleConfigRegistry(t *testing.T) {
	// Exactly 7 portal roles are registered.
	if got := len(SupportedDashboardRoles()); got != 7 {
		t.Fatalf("SupportedDashboardRoles len = %d, want 7", got)
	}

	cases := []struct {
		role          string
		wantFKColumn  string
		wantAssignCol bool
		wantTier      int
	}{
		{DashboardRoleTrustee, AssetFKColumnTrustee, true, 1},
		{DashboardRoleAssetManager, AssetFKColumnAssetManager, true, 1},
		{DashboardRoleAssetCustodian, AssetFKColumnCustodian, true, 1},
		{DashboardRoleLegalAdviser, AssetFKColumnLegalAdviser, false, 1},
		{DashboardRoleFinancialAdviser, AssetFKColumnFinancial, false, 1},
		{DashboardRoleIssuingHouse, AssetFKColumnIssuingHouse, false, 2},
		{DashboardRoleRatingAgency, AssetFKColumnRatingAgency, false, 2},
	}
	for _, tc := range cases {
		t.Run(tc.role, func(t *testing.T) {
			rc, ok := RoleConfigForRole(tc.role)
			if !ok {
				t.Fatalf("RoleConfigForRole(%q) not found", tc.role)
			}
			if rc.AssetFKColumn != tc.wantFKColumn {
				t.Fatalf("AssetFKColumn = %q, want %q", rc.AssetFKColumn, tc.wantFKColumn)
			}
			if AssetFKColumnForRole(tc.role) != tc.wantFKColumn {
				t.Fatalf("AssetFKColumnForRole mismatch for %q", tc.role)
			}
			if rc.HasAssignmentCols != tc.wantAssignCol {
				t.Fatalf("HasAssignmentCols = %v, want %v", rc.HasAssignmentCols, tc.wantAssignCol)
			}
			if RoleHasAssignmentColumns(tc.role) != tc.wantAssignCol {
				t.Fatalf("RoleHasAssignmentColumns mismatch for %q", tc.role)
			}
			if rc.Tier != tc.wantTier {
				t.Fatalf("Tier = %d, want %d", rc.Tier, tc.wantTier)
			}
		})
	}

	if _, ok := RoleConfigForRole("legal_and_professionals"); ok {
		t.Fatal("legal_and_professionals must NOT be a portal role")
	}
}

func TestRoleCanAccessDocument(t *testing.T) {
	if !RoleCanAccessDocument(DashboardRoleTrustee, nil) {
		t.Fatal("empty access roles should allow all supported roles")
	}
	if !RoleCanAccessDocument(DashboardRoleAssetManager, []string{"asset_manager"}) {
		t.Fatal("matching role should be allowed")
	}
	if RoleCanAccessDocument(DashboardRoleAssetCustodian, []string{"trustee"}) {
		t.Fatal("non-matching role should be rejected")
	}
}

func TestStakeholderEnums(t *testing.T) {
	if got := AllowedReportTypes(); len(got) != 2 || got[0] != ReportTypeIncome || got[1] != ReportTypeOperational {
		t.Fatalf("report types = %v", got)
	}
	if got := AllowedAssetOperationalStatuses(); len(got) != 4 {
		t.Fatalf("operational statuses = %v", got)
	}
}
