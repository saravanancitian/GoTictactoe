//go:build (android || ios || (darwin && arm) || (darwin && arm64)) && !js

package tictactoe

func IsMobileBuild() bool {
	return true
}
