package workwave

// Driver is a driver in WorkWave.
type Driver struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
