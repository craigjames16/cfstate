package utils

import (
	"testing"
)

func TestGetHash(t *testing.T) {
	type args struct {
		fileData []byte
	}
	tests := []struct {
		name         string
		args         args
		wantFileHash string
		wantErr      bool
	}{
		{
			name: "TestGetHash",
			args: args{
				fileData: []byte("test"),
			},
			wantFileHash: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileHash, err := GetHash(tt.args.fileData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotFileHash != tt.wantFileHash {
				t.Errorf("GetHash() = %v, want %v", gotFileHash, tt.wantFileHash)
			}
		})
	}
}
