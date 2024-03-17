package handler

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
)

func Test_validateRegistrationRequest(t *testing.T) {
	type args struct {
		request generated.RegistrationRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John Does",
					PhoneNumber: "+62811223344",
					Password:    "p4Ssword!",
				},
			},
			wantErr: false,
		},
		{
			name: "error - phone number less than 10",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John Does",
					PhoneNumber: "+6211223344",
					Password:    "p4Ssword!",
				},
			},
			wantErr: true,
		},
		{
			name: "error - phone number more than 13",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John Does",
					PhoneNumber: "+6211223344556677",
					Password:    "p4Ssword!",
				},
			},
			wantErr: true,
		},
		{
			name: "error - name less than 3",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "Ba",
					PhoneNumber: "+621122334455",
					Password:    "p4Ssword!",
				},
			},
			wantErr: true,
		},
		{
			name: "error - name more than 60",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John Imifdsmaofpdmsauifnodpsafundiopsanfdipsaomfmfduisaopmfdomfdosa",
					PhoneNumber: "+621122334455",
					Password:    "p4Ssword!",
				},
			},
			wantErr: true,
		},
		{
			name: "error - pass less than 6",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John",
					PhoneNumber: "+621122334455",
					Password:    "4Ssd!",
				},
			},
			wantErr: true,
		},
		{
			name: "error - pass more than 64",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John",
					PhoneNumber: "+621122334455",
					Password:    "4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!4Ssd!",
				},
			},
			wantErr: true,
		},
		{
			name: "error - pass not has upper",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John",
					PhoneNumber: "+621122334455",
					Password:    "fdasmfoda$##",
				},
			},
			wantErr: true,
		},
		{
			name: "error - pass not has special",
			args: args{
				request: generated.RegistrationRequest{
					FullName:    "John",
					PhoneNumber: "+621122334455",
					Password:    "fdasmfodanfdsaNLJH",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRegistrationRequest(tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("validateRegistrationRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hashingAndSalting(t *testing.T) {
	t.SkipNow()
	type args struct {
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal test",
			args: args{
				pass: "somepass",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashingAndSalting(tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashingAndSalting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("hashingAndSalting() = %v, want not empty", got)
			}
		})
	}
}

func Test_hashingSaltingAndMatchingPass(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal test",
			args: args{
				pass: "somepass123",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := hashingAndSalting(tt.args.pass)
			if err != nil {
				t.Errorf("hashingAndSalting() error = %v", err)
				return
			}
			if got := matchingPass(tt.args.pass, hashed); got != tt.want {
				t.Errorf("matchingPass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchingPass(t *testing.T) {
	type args struct {
		pass       string
		hashedPass string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal test",
			args: args{
				pass:       "aGsDfR!",
				hashedPass: "$2a$04$hq/dzj8jt4gjclRiL4vHj.j27JqNbQDZ5953YGzsbb/yH0nS6nf76",
			},
			want: true,
		},
		{
			name: "error test",
			args: args{
				pass:       "somewrongpass",
				hashedPass: "$2a$04$hq/dzj8jt4gjclRiL4vHj.j27JqNbQDZ5953YGzsbb/yH0nS6nf76",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchingPass(tt.args.pass, tt.args.hashedPass); got != tt.want {
				t.Errorf("matchingPass() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateJwtToken(t *testing.T) {
	type args struct {
		id     int
		secret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal test",
			args: args{
				id:     123,
				secret: "somesecret123",
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzfQ.5nJ-ejccbobBUr9v6_crlGAypd0RbmJNxg_99WQzlAw",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateJwtToken(tt.args.id, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJwtToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateJwtToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateToken(t *testing.T) {
	type args struct {
		tokenString string
		secret      string
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		{
			name: "normal test",
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzfQ.5nJ-ejccbobBUr9v6_crlGAypd0RbmJNxg_99WQzlAw",
				secret:      "somesecret123",
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "error - wrong jwt format",
			args: args{
				tokenString: "fdsafdsa",
				secret:      "somesecret123",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateToken(tt.args.tokenString, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
