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
			want: `<div class="grid grid-cols-1 gap-4 px-4 py-4"></div>`,
		},
		{
			name: "empty",
			args: args{
				values: []ContainerApp{},
			},
			want: `<div class="grid grid-cols-1 gap-4 px-4 py-4"></div>`,
		},
		{
			name: "one",
			args: args{
				values: []ContainerApp{
					ContainerTestApp{"one", "http://localhost:8080"},
				},
			},
			want: `
        <div class="grid grid-cols-1 gap-4 px-4 py-4">
          <a href="http://localhost:8080" target="_blank">
            <div class="flex-auto bg-sky-200 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
              <p class=" text-sm font-medium text-gray-900 truncate dark:text-white">
                one
              </p>
              <p class="text-sm truncate text-gray-500 dark:text-gray-400">
                http://localhost:8080
              </p>
            </div>
          </a>
        </div>
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
        <div class="grid grid-cols-1 gap-4 px-4 py-4">
          <a href="http://localhost:8080" target="_blank">
            <div class="flex-auto bg-sky-200 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
              <p class=" text-sm font-medium text-gray-900 truncate dark:text-white">
                one
              </p>
              <p class="text-sm truncate text-gray-500 dark:text-gray-400">
                http://localhost:8080
              </p>
            </div>
          </a>
          <a href="http://localhost:8081" target="_blank">
            <div class="flex-auto bg-sky-200 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
              <p class=" text-sm font-medium text-gray-900 truncate dark:text-white">
                two
              </p>
              <p class="text-sm truncate text-gray-500 dark:text-gray-400">
                http://localhost:8081
              </p>
            </div>
          </a>
          <a href="http://localhost:8082" target="_blank">
            <div class="flex-auto bg-sky-200 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
              <p class=" text-sm font-medium text-gray-900 truncate dark:text-white">
                three
              </p>
              <p class="text-sm truncate text-gray-500 dark:text-gray-400">
                http://localhost:8082
              </p>
            </div>
          </a>
        </div>
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
