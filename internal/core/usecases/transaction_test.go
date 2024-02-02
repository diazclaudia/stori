package usecases

import (
	"reflect"
	"stori/internal/core/domain"
	"testing"
)

func Test_getCalculations(t *testing.T) {
	type args struct {
		info *[]domain.Transaction
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Response
		wantErr bool
	}{
		{
			name: "success calculous",
			args: args{
				info: &[]domain.Transaction{
					{
						Id:    "1",
						Date:  "7/15",
						Value: "+60.5",
					},
					{
						Id:    "2",
						Date:  "1/28",
						Value: "-10.3",
					},
					{
						Id:    "3",
						Date:  "8/2",
						Value: "-20.46",
					},
					{
						Id:    "4",
						Date:  "12/13",
						Value: "+10",
					},
				},
			},
			want: &domain.Response{
				TotalBalance: 39.74,
				NumberTransactions: map[string]int{
					"1":  1,
					"12": 1,
					"7":  1,
					"8":  1,
				},
				AverageDebit:  -15.38,
				AverageCredit: 35.25,
			},
			wantErr: false,
		},
		{
			name: "error calculous",
			args: args{
				info: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCalculations(tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCalculations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCalculations() got = %v, want %v", got, tt.want)
			}
		})
	}
}
