# wzmux

[English](README.en.md)

Claude Code 에이전트를 WezTerm 탭으로 실행, 모니터링, 관리하는 멀티플렉서.

## 주요 기능

- **멀티 에이전트** — 여러 Claude Code 세션을 WezTerm 탭에서 동시 관리
- **실시간 대시보드** — Bubble Tea 기반 TUI로 에이전트 상태 모니터링
- **세션 재개** — 이름으로 이전 대화를 이어서 진행
- **자동 상태 표시** — 탭 타이틀에 상태 이모지 자동 반영 (⚙ 실행중, 🟡 대기, ✅ 완료, 🔴 에러)
- **단일 바이너리** — WezTerm 외 런타임 의존성 없음

## 요구사항

- [WezTerm](https://wezfurlong.org/wezterm/)
- [Claude Code](https://docs.anthropic.com/en/docs/claude-code)

## 설치

```bash
go install github.com/gnu-gnu/wzmux@latest
```

소스에서 빌드:

```bash
git clone https://github.com/gnu-gnu/wzmux.git
cd wzmux
go install .
```

## 셋업

Claude Code hook 등록 (최초 1회):

```bash
wzmux setup
```

`~/.claude/settings.json`에 hook을 추가한다. 기존 설정은 보존되며 `settings.json.bak`으로 백업된다.

## 사용법

### 새 에이전트 시작

```bash
wzmux new backend
wzmux new frontend "로그인 페이지 리팩토링해줘"
```

WezTerm 탭이 생성되고 오른쪽에 대시보드가 split된다. 이름이 중복되면 자동으로 suffix가 붙는다 (`backend`, `backend-1`, `backend-2`, ...).

### 세션 재개

```bash
wzmux resume backend
```

이전 대화를 이어서 진행한다.

### 에이전트 목록

```bash
wzmux ls
```

### 대시보드

```bash
wzmux dashboard
```

`1`-`9` 키로 에이전트 탭 이동, `q`로 종료.

### 에이전트 종료

```bash
wzmux kill backend
wzmux kill --all
```

## 동작 원리

```
Claude Code 에이전트
       ↓ (hook 이벤트)
wzmux hook <event>          ← Go 바이너리가 직접 처리
       ↓
/tmp/claude-agent-status/  ← 상태 JSON 파일
       ↑
wzmux dashboard             ← 상태 파일을 읽어 TUI 렌더링
```

1. `wzmux setup`이 `wzmux hook <event>`를 Claude Code hook으로 등록
2. WezTerm 안에서 Claude Code가 실행되면 hook이 탭 타이틀 변경 + 상태 파일 기록
3. 대시보드가 2초마다 상태 파일을 폴링
4. WezTerm 밖에서는 hook이 아무 동작도 하지 않음

## 제거

```bash
wzmux uninstall
```

`~/.claude/settings.json`에서 wzmux hook만 제거하고 상태 파일을 정리한다. 다른 설정은 건드리지 않는다.

## 라이선스

MIT
