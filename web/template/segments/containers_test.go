package segments

import (
	"context"
	"testing"

	"github.com/a-h/templ/generator/htmldiff"
)

func TestList(t *testing.T) {
	type args struct {
		values []ContainerApp
	}
	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "nil",
			args: args{
				values: nil,
			},
			want: `<ul class="list-inside list-disc"></ul>`,
		},
		{
			name: "empty",
			args: args{
				values: []ContainerApp{},
			},
			want: `<ul class="list-inside list-disc"></ul>`,
		},
		{
			name: "one",
			args: args{
				values: []ContainerApp{
					ContainerTestApp{"one", "http://localhost:8080"},
				},
			},
			want: `
        <ul class="list-inside list-disc">
          <li>
            <a href="http://localhost:8080" target="_blank">
              one
            </a>
          </li>
        </ul>
      `,
		},
		{
			name: "three",
			args: args{
				values: []ContainerApp{
					ContainerTestApp{"one", "http://localhost:8080"},
					ContainerTestApp{"two", "http://localhost:8081"},
					ContainerTestApp{"three", "http://localhost:8082"},
				},
			},
			want: `
        <ul class="list-inside list-disc">
          <li>
            <a href="http://localhost:8080" target="_blank">
              one
            </a>
          </li>
          <li>
            <a href="http://localhost:8081" target="_blank">
              two
            </a>
          </li>
          <li>
            <a href="http://localhost:8082" target="_blank">
              three
            </a>
          </li>
        </ul>
      `,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			component := List(test.args.values)
			diff, err := htmldiff.DiffCtx(context.Background(), component, test.want)
			if err != nil {
				t.Fatal(err)
			}
			if diff != "" {
				t.Error(diff)
			}
		})
	}
}

type ContainerTestApp struct {
	name string
	url  string
}

func (c ContainerTestApp) Name() string { return c.name }
func (c ContainerTestApp) URL() string  { return c.url }
