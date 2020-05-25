package views

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
	Yield interface{}
}

// SetAlert sets the Alert object on a data struct
func (d *Data) SetAlert(data *Data, level string, err error) {
	data.Alert = &Alert{
		Level:   level,
		Message: err.Error(),
	}
}