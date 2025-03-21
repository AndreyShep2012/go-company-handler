package simple

import "log/slog"

type Publisher struct {
}

func New() *Publisher {
	return new(Publisher)
}

func (Publisher) OnCreateCompany(e any) {
	slog.Info("OnCreateCompany event received", "event", e)
}

func (Publisher) OnPatchCompany(e any) {
	slog.Info("OnPatchCompany event received", "event", e)
}

func (Publisher) OnDeleteCompany(e any) {
	slog.Info("OnDeleteCompany event received", "event", e)
}
