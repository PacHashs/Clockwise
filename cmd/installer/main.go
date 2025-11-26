package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows/registry"
)

func isAdmin() bool {
	var sid *windows.SID
	_, err := windows.OpenCurrentProcessToken()
	return err == nil
}

func runElevated() error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(strings.Join(os.Args[1:], " "))

	var showCmd int32 = 1 // SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		return fmt.Errorf("failed to elevate privileges: %v", err)
	}
	return nil
}

func addToPath(dir string) error {
	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.QUERY_VALUE|registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()

	path, _, err := k.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("failed to read PATH: %v", err)
	}

	paths := strings.Split(path, ";")
	for _, p := range paths {
		if strings.EqualFold(p, dir) {
			return nil // Already in PATH
		}
	}

	if path != "" {
		path = dir + ";" + path
	} else {
		path = dir
	}

	if err := k.SetStringValue("Path", path); err != nil {
		return fmt.Errorf("failed to update PATH: %v", err)
	}

	sendMessageTimeout("Environment")
	return nil
}

func sendMessageTimeout(msg string) {
	hwnd := windows.HWND_BROADCAST
	msgPtr, _ := syscall.UTF16PtrFromString(msg)
	wparam := uintptr(0)
	lparam := uintptr(0)
	flags := uint32(0x0000) // SMTO_ABORTIFHUNG
	timeout := uint32(1000) // 1 second
	_, _, _ = syscall.Syscall6(
		procSendMessageTimeoutW.Addr(),
		6,
		uintptr(hwnd),
		uintptr(wmSettingChange),
		wparam,
		lparam,
		uintptr(flags),
		uintptr(unsafe.Pointer(&timeout)),
	)
}

var (
	moduser32                  = windows.NewLazyDLL("user32.dll")
	procSendMessageTimeoutW    = moduser32.NewProc("SendMessageTimeoutW")
	wmSettingChange    uint32 = 0x001A // WM_SETTINGCHANGE
)

