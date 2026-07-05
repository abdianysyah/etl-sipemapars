package reporter

type Reporter interface {
	Progress(jobUUID string, progress int, currentStep string) error
	Log(jobUUID string, level int, step string, message string) error
	Finish(jobUUID string, message string) error
	Failed(jobUUID string, step string, message string) error
}

type Noop struct{}

func (Noop) Progress(string, int, string) error { return nil }
func (Noop) Log(string, int, string, string) error { return nil } 
func (Noop) Finish(string, string) error { return nil } 
func (Noop) Failed(string, string, string) error { return nil } 