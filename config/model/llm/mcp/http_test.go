package mcp

import "testing"

func TestHttp_Validate(t *testing.T) {
	tests := []struct {
		name    string
		http    Http
		wantErr bool
	}{
		{
			name: "valid http",
			http: Http{
				Endpoint: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "valid http",
			http: Http{
				BaseUrl: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "empty endpoint",
			http: Http{
				Endpoint: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.http.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Http.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttp_GetTransport(t *testing.T) {

}
