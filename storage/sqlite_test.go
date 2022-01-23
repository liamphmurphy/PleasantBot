package storage

import "testing"

func TestDelete(t *testing.T) {
	tests := []struct {
		description string
		ID          string
		wantErr     error
	}{
		{
			description: "id should be found and deleted",
			ID:          "1",
			wantErr:     nil,
		},
	}
}
