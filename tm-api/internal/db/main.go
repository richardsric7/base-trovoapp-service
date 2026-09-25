package db

import (
	stakeholderModels "admin-panel-dashboard/internal/components/stakeholder/models"
	vaultSignerModels "admin-panel-dashboard/internal/components/vaultsigner/models"
	"admin-panel-dashboard/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"time"

	"github.com/ecnepsnai/discord"
	"gorm.io/driver/postgres"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var sqlDB *sql.DB // Set package-wide, but not exported
var onceAdmin sync.Once

// func OpenDb() (*gorm.DB, error) {
// 	dbType := os.Getenv("DB_TYPE")
// 	if len(dbType) == 0 {
// 		dbType = "postgres"
// 		log.Println("ENV DB_TYPE not set, using default postgres")
// 	}
// 	// if len(dbType) == 0 {
// 	// 	return nil, errors.New("empty db type. Check env variable: DB_TYPE")
// 	// }
// 	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
// 	if len(dbConnectionString) == 0 {
// 		return nil, errors.New("empty connection string. Check env variable: DB_CONNECTION_STRING")
// 	}

// 	var err error

// 	// if dbType == "sqlite" {
// 	// 	gormDB, err = gorm.Open(sqlite.Open(dbConnectionString), &gorm.Config{
// 	// 		QueryFields: true,
// 	// 	})
// 	// }

// 	if dbType == "postgres" {
// 		var maxoconn, idleCon string
// 		if os.Getenv("DB_MAX_OPEN_CONNECTIONS") != "" {
// 			maxoconn = os.Getenv("DB_MAX_OPEN_CONNECTIONS")
// 		} else {
// 			maxoconn = "50"
// 		}
// 		if os.Getenv("DB_MAX_IDLE_CONNECTIONS") != "" {
// 			idleCon = os.Getenv("DB_MAX_IDLE_CONNECTIONS")
// 		} else {
// 			idleCon = "50"
// 		}
// 		maxOpenConn, _ := strconv.ParseInt(maxoconn, 10, 64)
// 		maxIdleConns, _ := strconv.ParseInt(idleCon, 10, 64)
// 		// connsMaxIdleTime, _ := strconv.ParseInt(idletime, 10, 64)
// 		once.Do(func() {
// 			sqlDB, err = sql.Open("pgx", dbConnectionString)

// 			sqlDB.SetMaxIdleConns(int(maxIdleConns))
// 			sqlDB.SetConnMaxIdleTime(5 * time.Second)
// 			sqlDB.SetMaxOpenConns(int(maxOpenConn))
// 			sqlDB.SetConnMaxLifetime(24 * time.Hour)
// 		})

// 		gormDB, err = gorm.Open(postgres.New(postgres.Config{
// 			Conn: sqlDB,
// 		}), &gorm.Config{
// 			Logger:      logger.Default.LogMode(logger.Silent),
// 			QueryFields: true,
// 		})
// 		if err != nil {
// 			log.Printf("[OpenDb]failed to connect database, %s\n", err)
// 			return nil, err
// 		}

// 	}
// 	return gormDB, nil
// }

// TrovoWalletDb and P2PDb used to open separate connections (to
// WALLET_DB_CONNECTION_STRING / P2P_DB_CONNECTION_STRING respectively).
// The wallet, P2P, and admin schemas now all live in the same physical
// database, so AdminDB below is the only connection this service opens -
// see main.go, which uses it wherever these used to be called.

func AdminDB() (gormDB *gorm.DB, err error) {
	dbType := os.Getenv("DB_TYPE")
	if len(dbType) == 0 {
		dbType = "postgres"
		log.Println("ENV DB_TYPE not set, using default postgres")
	}
	// if len(dbType) == 0 {
	// 	return nil, errors.New("empty db type. Check env variable: DB_TYPE")
	// }
	dbConnectionString := os.Getenv("ADMIN_CONNECTION_STRING")
	if len(dbConnectionString) == 0 {
		return nil, errors.New("empty connection string. Check env variable: ADMIN_CONNECTION_STRING")
	}

	// var err error

	// if dbType == "sqlite" {
	// 	gormDB, err = gorm.Open(sqlite.Open(dbConnectionString), &gorm.Config{
	// 		QueryFields: true,
	// 	})
	// }

	if dbType == "postgres" {
		var maxoconn, idleCon string
		if os.Getenv("DB_MAX_OPEN_CONNECTIONS") != "" {
			maxoconn = os.Getenv("DB_MAX_OPEN_CONNECTIONS")
		} else {
			maxoconn = "50"
		}
		if os.Getenv("DB_MAX_IDLE_CONNECTIONS") != "" {
			idleCon = os.Getenv("DB_MAX_IDLE_CONNECTIONS")
		} else {
			idleCon = "50"
		}
		maxOpenConn, _ := strconv.ParseInt(maxoconn, 10, 64)
		maxIdleConns, _ := strconv.ParseInt(idleCon, 10, 64)
		// connsMaxIdleTime, _ := strconv.ParseInt(idletime, 10, 64)
		onceAdmin.Do(func() {
			sqlDB, err = sql.Open("pgx", dbConnectionString)

			sqlDB.SetMaxIdleConns(int(maxIdleConns))
			sqlDB.SetConnMaxIdleTime(5 * time.Second)
			sqlDB.SetMaxOpenConns(int(maxOpenConn))
			sqlDB.SetConnMaxLifetime(24 * time.Hour)
		})

		gormDB, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Silent),
			QueryFields: true,
		})
		if err != nil {
			log.Printf("[AdminDB]failed to connect database, %s\n", err)
			return nil, err
		}

	}
	if err = migrateAdminSchema(gormDB); err != nil {
		log.Println("admin db auto-migrate failed", err)
		return nil, err
	}

	// Seed permissions
	seedPermissions(gormDB)
	return gormDB, nil
}

