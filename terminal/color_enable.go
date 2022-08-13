//go:build !windows

package terminal

func enableColor() (Cleanup, error) {
	return nil, nil
}