func main() {
	if !isAdmin() {
		fmt.Println("This installer requires administrator privileges to modify system PATH.")
		if err := runElevated(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to request elevation: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	binPath := flag.String("binary", "", "Path to cw binary to install (default: will prompt)")
	nonInteractive := flag.Bool("y", false, "non-interactive (accept defaults)")
	bootstrap := flag.Bool("bootstrap", false, "build cw and cwfmt from source before installing")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
        if *bootstrap {
            fmt.Println("Bootstrapping: building cw and cwfmt from source...")
            // Target paths
            binDir := "bin"
            if err := os.MkdirAll(binDir, 0755); err != nil {
                fmt.Fprintln(os.Stderr, "failed to create bin dir:", err)
                os.Exit(1)
            }
            exeSuffix := ""
            if runtime.GOOS == "windows" {
                exeSuffix = ".exe"
            }
            cwOut := filepath.Join(binDir, "cw"+exeSuffix)
            cwfmtOut := filepath.Join(binDir, "cwfmt"+exeSuffix)
            // Build cw
            if err := exec.Command("go", "build", "-o", cwOut, "./cmd/cw").Run(); err != nil {
                fmt.Fprintln(os.Stderr, "failed to build cw:", err)
                os.Exit(1)
            }
            // Build cwfmt
            if err := exec.Command("go", "build", "-o", cwfmtOut, "./cmd/cwfmt").Run(); err != nil {
                fmt.Fprintln(os.Stderr, "failed to build cwfmt:", err)
                os.Exit(1)
            }
            *binPath = cwOut
            fmt.Println("Built cw at:", *binPath)
        }
        if *nonInteractive {
            fmt.Fprintln(os.Stderr, "binary path required in non-interactive mode")
            os.Exit(2)
        }
        fmt.Print("Path to cw binary to install (leave empty to look for ./cw or ./bin/cw/cw.exe): ")
        t, _ := reader.ReadString('\n')
        t = strings.TrimSpace(t)
        if t == "" {
            // try defaults
            try := []string{"./cw", "./bin/cw", "./bin/cw.exe", "cw.exe", "bin/cw.exe"}
            found := ""
            for _, p := range try {
                if _, err := os.Stat(p); err == nil {
                    found = p
                    break
                }
            }
            if found == "" {
                fmt.Fprintln(os.Stderr, "could not find cw binary; please provide path")
                os.Exit(1)
            }
            *binPath = found
        } else {
            *binPath = t
        }
    }

    if _, err := os.Stat(*binPath); err != nil {
        fmt.Fprintln(os.Stderr, "binary not found:", *binPath)
        os.Exit(1)
    }

    fmt.Println("Installing Clockwise using installer...")

    switch runtime.GOOS {
    case "windows":
        installWindows(*binPath, *nonInteractive)
    default:
        installUnix(*binPath, *nonInteractive)
    }
}

func installWindows(binPath string, nonInteractive bool) {
    programFiles := os.Getenv("ProgramFiles")
    if programFiles == "" {
        programFiles = `C:\\Program Files`
    }
    destDir := filepath.Join(programFiles, "Clockwise")
    if err := os.MkdirAll(destDir, 0755); err != nil {
        fmt.Fprintln(os.Stderr, "failed to create dest dir:", err)
        os.Exit(1)
    }
    dest := filepath.Join(destDir, "cw.exe")
    if err := copyFile(binPath, dest); err != nil {
        fmt.Fprintln(os.Stderr, "failed to copy binary:", err)
        os.Exit(1)
    }

    // Add to user PATH using setx
    // Read current user PATH
    curPath := os.Getenv("PATH")
    if !strings.Contains(strings.ToLower(curPath), strings.ToLower(destDir)) {
        if err := addToUserPathWindows(destDir); err != nil {
            fmt.Fprintln(os.Stderr, "warning: failed to update user PATH automatically (you may need to add it manually):", err)
        } else if !nonInteractive {
            fmt.Println("Updated user PATH to include:", destDir)
            fmt.Println("You may need to restart your shell or explorer for changes to take effect.")
        }
    }

    fmt.Println("Clockwise installed to:", dest)
}

// addToUserPathWindows attempts multiple ways to add destDir to the user's PATH.
// It first tries 'setx', then falls back to a PowerShell registry update using
// [Environment]::SetEnvironmentVariable('Path','...','User').
func addToUserPathWindows(destDir string) error {
    curPath := os.Getenv("PATH")
    newPath := curPath + ";" + destDir
    // Try setx first
    cmd := exec.Command("setx", "PATH", newPath)
    if out, err := cmd.CombinedOutput(); err == nil {
        _ = out
        return nil
    }

    // Fallback: use PowerShell to set the user environment variable directly.
    // Build a PowerShell one-liner that appends destDir if not present.
    // Use single quotes to avoid expansion and escape single quotes in destDir.
    esc := strings.ReplaceAll(destDir, "'", "''")
    ps := fmt.Sprintf("$old = [Environment]::GetEnvironmentVariable('Path','User'); if ($old -notlike '*%s*') { [Environment]::SetEnvironmentVariable('Path', $old + ';%s', 'User') }", esc, esc)
    psCmd := exec.Command("powershell", "-NoProfile", "-Command", ps)
    if out, err := psCmd.CombinedOutput(); err != nil {
        return fmt.Errorf("setx failed and powershell fallback failed: %v: %s", err, string(out))
    }
    return nil
}

func installUnix(binPath string, nonInteractive bool) {
    dest := "/usr/local/bin/cw"
    // If not root, ask for sudo
    if os.Geteuid() != 0 {
        fmt.Println("Elevated privileges are required to install to /usr/local/bin. Running via sudo...")
        cmd := exec.Command("sudo", "cp", binPath, dest)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            fmt.Fprintln(os.Stderr, "failed to copy via sudo:", err)
            os.Exit(1)
        }
        if err := exec.Command("sudo", "chmod", "+x", dest).Run(); err != nil {
            fmt.Fprintln(os.Stderr, "failed to chmod:", err)
        }
    } else {
        if err := copyFile(binPath, dest); err != nil {
            fmt.Fprintln(os.Stderr, "failed to copy binary:", err)
            os.Exit(1)
        }
        if err := os.Chmod(dest, 0755); err != nil {
            fmt.Fprintln(os.Stderr, "failed to chmod:", err)
        }
    }
    fmt.Println("Clockwise installed to:", dest)
    if !nonInteractive {
        fmt.Println("If '/usr/local/bin' is not on your PATH, add it in your shell profile.")
    }
}

func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer func() {
        _ = out.Close()
    }()
    if _, err := io.Copy(out, in); err != nil {
        return err
    }
    return out.Sync()
}