func migrateAdminSchema(gormDB *gorm.DB) error {
	return gormDB.Transaction(func(tx *gorm.DB) error {
		if tx.Dialector.Name() == "postgres" {
			if err := tx.Exec("SELECT pg_advisory_xact_lock(20260821, 1)").Error; err != nil {
				return fmt.Errorf("acquire admin schema migration lock: %w", err)
			}
		}
		return migrateAdminSchemaTransaction(tx)
	})
}

func migrateAdminSchemaTransaction(gormDB *gorm.DB) error {
	if err := prepareAdminSchemaForAutoMigrate(gormDB); err != nil {
		return err
	}

	if err := gormDB.AutoMigrate(&models.AdminUser{}, &models.RoleConfig{}, &models.AdminPermission{},
		&models.RolePermission{}, &models.AdminSuspensionHistory{}, &models.UserSuspensionReason{},
		&models.AdminChangeLog{}, &models.SuspensionReason{}, &models.RoleConfig{}, &models.CurrencyConfig{},
		&models.AdminAccessLog{},
		&models.Organization{}, &models.OrganizationMember{}, &models.OrganizationInvite{},
		&stakeholderModels.StakeholderAssetAssignment{},
		&stakeholderModels.StakeholderAuditLog{},
		&stakeholderModels.StakeholderNotification{},
		&stakeholderModels.StakeholderAuthorizationChallenge{},
		&stakeholderModels.FundReleaseRequest{},
		&stakeholderModels.DueDiligenceChecklist{},
		&stakeholderModels.DueDiligenceItem{},
		&stakeholderModels.RevenueRecord{},
		&stakeholderModels.Distribution{},
		&stakeholderModels.AssetValuation{},
		&stakeholderModels.SegregatedAccount{},
		&stakeholderModels.ComplianceItem{},
		&stakeholderModels.StakeholderDocument{},
		&stakeholderModels.StakeholderNotificationPreference{},
		&stakeholderModels.StakeholderAssetOperation{},
		&stakeholderModels.StakeholderReport{},
		&stakeholderModels.StakeholderStructuringStatus{},
		&stakeholderModels.ComplianceRequirementTemplate{},
		&stakeholderModels.ComplianceRequirementInstance{},
		&vaultSignerModels.VaultSignerManagedSecret{},
		&vaultSignerModels.VaultSignerAssignment{},
		&vaultSignerModels.VaultSignerAuditLog{},
		// Legal Adviser and Financial Adviser are the distinct A5
		// stakeholder roles' partner tables - previously migrated by the
		// separate wallet-DB connection (TrovoWalletDb), now folded in
		// here since admin/wallet/p2p all share this one database.
		// AutoMigrate is additive and only creates the two tables (or
		// adds missing columns); it never drops or alters existing ones.
		&models.LegalAdviser{}, &models.FinancialAdviser{},
	); err != nil {
		return err
	}
	if err := pruneOrphanedStakeholderEvidenceLinks(gormDB); err != nil {
		return err
	}
	if err := gormDB.AutoMigrate(
		&stakeholderModels.ComplianceItemDocument{},
		&stakeholderModels.DueDiligenceItemDocument{},
		&stakeholderModels.ComplianceRequirementTemplateItem{},
	); err != nil {
		return err
	}

	return ensureAdminSchemaDefaults(gormDB)
}

