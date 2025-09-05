package cmdx

import (
	"context"
	"errors"
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
func Launch(ctx context.Context, workDir string, args []string) (*exec.Cmd, error) {
	if len(args) == 0 {
		return nil, errors.New("args cannot be empty")
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = workDir

	// 关键：只 Start，不 Wait
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// 注意：不能调 cmd.Wait()，否则又会阻塞

	return cmd, nil
}
