package monitor

import (
	"encoding/json"
	"time"

	"github.com/khiemnd777/noah_framework/shared/logger"
	"github.com/khiemnd777/noah_framework/shared/redis"
)

func SetModuleStatus(name string, port int, status string) {
	data := map[string]interface{}{
		"name":      name,
		"port":      port,
		"status":    status,
		"updatedAt": time.Now(),
	}

	jsonData, _ := json.Marshal(data)
	if err := redis.Set("status", "module_status:"+name, string(jsonData), 0); err != nil {
		logger.Warn("❌ Failed to write module status to Redis", err)
	}
}