// prepareAdminSchemaForAutoMigrate makes upgrades safe when an earlier release
// created the receiving-account columns without the final NOT NULL constraint.
// UpdateColumn avoids changing the business record's updated_at timestamp.
func prepareAdminSchemaForAutoMigrate(gormDB *gorm.DB) error {
	migrator := gormDB.Migrator()
	model := &stakeholderModels.FundReleaseRequest{}
	if migrator.HasTable(model) {
		for _, column := range []string{"receiving_bank", "receiving_account_name", "receiving_account_number"} {
			if !migrator.HasColumn(model, column) {
				continue
			}
			if err := gormDB.Model(model).Where(column+" IS NULL").UpdateColumn(column, "").Error; err != nil {
				return fmt.Errorf("backfill fund release %s: %w", column, err)
			}
		}
	}

	return nil
}

func pruneOrphanedStakeholderEvidenceLinks(gormDB *gorm.DB) error {
	migrator := gormDB.Migrator()
	type evidenceLinkMigration struct {
		model        interface{}
		parent       interface{}
		parentExpr   string
		documentExpr string
	}
	migrations := []evidenceLinkMigration{
		{
			model:        &stakeholderModels.ComplianceItemDocument{},
			parent:       &stakeholderModels.ComplianceItem{},
			parentExpr:   "NOT EXISTS (SELECT 1 FROM compliance_items WHERE compliance_items.id = compliance_item_documents.compliance_item_id)",
			documentExpr: "NOT EXISTS (SELECT 1 FROM stakeholder_documents WHERE stakeholder_documents.id = compliance_item_documents.document_id)",
		},
		{
			model:        &stakeholderModels.DueDiligenceItemDocument{},
			parent:       &stakeholderModels.DueDiligenceItem{},
			parentExpr:   "NOT EXISTS (SELECT 1 FROM due_diligence_items WHERE due_diligence_items.id = due_diligence_item_documents.due_diligence_item_id)",
			documentExpr: "NOT EXISTS (SELECT 1 FROM stakeholder_documents WHERE stakeholder_documents.id = due_diligence_item_documents.document_id)",
		},
	}

	for _, migration := range migrations {
		if !migrator.HasTable(migration.model) || !migrator.HasTable(migration.parent) ||
			!migrator.HasTable(&stakeholderModels.StakeholderDocument{}) {
			continue
		}
		query := gormDB.Where(migration.parentExpr).Or(migration.documentExpr)
		if err := query.Delete(migration.model).Error; err != nil {
			return fmt.Errorf("prune orphaned %T rows: %w", migration.model, err)
		}
	}
	return nil
}

// ensureAdminSchemaDefaults compensates for GORM v1.25 not detecting a missing
// empty-string default on an existing PostgreSQL column. New databases get the
// same defaults from the model tags; upgraded databases are corrected here.
func ensureAdminSchemaDefaults(gormDB *gorm.DB) error {
	migrator := gormDB.Migrator()
	model := &stakeholderModels.FundReleaseRequest{}
	if !migrator.HasTable(model) {
		return nil
	}
	columnTypes, err := migrator.ColumnTypes(model)
	if err != nil {
		return fmt.Errorf("inspect fund release column defaults: %w", err)
	}
	hasDefault := make(map[string]bool, len(columnTypes))
	for _, columnType := range columnTypes {
		_, hasDefault[columnType.Name()] = columnType.DefaultValue()
	}

	for _, column := range []string{"receiving_bank", "receiving_account_name", "receiving_account_number"} {
		if !migrator.HasColumn(model, column) {
			continue
		}
		if hasDefault[column] {
			continue
		}
		statement := fmt.Sprintf("ALTER TABLE fund_release_requests ALTER COLUMN %s SET DEFAULT ''", column)
		if err := gormDB.Exec(statement).Error; err != nil {
			return fmt.Errorf("set fund release %s default: %w", column, err)
		}
	}
	return nil
}

///:{currencies}
///:{role}

// func TrovoWalletDb1() (*gorm.DB, error) {
// 	// Get the connection string from the environment variable
// 	connectionString := os.Getenv("WALLET_DB_CONNECTION_STRING")
// 	if connectionString == "" {
// 		return nil, fmt.Errorf("WALLET_DB_CONNECTION_STRING environment variable is not set")
// 	}
// 	newLogger := logger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
// 		logger.Config{
// 			SlowThreshold:             time.Second, // Slow SQL threshold
// 			LogLevel:                  logger.Info, // Log level Info, Silent, Warn, Error
// 			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
// 			Colorful:                  false,       // Disable color
// 		},
// 	)
// 	gormConfig := &gorm.Config{
// 		Logger: newLogger,
// 	}

