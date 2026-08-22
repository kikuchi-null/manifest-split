package ms

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

// ---------- Args.validate() (正常系) ----------

func TestArgs_validate_success(t *testing.T) {

	tests := []struct {
		name     string
		in       Args
		wantMode string // 検証後に期待するMode(空指定はModeDefaultに補正されることを含めて確認)
	}{
		{
			name:     "modeが空文字ならdefaultに補正される",
			in:       Args{Mode: "", Input: "in.xml", Output: "out", Num: 1},
			wantMode: ModeDefault,
		},
		{
			name:     "defaultモードの下限境界値(1)",
			in:       Args{Mode: ModeDefault, Input: "in.xml", Output: "out", Num: 1},
			wantMode: ModeDefault,
		},
		{
			name:     "defaultモードの上限境界値(MemberLimit)",
			in:       Args{Mode: ModeDefault, Input: "in.xml", Output: "out", Num: MemberLimit},
			wantMode: ModeDefault,
		},
		{
			name:     "filesモードの下限境界値(1)",
			in:       Args{Mode: ModeFiles, Input: "in.xml", Output: "out", Num: 1},
			wantMode: ModeFiles,
		},
		{
			name:     "typesモードはNumを問わない",
			in:       Args{Mode: ModeTypes, Input: "in.xml", Output: "out"},
			wantMode: ModeTypes,
		},
		{
			name:     "sampleモードはInputを問わない",
			in:       Args{Mode: ModeSample, Output: "out"},
			wantMode: ModeSample,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.in
			a.validate() // ここでは正常系のみを扱うため os.Exit は発生しない
			if a.Mode != tt.wantMode {
				t.Errorf("Mode after validate() = %q, want %q", a.Mode, tt.wantMode)
			}
		})
	}
}

// ---------- Args.validate() (異常系) ----------
//
// validate()はエラー時にos.Exit(1)を呼ぶため、同一プロセス内では検証できない。
// テストバイナリ自身をサブプロセスとして再実行し、終了コードで判定する。

func TestArgs_validate_failure(t *testing.T) {

	if os.Getenv("MS_TEST_VALIDATE_SUBPROCESS") == "1" {
		n, _ := strconv.Atoi(os.Getenv("MS_TEST_NUM"))
		a := Args{
			Mode:   os.Getenv("MS_TEST_MODE"),
			Input:  os.Getenv("MS_TEST_INPUT"),
			Output: os.Getenv("MS_TEST_OUTPUT"),
			Num:    n,
		}
		a.validate()
		// 期待通りならvalidate()内のos.Exit(1)でここには到達しない
		os.Exit(0)
	}

	tests := []struct {
		name                string
		mode, input, output string
		num                 int
	}{
		{"不正なmode", "bogus", "in.xml", "out", 1},
		{"defaultモードでinput未指定", ModeDefault, "", "out", 1},
		{"defaultモードでoutput未指定", ModeDefault, "in.xml", "", 1},
		{"sampleモードでoutput未指定", ModeSample, "", "", 0},
		{"defaultモードでnum=0", ModeDefault, "in.xml", "out", 0},
		{"defaultモードでnum=MemberLimit+1", ModeDefault, "in.xml", "out", MemberLimit + 1},
		{"filesモードでnum=0", ModeFiles, "in.xml", "out", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestArgs_validate_failure")
			cmd.Env = append(os.Environ(),
				"MS_TEST_VALIDATE_SUBPROCESS=1",
				"MS_TEST_MODE="+tt.mode,
				"MS_TEST_INPUT="+tt.input,
				"MS_TEST_OUTPUT="+tt.output,
				"MS_TEST_NUM="+strconv.Itoa(tt.num),
			)

			err := cmd.Run()

			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				t.Fatalf("expected validate() to terminate the process with an error, got err=%v", err)
			}
			if exitErr.ExitCode() != 1 {
				t.Errorf("exit code = %d, want 1", exitErr.ExitCode())
			}
		})
	}
}

// ---------- isBatchMode / RecieveArgs(バッチ経路) ----------

// フラグはパッケージ変数のため、テスト間で汚染しないよう値を退避・復元する
func withFlags(t *testing.T, mode, input, output string, num int, clean bool) {
	t.Helper()

	origMode, origInput, origOutput, origNum, origClean := *fMode, *fInput, *fOutput, *fNum, *fClean
	*fMode, *fInput, *fOutput, *fNum, *fClean = mode, input, output, num, clean

	t.Cleanup(func() {
		*fMode, *fInput, *fOutput, *fNum, *fClean = origMode, origInput, origOutput, origNum, origClean
	})
}

func TestIsBatchMode(t *testing.T) {

	tests := []struct {
		name                string
		mode, input, output string
		num                 int
		want                bool
	}{
		{"フラグが何も指定されていなければ対話モード", "", "", "", 0, false},
		{"modeだけ指定されていればバッチモード", "default", "", "", 0, true},
		{"inputだけ指定されていればバッチモード", "", "in.xml", "", 0, true},
		{"outputだけ指定されていればバッチモード", "", "", "out", 0, true},
		{"numだけ指定されていればバッチモード", "", "", "", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withFlags(t, tt.mode, tt.input, tt.output, tt.num, false)
			if got := isBatchMode(); got != tt.want {
				t.Errorf("isBatchMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecieveArgs_batchMode(t *testing.T) {
	withFlags(t, ModeFiles, "in.xml", "out", 3, true)

	// isBatchMode()がtrueになるため、標準入力からのScanlnは発生しない
	a := RecieveArgs()

	want := Args{Mode: ModeFiles, Input: "in.xml", Output: "out", Num: 3, Clean: true}
	if a != want {
		t.Errorf("RecieveArgs() = %+v, want %+v", a, want)
	}
}
