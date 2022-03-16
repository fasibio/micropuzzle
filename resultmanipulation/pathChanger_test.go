package resultmanipulation

import (
	"testing"
)

func TestChangePathOfRessources(t *testing.T) {
	type args struct {
		html   string
		prefix string
	}
	type test struct {
		name    string
		args    args
		want    string
		wantErr bool
	}
	tests := []test{
		{
			name: "no to change tag in html input equal output",
			args: args{
				html:   "<html><body><div></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div></div></body></html>",
			wantErr: false,
		},
		{
			name: "to change link tag in html",
			args: args{
				html:   "<html><body><div><link href='/lala'>test</link></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><link href=\"mock/lala\"/>test</div></body></html>",
			wantErr: false,
		},
		{
			name: "to change link tag in html",
			args: args{
				html:   "<html><body><div><link href='/lala'>test</link></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><link href=\"mock/lala\"/>test</div></body></html>",
			wantErr: false,
		},
		{
			name: "to change script tag in html",
			args: args{
				html:   "<html><body><div><script src='/lala'>test</script></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><script src=\"mock/lala\">test</script></div></body></html>",
			wantErr: false,
		},
		{
			name: "to change img tag in html",
			args: args{
				html:   "<html><body><div><img src='/lala'>test</img></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><img src=\"mock/lala\"/>test</div></body></html>",
			wantErr: false,
		},
		{
			name: "to change img tag in html",
			args: args{
				html:   "<html><body><div><img src='/lala'>test</img></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><img src=\"mock/lala\"/>test</div></body></html>",
			wantErr: false,
		},
		{
			name: "to change a tag in html",
			args: args{
				html:   "<html><body><div><a href='/lala'>test</a></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><a href=\"mock/lala\">test</a></div></body></html>",
			wantErr: false,
		},
		{
			name: "to change iframe tag in html",
			args: args{
				html:   "<html><body><div><iframe src='/lala'>test</iframe></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><iframe src=\"mock/lala\">test</iframe></div></body></html>",
			wantErr: false,
		},
		{
			name: "to change embed tag in html",
			args: args{
				html:   "<html><body><div><embed src='/lala'>test</embed></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><embed src=\"mock/lala\"/>test</div></body></html>",
			wantErr: false,
		},
		{
			name: "to change source tag in html",
			args: args{
				html:   "<html><body><div><source src='/lala'>test</source></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><source src=\"mock/lala\"/>test</div></body></html>",
			wantErr: false,
		},
		{
			name: "to change track tag in html",
			args: args{
				html:   "<html><body><div><track src='/lala'>test</track></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><track src=\"mock/lala\"/>test</div></body></html>",
			wantErr: false,
		},
		{
			name: "to change video tag in html",
			args: args{
				html:   "<html><body><div><video src='/lala'>test</video></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><video src=\"mock/lala\">test</video></div></body></html>",
			wantErr: false,
		},
		{
			name: "to change audio tag in html",
			args: args{
				html:   "<html><body><div><audio src='/lala'>test</audio></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><audio src=\"mock/lala\">test</audio></div></body></html>",
			wantErr: false,
		},
		{
			name: "to change form tag in html",
			args: args{
				html:   "<html><body><div><form action='/lala'>test</form></div></body></html>",
				prefix: "mock",
			},
			want:    "<html><head></head><body><div><form action=\"mock/lala\">test</form></div></body></html>",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChangePathOfRessources(tt.args.html, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangePathOfRessources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChangePathOfRessources() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangePathOfRessourcesCss(t *testing.T) {
	type args struct {
		css    string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no url to change in css",
			args: args{
				css:    "body { background-color: #f0f0f0; }",
				prefix: "mock",
			},
			want: "body { background-color: #f0f0f0; }",
		},
		{
			name: "url to change in css",
			args: args{
				css:    "body { background-color: #f0f0f0;\nbackground-image: url(/lala); }",
				prefix: "mock",
			},
			want: "body { background-color: #f0f0f0;\nbackground-image: url(mock/lala); }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChangePathOfRessourcesCss(tt.args.css, tt.args.prefix); got != tt.want {
				t.Errorf("ChangePathOfRessourcesCss() = %v, want %v", got, tt.want)
			}
		})
	}
}