// 	//postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos",
// 	//	"localhost", "postgres", "toluwase", "trovowallet", 5432)
// 	// Open a connection to the database
// 	walletDb, err := gorm.Open(postgres.Open(connectionString), gormConfig)
// 	if err != nil {
// 		log.Println("wallet db connection error: ", err)
// 		return nil, err
// 	}
// 	log.Println("checking walletDb: ", walletDb)
// 	return walletDb, nil
// }

//	func P2PDb1() (*gorm.DB, error) {
//		// Get the connection string from the environment variable
//		connectionString := os.Getenv("P2P_DB_CONNECTION_STRING")
//		if connectionString == "" {
//			return nil, fmt.Errorf("DB_CONNECTION_STRING environment variable is not set")
//		}
//		newLogger := logger.New(
//			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
//			logger.Config{
//				SlowThreshold:             time.Second, // Slow SQL threshold
//				LogLevel:                  logger.Info, // Log level Info, Silent, Warn, Error
//				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
//				Colorful:                  false,       // Disable color
//			},
//		)
//		gormConfig := &gorm.Config{
//			Logger: newLogger,
//		}
//		// Open a connection to the database
//		log.Println("connectionString: ", connectionString)
//		p2pdb, err := gorm.Open(postgres.Open(connectionString), gormConfig)
//		if err != nil {
//			log.Println("p2pdb db connection error: ", err)
//			return nil, err
//		}
//		log.Println("checking p2pdb: ", p2pdb)
//		return p2pdb, nil
//	}
func PrintDBStats(tag string, db *gorm.DB) {
	sqlDB, _ := db.DB()
	dbStats := sqlDB.Stats()
	log.Printf("[%v], open connections: %v, Max Open Conns: %v, idle connections: %v, Max Idle Closed: %v", tag, dbStats.OpenConnections, dbStats.MaxOpenConnections, dbStats.Idle, dbStats.MaxIdleClosed)
	if dbStats.OpenConnections > (dbStats.MaxOpenConnections - 25) {
		logDiscordDBWarning(fmt.Sprintf("[DB CONNECTION FILLUP WARNING] open connections have reached %v/%v. increase limit!", dbStats.OpenConnections, dbStats.MaxOpenConnections))
	}
}

func logDiscordDBWarning(msg string) {
	discord.WebhookURL = "https://discord.com/api/webhooks/827986576415129663/wqMKp9wxB_fxs9Q3zlMKCNPGENXmD_ueUnL8hVCu1wmRfD2wkXAjfP85k1Ro_2_wGfiY"
	if len(os.Getenv("CONNECTION_WARNING_WEBHOOK")) > 50 {
		discord.WebhookURL = os.Getenv("CONNECTION_WARNING_WEBHOOK")
	}
	err := discord.Say(msg)
	if err != nil {
		log.Println("error sending discord message: ", err)
		return
	}
}

func seedPermissions(db *gorm.DB) {
	permissions := []models.AdminPermission{
		{Name: "VIEW_ADMIN_LIST"},
		{Name: "REMOVE_ADMIN_USER"},
		{Name: "GRANT_PERMISSIONS"},
		{Name: "INVITE_ADMIN"},
		{Name: "MANAGE_TROVO_WALLET"},
		{Name: "VERIFY_DOCUMENTS"},
		{Name: "UPDATE_PROFILE"},
		{Name: "ACCESS_REPORTS"},
		{Name: "MANAGE_SETTINGS"},
		{Name: "VIEW_USER_DATA"},
	}

	for _, permission := range permissions {
		db.FirstOrCreate(&permission, models.AdminPermission{Name: permission.Name})
	}

	rolePermissions := map[models.Role][]string{
		models.SuperAdmin: {
			"VIEW_ADMIN_LIST", "REMOVE_ADMIN_USER", "GRANT_PERMISSIONS", "INVITE_ADMIN",
			"MANAGE_TROVO_WALLET", "VERIFY_DOCUMENTS", "UPDATE_PROFILE", "ACCESS_REPORTS",
			"MANAGE_SETTINGS",
		},
		models.ViewOnlyAdmin: {"VIEW_ADMIN_LIST", "VIEW_USER_DATA"},
		models.EditLevelAdmin: {
			"VIEW_ADMIN_LIST", "VIEW_USER_DATA", "UPDATE_PROFILE", "ACCESS_REPORTS",
			"MANAGE_SETTINGS", "MANAGE_TROVO_WALLET", "VERIFY_DOCUMENTS",
		},
	}

	for role, permissions := range rolePermissions {
		for _, permName := range permissions {
			var permission models.AdminPermission
			db.First(&permission, "name = ?", permName)
			db.FirstOrCreate(&models.RolePermission{Role: role, PermissionID: permission.ID})
		}
	}
}
