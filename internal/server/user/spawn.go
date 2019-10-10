package user

import (
	"io"
	"os"
	"os/exec"

	"github.com/google/shlex"

	"github.com/BBVA/kapow/internal/server/model"
)

func spawn(h *model.Handler, out io.Writer) error {
	args, err := shlex.Split(h.Route.Entrypoint)
	if err != nil {
		return err
	}

	if h.Route.Command != "" {
		args = append(args, h.Route.Command)
	}

	cmd := exec.Command(args[0], args[1:]...)
	if out != nil {
		cmd.Stdout = out
	}
	cmd.Env = append(os.Environ(), "KAPOW_URL=http://localhost:8081")
	cmd.Env = append(cmd.Env, "KAPOW_HANDLER_ID="+h.ID)

	err = cmd.Run()

	return err
}
