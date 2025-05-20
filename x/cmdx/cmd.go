package cmdx

import (
	"context"
	"os/exec"
)

func Run(ctx context.Context, workDir string, args []string) ([]byte, *exec.Cmd, error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, nil, err
	}
	return output, cmd, nil
}
