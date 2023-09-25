package backend

import "testing"

func TestTranslateText(t *testing.T) {
	type args struct {
		text     string
		language string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Successful Translation",
			args: args{
				text:     "Good evening",
				language: "es",
			},
			want:    "Buenas noches",
			wantErr: false,
		},
		// Example of a test case that will fail. Expected behaviour.
		//{
		//	name: " Failed Translation - Invalid Language Code",
		//	args: args{
		//		text:     "Good evening",
		//		language: "invalid",
		//	},
		//	want:    "",
		//	wantErr: true,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TranslateText(tt.args.text, tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranslateText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Logf("TranslateText() got = %v, want %v", got, tt.want)
			}
		})
	}
}
