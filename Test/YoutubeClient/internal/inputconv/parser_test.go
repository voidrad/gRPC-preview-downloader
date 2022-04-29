package inputconv

import (
	"reflect"
	"testing"
)

func Test_AddToMap(t *testing.T) {

	type want struct {
		links []string
		async bool
	}

	tests := []struct {
		name string
		args []string
		want want
	}{
		{
			name: "bad link in raw",
			args: []string{"www.NOTyoutube.com/watch?v=vZwSem7EpZ0", "www.youtube.com/watch?v=MYVlU4HFepk"},
			want: want{
				links: []string{"MYVlU4HFepk"},
				async: false,
			},
		},
		{
			name: "standart whithout async",
			args: []string{"www.youtubge.com/watch?v=vZwSem7EpZ0"},
			want: want{
				links: []string{},
				async: false,
			},
		},
		{
			name: "standart whith async",
			args: []string{"www.youtube.com/watch?v=vZwSem7EpZ0", "www.youtube.com/watch?v=MYVlU4HFepk", "--async"},
			want: want{
				links: []string{"vZwSem7EpZ0", "MYVlU4HFepk"},
				async: true,
			},
		},
		{
			name: "async in the midle",
			args: []string{"www.youtube.com/watch?v=vZwSem7EpZ0", "--async", "www.youtube.com/watch?v=MYVlU4HFepk"},
			want: want{
				links: []string{"vZwSem7EpZ0", "MYVlU4HFepk"},
				async: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := want{}
			if got.links, got.async = ParseCommand(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
