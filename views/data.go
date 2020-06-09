package views

import (
	"profile.com/models"
)

const (
	// ErrLevelDanger is error color for bootstrap alert
	ErrLevelDanger = "danger"
)

// Alert defines the shape of the alert object
type Alert struct {
	Level   string
	Message string
}

// Data defines the shape of the page data
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

// SetAlert sets the Alert object on a data struct
func (d *Data) SetAlert(level string, err error) {
	d.Alert = &Alert{
		Level:   level,
		Message: err.Error(),
	}
}
