package user

import (
	"io"
	"os"
	"os/exec"

	"github.com/BBVA/kapow/internal/server/model"
)

func spawn(h *model.Handler, out io.Writer) error {
	cmd := exec.Command(h.Route.Entrypoint)
	if out != nil {
		cmd.Stdout = out
	}
	cmd.Env = append(os.Environ(), "KAPOW_URL=http://localhost:8081")
	cmd.Env = append(cmd.Env, "KAPOW_HANDLER_ID="+h.ID)

	err := cmd.Run()

	return err
}
