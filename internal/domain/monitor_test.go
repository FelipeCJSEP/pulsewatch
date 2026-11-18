package domain

import (
	"testing"
)

func TestNewMonitor_Validation(t *testing.T) {
	// Deve falhar se o nome for vazio
	_, err := NewMonitor("", "https://api.exemplo.com", 10, 3)

	if err != ErrMonitorNameRequired {
		t.Fatalf("esperava o erro %v, mas obteve %v", ErrMonitorNameRequired, err)
	}

	// Deve falhar se URL for inválida
	_, err = NewMonitor("API Pagamentos", "api exemplo", 10, 3)

	if err != ErrMonitorInvalidURL {
		t.Fatalf("esperava o erro %v, mas obteve %v", ErrMonitorInvalidURL, err)
	}

	// Deve funcionar
	_, err = NewMonitor("API Pagamentos", "https://api.exemplo.com", 10, 3)

	if err != nil {
		t.Fatalf("Esperava sucesso, mas obteve o erro %v", err)
	}
}

func TestMonitor_FailureThresholdTransitionsToDown(t *testing.T) {
	monitor, _ := NewMonitor("API Pagamentos", "https://api.exemplo.com", 10, 3)
	var recordFailed bool
	i := 1

	for {
		recordFailed = monitor.RecordFailed()

		if i > 2 {
			break
		}

		if recordFailed || monitor.Status == StatusDown {
			t.Fatalf("não deveria cair na tentativa %d, mas retornou %v com status %s", i, recordFailed, monitor.Status)
		}

		i++
	}

	if !recordFailed || monitor.Status != StatusDown {
		t.Fatalf("esperava que o monitor caísse (DOWN) na 3ª falha, mas status atual é %s", monitor.Status)
	}

	monitor.RecordSuccess()

	if monitor.Status != StatusUp || monitor.ConsecutiveFailures > 0 {
		t.Fatalf("esperava status UP e 0 falhas após recuperação, mas obteve status %s e %d falhas", monitor.Status, monitor.ConsecutiveFailures)
	}
}
