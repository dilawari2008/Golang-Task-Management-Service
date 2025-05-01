package consumers

import (
	"task-management-system/models"
)

type Consumer interface {
	Process(task *models.Task) (bool, error)
}

