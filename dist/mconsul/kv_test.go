package mconsul

import "testing"

func TestGetValue(t *testing.T) {
	type args struct {
		key string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test1", args{"service/coast/mysql/coast_mysql_host"}, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := DefaultDatacenter
			value, err := dc.GetValue(tt.args.key)
			if (err != nil) == tt.wantErr {
				t.Errorf("Datacenter.GetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(string(value))
		})
	}
}
