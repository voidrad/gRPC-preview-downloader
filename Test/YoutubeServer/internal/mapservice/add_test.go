package mapservice

import (
	"reflect"
	"sync"
	"testing"
)

func Test_AddToMap(t *testing.T) {
	type args struct {
		m           map[string][]byte
		youtubeLink string
		preview     []byte
		mutex       *sync.Mutex
	}
	tests := []struct {
		name string
		args args
		want map[string][]byte
	}{
		{
			name: "nil slice",
			args: args{
				m:           map[string][]byte{"test": []byte{}},
				youtubeLink: "testing",
				preview:     nil,
				mutex:       &sync.Mutex{},
			},
			want: map[string][]byte{"test": []byte{},
				"testing": nil},
		},
		{
			name: "defaut",
			args: args{
				m:           map[string][]byte{"test": []byte{'f', 'g'}},
				youtubeLink: "testing",
				preview:     []byte{'q', 'o'},
				mutex:       &sync.Mutex{},
			},
			want: map[string][]byte{"test": []byte{'f', 'g'},
				"testing": []byte{'q', 'o'}},
		},
		{
			name: "dual",
			args: args{
				m:           map[string][]byte{"test": []byte{'f', 'g'}},
				youtubeLink: "test",
				preview:     []byte{'f', 'g'},
				mutex:       &sync.Mutex{},
			},
			want: map[string][]byte{"test": []byte{'f', 'g'}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddToMap(tt.args.m, tt.args.youtubeLink, tt.args.preview, tt.args.mutex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
