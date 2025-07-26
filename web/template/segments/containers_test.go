package segments

import (
	"context"
	"testing"

	"github.com/a-h/templ/generator/htmldiff"
)

func TestList(t *testing.T) {
	type args struct {
		values []*ContainerApp
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
			want: `<div class="grid grid-cols-1 items-center gap-y-4"></div>`,
		},
		{
			name: "empty",
			args: args{
				values: []*ContainerApp{},
			},
			want: `<div class="grid grid-cols-1 items-center gap-y-4"></div>`,
		},
		{
			name: "one",
			args: args{
				values: []*ContainerApp{
					{"one", "http://localhost:8080"},
				},
			},
			want: `
        <div class="grid grid-cols-1 items-center gap-y-4">
          <a href="http://localhost:8080" target="_blank">
						<div class="bg-sky-400 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
							<div class="flex items-baseline uppercase">
              	<p>o</p>
                <p class="text-sm">ne</p>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">http://localhost:8080</p>
            </div>
          </a>
        </div>
      `,
		},
		{
			name: "three",
			args: args{
				values: []*ContainerApp{
					{"one", "http://localhost:8080"},
					{"two", "http://localhost:8081"},
					{"three", "http://localhost:8082"},
				},
			},
			want: `
        <div class="grid grid-cols-1 items-center gap-y-4">
          <a href="http://localhost:8080" target="_blank">
						<div class="bg-sky-400 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
							<div class="flex items-baseline uppercase">
              	<p>o</p>
                <p class="text-sm">ne</p>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">http://localhost:8080</p>
            </div>
          </a>
          <a href="http://localhost:8081" target="_blank">
						<div class="bg-sky-400 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
							<div class="flex items-baseline uppercase">
              	<p>t</p>
                <p class="text-sm">wo</p>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">http://localhost:8081</p>
            </div>
          </a>
          <a href="http://localhost:8082" target="_blank">
						<div class="bg-sky-400 dark:bg-sky-800 rounded-lg px-4 py-2 selectable">
							<div class="flex items-baseline uppercase">
              	<p>t</p>
                <p class="text-sm">hree</p>
              </div>
              <p class="text-sm text-gray-600 dark:text-gray-400">http://localhost:8082</p>
            </div>
          </a>
        </div>
      `,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			component := Containers(test.args.values)
			_, diff, err := htmldiff.DiffCtx(context.Background(), component, test.want)
			if err != nil {
				t.Fatal(err)
			}
			if diff != "" {
				t.Error(diff)
			}
		})
	}
}
