package main

import (
	"testing"
)

func TestReader_getMicroPuzzleElement(t *testing.T) {
	type args struct {
		name    string
		content string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test result by given params",
			args: args{
				name:    "test_name",
				content: "test:content",
			},
			want: "<micro-puzzle-element name=\"test_name\"><template>test:content</template></micro-puzzle-element>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reader{}
			if got := r.getMicroPuzzleElement(tt.args.name, tt.args.content); got != tt.want {
				t.Errorf("Reader.getMicroPuzzleElement() = %v, want %v", got, tt.want)
			}
		})
	}
}
