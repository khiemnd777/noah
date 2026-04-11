package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/khiemnd777/noah_api/modules/i18n/repository"
	"github.com/khiemnd777/noah_api/modules/i18n/service"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/utils"
)

func SyncLanguageXMLOnBoot(db *generated.Client) error {
	if db == nil {
		return nil
	}

	xmlDir := utils.GetFullPath("languages")
	log.Printf("🌐 Syncing language XML from %s", xmlDir)

	svc := service.NewLanguageService(repository.NewLanguageRepository(db))
	if err := svc.SyncFromDirectory(context.Background(), xmlDir); err != nil {
		return fmt.Errorf("sync language xml on boot failed: %w", err)
	}

	log.Printf("✅ Language XML synced successfully")
	return nil
}
