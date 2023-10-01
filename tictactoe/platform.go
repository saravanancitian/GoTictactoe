//go:build ((darwin && !arm && !arm64) || freebsd || linux || windows || js) && !android && !ios

package tictactoe

func IsMobileBuild() bool {
	return false
}
