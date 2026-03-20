package agent

import (
"bytes"
"fmt"
"os/exec"
"strings"
)

type SymbolEngine struct {
binaryPath string
}

func NewSymbolEngine() *SymbolEngine {
return &SymbolEngine{binaryPath: "pil"}
}

func (e *SymbolEngine) Evaluate(scriptPath string, args ...string) (string, error) {
cmdArgs := append([]string{scriptPath, "+"}, args...) // + to exit after execution if needed or handle scripts
cmd := exec.Command(e.binaryPath, cmdArgs...)

var out bytes.Buffer
var stderr bytes.Buffer
cmd.Stdout = &out
cmd.Stderr = &stderr

err := cmd.Run()
if err != nil {
 "", fmt.Errorf("picolisp error: %w (stderr: %s)", err, stderr.String())
}

return strings.TrimSpace(out.String()), nil
}
