package backend

import "testing"

func TestGetTargetVoice(t *testing.T) {
	type args struct {
		language string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "English - Successful Match",
			args: args{
				language: "en",
			},
			want:    "Emma",
			wantErr: false,
		},
		{
			name: "Not supported language - Want Err",
			args: args{
				language: "xxxxxxxx",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTargetVoice(tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTargetVoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTargetVoice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
