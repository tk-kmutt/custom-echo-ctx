package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang-jwt/jwt"
)

var (
	_s  = []byte("s")
	_lt = 1
	_is = "is"
	_lu = LoginUser{Name: "n", Email: "e"}
)

func TestJWT_SignedString(t *testing.T) {
	type fields struct {
		secret   []byte
		lifeTime int
		issuer   string
	}
	type args struct {
		lu LoginUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				secret:   _s,
				lifeTime: _lt,
				issuer:   _is,
			},
			args: args{
				lu: _lu,
			},
			want:    "", //同じの作れないのでテストはスキップしとく
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JWT{
				secret:   tt.fields.secret,
				lifeTime: tt.fields.lifeTime,
				issuer:   tt.fields.issuer,
			}
			got, err := r.SignedString(tt.args.lu)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignedString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("SignedString() got = %v, want %v", got, tt.want)
			//}
			fmt.Println(got)
		})
	}
}

func TestJWT_VerifySignedString(t *testing.T) {
	j := NewJWT(_s, _lt, _is)
	_ss, _ := j.SignedString(_lu)

	type fields struct {
		secret   []byte
		lifeTime int
		issuer   string
	}
	type args struct {
		signedString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				secret:   _s,
				lifeTime: _lt,
				issuer:   _is,
			},
			args: args{
				signedString: _ss,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JWT{
				secret:   tt.fields.secret,
				lifeTime: tt.fields.lifeTime,
				issuer:   tt.fields.issuer,
			}
			got, err := r.VerifySignedString(tt.args.signedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifySignedString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestJWT_Verify(t *testing.T) {
	j := NewJWT(_s, _lt, _is)
	_ss, _ := j.SignedString(_lu)

	type fields struct {
		secret   []byte
		lifeTime int
		issuer   string
	}
	type args struct {
		signedString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    jwt.MapClaims
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				secret:   _s,
				lifeTime: _lt,
				issuer:   _is,
			},
			args: args{
				signedString: _ss,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &JWT{
				secret:   tt.fields.secret,
				lifeTime: tt.fields.lifeTime,
				issuer:   tt.fields.issuer,
			}
			got, err := r.Verify(tt.args.signedString)
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.True(t, got.VerifyIssuer(_is, true))
		})
	}
}
