package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RunTapeInput struct {
	TapeText   *string `json:"tape_text,omitempty" jsonschema:"VHS tape text"`
	TapePath   *string `json:"tape_path,omitempty" jsonschema:"Path to a .tape file"`
	TimeoutSec int     `json:"timeout_sec,omitempty" jsonschema:"Timeout seconds (default 120)"`
}

type TuiSnapInput struct {
	Command        string   `json:"command" jsonschema:"Command to run (e.g., 'htop')"`
	WaitRegex      string   `json:"wait_regex" jsonschema:"Regex to wait for (VHS Wait+Screen /.../)"`
	Keys           []string `json:"keys,omitempty" jsonschema:"Extra keystrokes (e.g., [\"Down\",\"Enter\"])"`
	ScreenshotName string   `json:"screenshot_name,omitempty" jsonschema:"Output PNG filename (default screenshot.png)"`
	Width          int      `json:"width,omitempty" jsonschema:"Terminal width (default 1000)"`
	Height         int      `json:"height,omitempty" jsonschema:"Terminal height (default 600)"`
	TypingMS       int      `json:"typing_ms,omitempty" jsonschema:"Typing speed in ms (default 0)"`
	TimeoutSec     int      `json:"timeout_sec,omitempty" jsonschema:"Timeout seconds (default 120)"`
}

type ArtifactList struct {
	Artifacts []string `json:"artifacts"`
}

func runVHS(tapeText string, timeout time.Duration) ([]string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	wd, err := os.MkdirTemp(cacheDir, "vhs-mcp-*")
	if err != nil {
		return nil, err
	}
	tapePath := filepath.Join(wd, "demo.tape")
	if err := os.WriteFile(tapePath, []byte(tapeText), 0o644); err != nil {
		return nil, err
	}

	cmd := exec.Command("vhs", tapePath)
	cmd.Dir = wd
	// 環境によっては PATH を明示したい場合がある
	cmd.Env = os.Environ()

	// 標準出力/標準エラーをログへ（必要なら外してOK）
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("vhs run error: %w", err)
		}
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return nil, errors.New("vhs timed out")
	}

	var uris []string
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".png", ".gif", ".mp4", ".webm", ".txt", ".ascii", ".json":
			u := "file://" + filepath.ToSlash(path)
			uris = append(uris, u)
		}
		return nil
	})
	return uris, err
}

func toolRunTape(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in RunTapeInput,
) (*mcp.CallToolResult, ArtifactList, error) {
	if (in.TapeText == nil && in.TapePath == nil) || (in.TapeText != nil && in.TapePath != nil) {
		return nil, ArtifactList{}, errors.New("provide exactly one of tape_text or tape_path")
	}
	var tape string
	if in.TapeText != nil {
		tape = *in.TapeText
	} else {
		b, err := os.ReadFile(*in.TapePath)
		if err != nil {
			return nil, ArtifactList{}, err
		}
		tape = string(b)
	}
	timeout := 120
	if in.TimeoutSec > 0 {
		timeout = in.TimeoutSec
	}
	uris, err := runVHS(tape, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, ArtifactList{}, err
	}
	// テキストとしても一覧を返す
	text := "Produced:\n" + strings.Join(uris, "\n")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}, ArtifactList{Artifacts: uris}, nil
}

func toolTuiSnap(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in TuiSnapInput,
) (*mcp.CallToolResult, ArtifactList, error) {
	if in.Command == "" || in.WaitRegex == "" {
		return nil, ArtifactList{}, errors.New("command and wait_regex are required")
	}
	w := in.Width
	h := in.Height
	if w == 0 {
		w = 1000
	}
	if h == 0 {
		h = 600
	}
	png := in.ScreenshotName
	if png == "" {
		png = "screenshot.png"
	}
	typing := in.TypingMS
	var b strings.Builder
	fmt.Fprintf(&b, "Set Width %d\n", w)
	fmt.Fprintf(&b, "Set Height %d\n", h)
	fmt.Fprintf(&b, "Set TypingSpeed %d ms\n", typing)
	// コマンド起動
	fmt.Fprintf(&b, "Type \"%s\"\n", strings.ReplaceAll(in.Command, "\"", "`"))
	fmt.Fprintf(&b, "Enter\n")
	// 追加キー
	for _, k := range in.Keys {
		switch k {
		case "Up",
			"Down",
			"Left",
			"Right",
			"Enter",
			"Tab",
			"Space",
			"Backspace",
			"PageUp",
			"PageDown",
			"Escape":
			b.WriteString(k + "\n")
		default:
			// 文字列としてタイプ
			fmt.Fprintf(&b, "Type \"%s\"\n", strings.ReplaceAll(k, "\"", "`"))
		}
	}
	// 待ち → スクショ
	// VHS: Wait[+Screen][+Line] /regex/ と Screenshot <path> が仕様。:contentReference[oaicite:3]{index=3}
	fmt.Fprintf(&b, "Wait+Screen /%s/\n", in.WaitRegex)
	fmt.Fprintf(&b, "Screenshot %s\n", png)

	timeout := 120
	if in.TimeoutSec > 0 {
		timeout = in.TimeoutSec
	}
	uris, err := runVHS(b.String(), time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, ArtifactList{}, err
	}
	text := "Screenshot:\n" + strings.Join(uris, "\n")
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, ArtifactList{Artifacts: uris}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "vhs-mcp-go",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "vhs_run_tape",
		Description: "Run a VHS tape and return produced file URIs (png/gif/etc).",
	}, toolRunTape)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "tui_snap_after",
		Description: "Run a command via VHS, wait for regex, take a PNG screenshot.",
	}, toolTuiSnap)

	// Stdio でホストに接続
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
