package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/sony/sonyflake/v2"
)

type ShortKeyGenerator struct {
	sf *sonyflake.Sonyflake
}

var (
	sfGenerator *ShortKeyGenerator
	sfOnce      sync.Once
)

// getMachineID returns a unique machine ID
func getMachineID() (int, error) {
	// In production, you might want to derive this from:
	// - IP address
	// - MAC address
	// - Docker container ID
	// - Cloud instance ID
	return 1, nil
}

func CreateSnowFlakeInstance() {
	sfOnce.Do(func() {
		st := sonyflake.Settings{
			// Start time is important - choose a recent timestamp
			StartTime: time.Date(2024, 3, 28, 0, 0, 0, 0, time.UTC),
			// Optional: Provide a machine ID
			MachineID: getMachineID,
		}

		sf, err := sonyflake.New(st)
		if err != nil {
			panic("sonyflake couldn't be created")
		}

		sfGenerator = &ShortKeyGenerator{
			sf: sf,
		}
	})
}

func NewSnowflakeID() (int64, error) {
	// Idempotent via sync.Once; calling every time keeps the read of
	// sfGenerator ordered after its write, so there is no data race.
	CreateSnowFlakeInstance()

	id, err := sfGenerator.sf.NextID()
	if err != nil {
		return -1, fmt.Errorf("ERROR WHILE GENERATING NEW ID - %v", err)
	}

	return id, nil
}
