package save

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Beeketing/consul/agent"
	"github.com/Beeketing/consul/testutil"
	"github.com/mitchellh/cli"
)

func TestSnapshotSaveCommand_noTabs(t *testing.T) {
	t.Parallel()
	if strings.ContainsRune(New(cli.NewMockUi()).Help(), '\t') {
		t.Fatal("help has tabs")
	}
}
func TestSnapshotSaveCommand_Validation(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		args   []string
		output string
	}{
		"no file": {
			[]string{},
			"Missing FILE argument",
		},
		"extra args": {
			[]string{"foo", "bar", "baz"},
			"Too many arguments",
		},
	}

	for name, tc := range cases {
		ui := cli.NewMockUi()
		c := New(ui)

		// Ensure our buffer is always clear
		if ui.ErrorWriter != nil {
			ui.ErrorWriter.Reset()
		}
		if ui.OutputWriter != nil {
			ui.OutputWriter.Reset()
		}

		code := c.Run(tc.args)
		if code == 0 {
			t.Errorf("%s: expected non-zero exit", name)
		}

		output := ui.ErrorWriter.String()
		if !strings.Contains(output, tc.output) {
			t.Errorf("%s: expected %q to contain %q", name, output, tc.output)
		}
	}
}

func TestSnapshotSaveCommand(t *testing.T) {
	t.Parallel()
	a := agent.NewTestAgent(t.Name(), ``)
	defer a.Shutdown()
	client := a.Client()

	ui := cli.NewMockUi()
	c := New(ui)

	dir := testutil.TempDir(t, "snapshot")
	defer os.RemoveAll(dir)

	file := path.Join(dir, "backup.tgz")
	args := []string{
		"-http-addr=" + a.HTTPAddr(),
		file,
	}

	code := c.Run(args)
	if code != 0 {
		t.Fatalf("bad: %d. %#v", code, ui.ErrorWriter.String())
	}

	f, err := os.Open(file)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer f.Close()

	if err := client.Snapshot().Restore(nil, f); err != nil {
		t.Fatalf("err: %v", err)
	}
}
