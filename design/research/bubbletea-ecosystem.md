# Bubble Tea ecosystem — top 1000 Go repos by stars

Captured 2026-08-31 via:

```
gh search repos --language=go bubbletea --sort=stars --limit=1000 \
  --json fullName,stargazersCount,description,url
```

Full-text search for "bubbletea" across Go repos — **not** the GitHub dependents graph.
It catches projects that name bubbletea in their name/description/readme and misses ones that
merely import it. 1000 is GitHub search's hard result cap, so this is everything the query can
reach. Raw JSON alongside in `bubbletea-ecosystem.json`.

| # | Stars | Repo | Description |
| ---: | ---: | --- | --- |
| 1 | 44725 | [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | A powerful little TUI framework 🏗 |
| 2 | 23018 | [yorukot/superfile](https://github.com/yorukot/superfile) | Pretty fancy and modern terminal file manager |
| 3 | 12448 | [dlvhdr/gh-dash](https://github.com/dlvhdr/gh-dash) | A rich terminal UI for GitHub that doesn't break your flow. |
| 4 | 3584 | [Gaurav-Gosain/tuios](https://github.com/Gaurav-Gosain/tuios) | Terminal UI OS (Terminal Multiplexer) |
| 5 | 2534 | [go-musicfox/go-musicfox](https://github.com/go-musicfox/go-musicfox) | go-musicfox是用Go写的又一款网易云音乐命令行客户端，支持UnblockNeteaseMusic、各种音质级别、lastfm、MPRIS、MacOS交互响应（睡眠暂停、蓝牙耳机连接断开响应、菜单栏控制等）... |
| 6 | 2128 | [idursun/jjui](https://github.com/idursun/jjui) | jjui is a TUI designed for interacting with the Jujutsu version control system. |
| 7 | 1078 | [floatpane/matcha](https://github.com/floatpane/matcha) | A beautiful and functional email client for your terminal, built with Go and the charming Bubble Tea TUI library. Never leave your command line to check your inbox or send an email again! |
| 8 | 910 | [lrstanley/bubblezone](https://github.com/lrstanley/bubblezone) | helper utility for BubbleTea, allowing easy mouse event tracking |
| 9 | 869 | [ekkinox/yai](https://github.com/ekkinox/yai) | Your AI powered terminal assistant. |
| 10 | 833 | [0xjuanma/golazo](https://github.com/0xjuanma/golazo) | The beautiful game in your terminal. Minimal TUI app to keep up with live & recent football/soccer matches written in Go. |
| 11 | 788 | [NimbleMarkets/ntcharts](https://github.com/NimbleMarkets/ntcharts) | Nimble Terminal Charts for the Golang BubbleTea framework and your TUIs |
| 12 | 742 | [guyfedwards/nom](https://github.com/guyfedwards/nom) | RSS reader for the terminal |
| 13 | 682 | [Broderick-Westrope/tetrigo](https://github.com/Broderick-Westrope/tetrigo) | Play Tetris in your terminal. |
| 14 | 635 | [mistakenelf/fm](https://github.com/mistakenelf/fm) | A terminal based file manager |
| 15 | 622 | [neur0map/glazepkg](https://github.com/neur0map/glazepkg) | See all your installed packages in one place. (Official package manage for RyokuArch) |
| 16 | 618 | [xdagiz/xytz](https://github.com/xdagiz/xytz) | A Beautiful YouTube Downloader/Player TUI |
| 17 | 568 | [siddhantac/puffin](https://github.com/siddhantac/puffin) | A beautiful terminal dashboard for hledger 💰 |
| 18 | 536 | [bjarneo/ku](https://github.com/bjarneo/ku) | A fast, keyboard-driven Kubernetes TUI. Browse any resource, edit objects, follow logs, and shell into pods.  |
| 19 | 493 | [Bahaaio/pomo](https://github.com/Bahaaio/pomo) | Customizable TUI Pomodoro timer with ASCII art, progress bar, desktop notifications, and productivity statistics. |
| 20 | 487 | [TypicalAM/goread](https://github.com/TypicalAM/goread) | Beautiful program to read your RSS/Atom feeds right in the terminal! |
| 21 | 483 | [termkit/gama](https://github.com/termkit/gama) | Manage your GitHub Actions from Terminal with great UI 🧪 |
| 22 | 480 | [robinovitch61/wander](https://github.com/robinovitch61/wander) | A terminal app/TUI for HashiCorp Nomad |
| 23 | 448 | [cedricblondeau/world-cup-2022-cli-dashboard](https://github.com/cedricblondeau/world-cup-2022-cli-dashboard) | Watch live World Cup 2022 matches in your terminal. ⚽🏆 |
| 24 | 441 | [jonas-grgt/ktea](https://github.com/jonas-grgt/ktea) | Kafka TUI client |
| 25 | 426 | [paulilaaso/bit](https://github.com/paulilaaso/bit) | CLI / TUI Logo Designer + ANSI Font Library with Gradients, Shadows, and Multi-Format Export |
| 26 | 404 | [ssleert/zfxtop](https://github.com/ssleert/zfxtop) | [WIP] fetch top for gen Z with X written by bubbletea enjoyer |
| 27 | 401 | [jon4hz/fztea](https://github.com/jon4hz/fztea) | 🐬🧋  Remote control your flipper from the local terminal or over SSH |
| 28 | 398 | [kiki-ki/go-qo](https://github.com/kiki-ki/go-qo) | qo is an interactive minimalist TUI to query JSON and CSV using SQL. |
| 29 | 387 | [Nomadcxx/sysc-greet](https://github.com/Nomadcxx/sysc-greet) | A tui greeter (not built in rust)  |
| 30 | 378 | [YoanWai/agent-manager](https://github.com/YoanWai/agent-manager) | The fastest workflow for every AI coding agent. Live status, quick prompts, worktrees, and diff review from one tmux TUI. |
| 31 | 374 | [c-grimshaw/gosniff](https://github.com/c-grimshaw/gosniff) | A fancy-schmancy tcpdump-esque TUI, programmed in Go. |
| 32 | 366 | [brianstrauch/solitaire-tui](https://github.com/brianstrauch/solitaire-tui) | Klondike solitaire for the terminal |
| 33 | 361 | [xguot/difi](https://github.com/xguot/difi) | Review and refine Git diffs before you push |
| 34 | 358 | [textfuel/lazyjira](https://github.com/textfuel/lazyjira) | Lazygit but for Jira |
| 35 | 345 | [kerlenton/mcpsnoop](https://github.com/kerlenton/mcpsnoop) | Wireshark for MCP. A transparent proxy that shows every real tool call between your AI client and your MCP servers, live in your terminal. |
| 36 | 343 | [Achno/gocheat](https://github.com/Achno/gocheat) | A beautiful customizable TUI Cheatsheet for keybindings,hotkeys and more in the terminal  |
| 37 | 339 | [nore-dev/fman](https://github.com/nore-dev/fman) | TUI File Manager |
| 38 | 331 | [dhth/omm](https://github.com/dhth/omm) | on-my-mind: a keyboard-driven task manager for the command line |
| 39 | 323 | [babycommando/zuse](https://github.com/babycommando/zuse) | ZUSE is an irc client for the terminal made in Go with Bubbletea |
| 40 | 318 | [heymaikol/network-doctor](https://github.com/heymaikol/network-doctor) | Network Doctor is a cross-platform network troubleshooting TUI that turns interface, DNS, TCP, TLS, HTTP, proxy, and path-MTU checks into one plain-English diagnosis. |
| 41 | 310 | [rusinikita/trainer](https://github.com/rusinikita/trainer) | GoLang interview prep questions. Terminal app with Go challenges and learning links |
| 42 | 287 | [chmouel/lazyworktree](https://github.com/chmouel/lazyworktree) | Easy Git worktree management CLI and TUI for the terminal. |
| 43 | 270 | [rasjonell/dashbrew](https://github.com/rasjonell/dashbrew) | TUI dashboard builder that lets you visualize data from scripts and APIs right in your console |
| 44 | 268 | [eduardofuncao/squix](https://github.com/eduardofuncao/squix) | A CLI tool for managing and executing SQL queries across multiple databases. Written in Go, made beautiful with BubbleTea |
| 45 | 267 | [mistakenelf/teacup](https://github.com/mistakenelf/teacup) | A collection of bubbles and utilities for bubbletea applications |
| 46 | 259 | [emprcl/signls](https://github.com/emprcl/signls) | a non-linear, generative midi sequencer in the terminal :infinity: |
| 47 | 258 | [charmbracelet/bubbletea-app-template](https://github.com/charmbracelet/bubbletea-app-template) | A template repository to create Bubble Tea apps. |
| 48 | 245 | [kencx/keyb](https://github.com/kencx/keyb) | Create and view custom hotkey cheatsheets in the terminal |
| 49 | 243 | [Microck/moji](https://github.com/Microck/moji) | find and download fonts from the terminal |
| 50 | 242 | [ashish0kumar/typtea](https://github.com/ashish0kumar/typtea) | minimal terminal typing tester |
| 51 | 238 | [hedhyw/json-log-viewer](https://github.com/hedhyw/json-log-viewer) | Interactive viewer for JSON logs. |
| 52 | 238 | [programmersd21/kairo](https://github.com/programmersd21/kairo) | 🤩 Kairo is a fast, keyboard-first terminal task manager in Go 🐹 with offline-first SQLite, Git sync 🔁, fuzzy search 🔍 & Lua plugins 🧩 |
| 53 | 238 | [zhh2001/rote](https://github.com/zhh2001/rote) | A cron that remembers what it did—a CLI job scheduler with run history, captured output, and a live terminal dashboard. |
| 54 | 232 | [benbusby/colorstorm](https://github.com/benbusby/colorstorm) | An interactive TUI for creating color themes for Vim, VSCode, and Sublime |
| 55 | 220 | [gigagrug/schema](https://github.com/gigagrug/schema) | All in one CLI tool for the database \\| Migration, Studio, LSP |
| 56 | 219 | [junhoyeo/contrabass](https://github.com/junhoyeo/contrabass) | 🎸 A project-level orchestrator for AI coding agents — Go & Charm stack implementation of OpenAI's Symphony |
| 57 | 215 | [SayYoungMan/tfui](https://github.com/SayYoungMan/tfui) | Interactive TUI for performing Terraform workflows |
| 58 | 211 | [rshelekhov/lazymake](https://github.com/rshelekhov/lazymake) | Modern TUI for Makefiles with interactive target selection, dependency    visualization, and command safety analysis |
| 59 | 204 | [SourcewareLab/Toney](https://github.com/SourcewareLab/Toney) | Toney is a fast, lightweight, terminal-based note-taking app for the modern developer.  |
| 60 | 200 | [elewis787/boa](https://github.com/elewis787/boa) | A Cobra command styled usage and help component powered by bubbletea  |
| 61 | 198 | [Mayowa-Ojo/chmod-cli](https://github.com/Mayowa-Ojo/chmod-cli) | Effortlessly generate chmod commands |
| 62 | 198 | [context-labs/whip](https://github.com/context-labs/whip) | A fast coding-agent harness in Go. Tool-use loop, bubbletea TUI, provider-routable models with live catalog discovery, MCP support, background subagents. One binary, no runtime, built for open-source models. |
| 63 | 195 | [armandsauzay/note](https://github.com/armandsauzay/note) | ✍️ take notes in your terminal ✍️ |
| 64 | 194 | [backendsystems/nibble](https://github.com/backendsystems/nibble) | easy to use network scanner, with a clickable tui interface |
| 65 | 189 | [CamiloGarciaLaRotta/kboard](https://github.com/CamiloGarciaLaRotta/kboard) | Terminal game to practice keyboard typing |
| 66 | 184 | [dhth/prs](https://github.com/dhth/prs) | Stay updated on PRs from your terminal |
| 67 | 183 | [cruise-org/cruise](https://github.com/cruise-org/cruise) | Cruise is a powerful, intuitive, and fully-featured TUI (Terminal User Interface) for managing containers. Built with Go and Bubbletea, it offers a visually rich, keyboard-first experience for managing containers, images, volumes, networks, logs and more — all from your terminal. |
| 68 | 170 | [Gaurav-Gosain/gollama](https://github.com/Gaurav-Gosain/gollama) | Gollama: Your offline conversational AI companion. An interactive tool for generating creative responses from various models, right in your terminal. Ideal for brainstorming, creative writing, or seeking inspiration. |
| 69 | 169 | [llehouerou/waves](https://github.com/llehouerou/waves) | Terminal music player with Soulseek integration, MusicBrainz tagging, Last.fm scrobbling, and radio mode. Built with Go and Bubble Tea. |
| 70 | 166 | [LeperGnome/bt](https://github.com/LeperGnome/bt) | Interactive tree-like terminal file manager |
| 71 | 165 | [matteo-pacini/RadioGoGo](https://github.com/matteo-pacini/RadioGoGo) | 📻 Go-powered CLI to surf global radio waves via a sleek TUI. Tune in & let's Go 🚀! |
| 72 | 163 | [szktkfm/mdtt](https://github.com/szktkfm/mdtt) | 🗓️ Markdown Table Editor TUI |
| 73 | 155 | [Zebbeni/ansizalizer](https://github.com/Zebbeni/ansizalizer) | A TUI to convert Images to ANSI strings using bubbletea |
| 74 | 154 | [clawscli/claws](https://github.com/clawscli/claws) | k9s-inspired TUI for AWS resource management with vim-style navigation |
| 75 | 153 | [hjr265/gittop](https://github.com/hjr265/gittop) | A beautiful terminal UI for visualizing Git repository statistics, inspired by htop/btop. |
| 76 | 148 | [lrstanley/bubbletint](https://github.com/lrstanley/bubbletint) | Terminal tints for everyone |
| 77 | 147 | [Nomadcxx/sysc-Go](https://github.com/Nomadcxx/sysc-Go) | Terminal animation library for Go. Pure Go animations ready to use in your TUI applications. |
| 78 | 147 | [0xjuanma/helm](https://github.com/0xjuanma/helm) | A minimalistic and customizable pomodoro-like timer for your terminal written in Go. Includes default pomodoro timer and lets users define their own. |
| 79 | 144 | [chip/pathos](https://github.com/chip/pathos) | pathos - CLI for editing a PATH env variable |
| 80 | 143 | [pidanou/helm-tui](https://github.com/pidanou/helm-tui) | A simple terminal UI for Helm |
| 81 | 141 | [renatoworks/oh-my-reddit](https://github.com/renatoworks/oh-my-reddit) | Beautiful Reddit threads, live in your terminal |
| 82 | 141 | [abeleinin/goki](https://github.com/abeleinin/goki) | Anki-like flashcard management tool for the terminal! |
| 83 | 133 | [samyakbardiya/trex](https://github.com/samyakbardiya/trex) | A Terminal app for RegEx visualization, :t-rex: roar! |
| 84 | 127 | [rmhubbert/bubbletea-overlay](https://github.com/rmhubbert/bubbletea-overlay) | An overlay / modal window component for Charm's v1 Bubble Tea TUI framework.  |
| 85 | 126 | [joao-zanutto/easydocker](https://github.com/joao-zanutto/easydocker) | EasyDocker is a TUI focused on investigating and troubleshooting Docker resources. Highly inspired by lazydocker and k9s while leveraging beautiful graphics from BubbleTea |
| 86 | 122 | [sumant1122/perfdeck](https://github.com/sumant1122/perfdeck) | 📊 A modern, lightweight, and customizable TUI performance monitor for your terminal. Built with Go and Bubble Tea. |
| 87 | 121 | [gitsocial-org/gitsocial](https://github.com/gitsocial-org/gitsocial) | Cross-forge collaboration platform |
| 88 | 121 | [swadhinbiswas/veet](https://github.com/swadhinbiswas/veet) | Universal Linux Application Uninstaller & Deep-Clean Residual Purger |
| 89 | 117 | [Bharath-code/git-scope](https://github.com/Bharath-code/git-scope) | A fast TUI to see the status of all git repositories on your machine. |
| 90 | 108 | [BigJk/crt](https://github.com/BigJk/crt) | Minimal terminal emulator for Bubbletea. |
| 91 | 107 | [0xMassi/pocketdev](https://github.com/0xMassi/pocketdev) | One command to run the AI coding CLI you already pay for (Claude Code, Codex, Cursor, Gemini, Grok, aider) on a Tailscale-only Hetzner box — code from your phone. |
| 92 | 104 | [chronosphereio/calyptia-go-bubble-table](https://github.com/chronosphereio/calyptia-go-bubble-table) | TUI table component for Bubbletea. |
| 93 | 104 | [stigoleg/keep-alive](https://github.com/stigoleg/keep-alive) | Keep-Alive is a lightweight, cross-platform utility to prevent your system from sleeping. Perfect for uninterrupted downloads, active connections, or long-running tasks, like vibe coding. |
| 94 | 103 | [Rtarun3606k/TakaTime](https://github.com/Rtarun3606k/TakaTime) | TakaTime is a blazingly fast, privacy-first coding activity tracker and the open-source, self-hosted alternative to WakaTime.  Track your development activity across Neovim, VS Code, JetBrains IDEs, and Antigravity while generating rich analytics, beautiful dashboards, and GitHub profile statistics—all powered by your own MongoDB database. |
| 95 | 102 | [raiyanyahya/llmaker](https://github.com/raiyanyahya/llmaker) | Selfhost modern LLM stacks. Run the whole fleet from your terminal |
| 96 | 100 | [nnnkkk7/lazyactions](https://github.com/nnnkkk7/lazyactions) | Lazygit-style TUI for GitHub Actions — monitor, trigger, and manage workflows from your terminal |
| 97 | 99 | [BOTbkcd/mayhem](https://github.com/BOTbkcd/mayhem) | A minimal TUI based task tracker 📝 |
| 98 | 98 | [handlebargh/yatto](https://github.com/handlebargh/yatto) | Interactive version-controlled todo-list for the command-line |
| 99 | 98 | [savannahostrowski/ghost](https://github.com/savannahostrowski/ghost) | Ghost 👻 is an experimental CLI that uses AI to generate GitHub Actions workflows, using OpenAI |
| 100 | 96 | [epilande/codegrab](https://github.com/epilande/codegrab) | ✋ Interactive CLI tool for selecting and bundling code into a single, LLM-ready output file |
| 101 | 95 | [trashhalo/readcli](https://github.com/trashhalo/readcli) | Tool that lets you read website content on the command line |
| 102 | 93 | [mrusme/mercator](https://github.com/mrusme/mercator) | OpenStreetMap but as terminal user interface (TUI) program (https://tty.fail/mrus/mercator) |
| 103 | 90 | [knz/bubbline](https://github.com/knz/bubbline) | Line editor based on the Bubbletea library. |
| 104 | 90 | [deemkeen/stegodon](https://github.com/deemkeen/stegodon) | Blog directly from ssh! |
| 105 | 89 | [pivovarit/tdocker](https://github.com/pivovarit/tdocker) | minimalistic terminal UI for everyday Docker operations |
| 106 | 89 | [SakshhamTheCoder/adbt](https://github.com/SakshhamTheCoder/adbt) | A modern, keyboard-driven Android Debug Bridge TUI. |
| 107 | 86 | [nimblebun/wordle-cli](https://github.com/nimblebun/wordle-cli) | play wordle in your command line |
| 108 | 85 | [kostyay/httpmon](https://github.com/kostyay/httpmon) | Terminal-native HTTP/HTTPS debugging proxy with MITM intercept, ring buffer store, and Bubble Tea TUI |
| 109 | 84 | [dancnb/sonicradio](https://github.com/dancnb/sonicradio) | A TUI radio player making use of Radio Browser API and Bubbletea. |
| 110 | 84 | [MonsieurTib/service-bus-tui](https://github.com/MonsieurTib/service-bus-tui) | A TUI explorer for Azure Service Bus |
| 111 | 81 | [litescript/ls-horizons](https://github.com/litescript/ls-horizons) | Terminal UI for visualizing NASA's Deep Space Network in real-time |
| 112 | 79 | [lunemis/mux](https://github.com/lunemis/mux) | TUI tmux session manager with live preview — built for AI coding workflows |
| 113 | 78 | [allyring/pvw](https://github.com/allyring/pvw) | A port viewer TUI made with BubbleTea in Go |
| 114 | 78 | [stanlyzoolo/keepkit](https://github.com/stanlyzoolo/keepkit) | Lightweight TUI for tracking versions of your favorite tools |
| 115 | 77 | [emprcl/sektron](https://github.com/emprcl/sektron) | a midi step sequencer in the terminal, made with live performance in mind :loop: |
| 116 | 76 | [serene-brew/Kaizen](https://github.com/serene-brew/Kaizen) | Why Watch Anime Like a Normie? |
| 117 | 72 | [KevM/bubbleo](https://github.com/KevM/bubbleo) | BubbleO is a collection of components for the excellent terminal UI tool bubbletea. Includes: NavStack, Breadcrumbs, and Menu |
| 118 | 72 | [juanibiapina/gob](https://github.com/juanibiapina/gob) | Process manager for AI agents (and humans) |
| 119 | 71 | [dangerous-person/dopogoto-terminal-player](https://github.com/dangerous-person/dopogoto-terminal-player) | Terminal-native music and video player: 15 albums. 30+ hours of music. Live chat. All in your command line. |
| 120 | 70 | [Ritiksuman07/quant-whisper](https://github.com/Ritiksuman07/quant-whisper) | A hedge fund on your local machine. Terminal-native algo trading engine in Go with local LLM inference (Ollama), Bubble Tea TUI, paper/live trading, and zero cloud dependency. BYOB. |
| 121 | 69 | [satisfactorymodding/ficsit-cli](https://github.com/satisfactorymodding/ficsit-cli) | A CLI tool for managing mods for the game Satisfactory |
| 122 | 68 | [Harshil-Anuwadia/archwiki-tui](https://github.com/Harshil-Anuwadia/archwiki-tui) | A terminal browser for the Arch Wiki. No bloat. No browser. Just the wiki. |
| 123 | 66 | [parfenovvs/lazylogcat](https://github.com/parfenovvs/lazylogcat) | LazyLogcat - TUI to view Android logs from adb logcat |
| 124 | 66 | [mr-smith-org/mr](https://github.com/mr-smith-org/mr) | Framework for creating scaffolds for any existing programming language with a customizable TUI. |
| 125 | 65 | [termkit/skeleton](https://github.com/termkit/skeleton) | The Multi-tab framework of Bubbletea programs! |
| 126 | 64 | [madalinpopa/gocost](https://github.com/madalinpopa/gocost) | Simple TUI application to manage monthly expenses |
| 127 | 64 | [julez-dev/chatuino](https://github.com/julez-dev/chatuino) | A feature rich TUI Twitch IRC Client |
| 128 | 63 | [gfazioli/octoscope](https://github.com/gfazioli/octoscope) | Terminal dashboard for GitHub — followers, stars, PRs and issues at a glance, auto-refreshed |
| 129 | 62 | [contriboss/ore-light](https://github.com/contriboss/ore-light) | Lean, Bundler-compatible gem manager written in Go |
| 130 | 61 | [dhth/cueitup](https://github.com/dhth/cueitup) | Inspect messages in an AWS SQS queue in a simple and deliberate manner |
| 131 | 61 | [londek/reactea](https://github.com/londek/reactea) | Rather simple Bubble Tea companion for handling hierarchy and support for lifting state up 🚀 |
| 132 | 61 | [Cladamos/clawea](https://github.com/Cladamos/clawea) | weather forecast tui |
| 133 | 60 | [mikecsmith/ihj](https://github.com/mikecsmith/ihj) | An fzf-inspired issue tracker with vim mode and LLM-assisted backlog refinement. Jira today, more providers coming. |
| 134 | 60 | [mrusme/wth](https://github.com/mrusme/wth) | What The Heck: The better personal information dashboard for your terminal |
| 135 | 58 | [ni5arga/stock-tui](https://github.com/ni5arga/stock-tui) | Real-time stock and cryptocurrency tracker for the terminal |
| 136 | 58 | [DavidMiserak/GoCard](https://github.com/DavidMiserak/GoCard) | A lightweight file-based spaced repetition system (SRS) that uses plain Markdown files for flashcards. Perfect for developers who prefer text files, Git version control, and keyboard-driven interfaces. |
| 137 | 56 | [mujib77/cosmo](https://github.com/mujib77/cosmo) | Cosmo is a terminal user interface (TUI) that provides real-time visibility into PostgreSQL database health, query performance, lock contention, and WAL/MVCC internals. Built in Go. |
| 138 | 55 | [shafreeck/guru](https://github.com/shafreeck/guru) | A ChatGPT command line client |
| 139 | 54 | [zigzter/chatterm](https://github.com/zigzter/chatterm) | A Twitch.tv chat client in the terminal, supporting moderator actions |
| 140 | 54 | [aymanbagabas/gh-stars](https://github.com/aymanbagabas/gh-stars) | GitHub stargazers in your terminal 🌟 |
| 141 | 53 | [MohsenBg/bgscan](https://github.com/MohsenBg/bgscan) | bgscan is Ultra‑fast multi‑protocol scanner with a modular chained engine and interactive BubbleTea terminal UI. |
| 142 | 53 | [mdsakalu/zmx-session-manager](https://github.com/mdsakalu/zmx-session-manager) | TUI session manager for zmx |
| 143 | 53 | [JanSmrcka/differ](https://github.com/JanSmrcka/differ) | Terminal UI git diff viewer with two-panel layout, syntax highlighting, vim keybindings, and AI-powered commits. Built with Go + Bubble Tea. |
| 144 | 52 | [knz/catwalk](https://github.com/knz/catwalk) | test library for Bubbletea TUI models |
| 145 | 52 | [dhth/kplay](https://github.com/dhth/kplay) | Inspect messages in a Kafka topic in a simple and deliberate manner |
| 146 | 52 | [mieubrisse/teact](https://github.com/mieubrisse/teact) | A React-like component/layout framework for Charm's Bubbletea |
| 147 | 52 | [zdyxry/tokui](https://github.com/zdyxry/tokui) | An interactive TUI for visualizing code statistics from tokei. |
| 148 | 52 | [KarolosLykos/hackertea](https://github.com/KarolosLykos/hackertea) | #Hackertea is a sleek terminal user interface written in Golang, designed to bring HackerNews articles directly to your fingertips. Powered by the magic of Golang and the #bubbletea library. Just fire up your terminal, launch Hackertea, and dive into the latest HackerNews stories with ease. Happy hacking! 🚀 |
| 149 | 52 | [doeixd/nmtui-go](https://github.com/doeixd/nmtui-go) | A TUI for managing NetworkManager (nmcli) Wi-Fi connections on Linux, built with Go & Bubble Tea. |
| 150 | 51 | [zackproser/bubbletea-stages](https://github.com/zackproser/bubbletea-stages) | Simple POC for bubbletea CLI / TUI that leverages sequential "stages" of work |
| 151 | 50 | [purpleclay/dns53](https://github.com/purpleclay/dns53) | Dynamic DNS within Amazon Route53. Expose your EC2 quickly, easily and privately |
| 152 | 48 | [SteveMCWin/ttune](https://github.com/SteveMCWin/ttune) | A guitar tuner for your terminal |
| 153 | 48 | [hacel/jfsh](https://github.com/hacel/jfsh) | A terminal-based client for Jellyfin |
| 154 | 47 | [aytzey/showagent](https://github.com/aytzey/showagent) | Every AI coding session on your machine, in one TUI: browse, search, resume, branch — and convert conversations between Codex, Claude Code and more. |
| 155 | 46 | [DaltonSW/BubbleUp](https://github.com/DaltonSW/BubbleUp) | Floats your alerts to the front of your BubbleTea application |
| 156 | 46 | [brunoluiz/xpdig](https://github.com/brunoluiz/xpdig) | 🧰 Dig into Crossplane traces via TUI (a là k9s) |
| 157 | 46 | [fabio42/sasqwatch](https://github.com/fabio42/sasqwatch) | A modern take on the classic watch command |
| 158 | 44 | [EthanEFung/bubble-datepicker](https://github.com/EthanEFung/bubble-datepicker) | A datepicker component for charmbracelet/bubbletea TUI applications |
| 159 | 44 | [dhotrey/pman](https://github.com/dhotrey/pman) | A CLI project manager |
| 160 | 44 | [gateway-of-last-resort/keyward](https://github.com/gateway-of-last-resort/keyward) | Terminal UI to manage SSH keys, edit ~/.ssh/config, and audit your SSH security. Encrypted metadata + backups, an A–F security grade, byte-identical config editing. Single static binary, offline, no daemon. |
| 161 | 42 | [ChrisEdwards/abacus](https://github.com/ChrisEdwards/abacus) | Interactive terminal UI for visualizing and navigating Beads issue tracking projects |
| 162 | 42 | [Zingzy/diskbloom](https://github.com/Zingzy/diskbloom) | 🌸 a pastel treemap TUI that shows what's eating your disk |
| 163 | 42 | [olivier-w/climp](https://github.com/olivier-w/climp) | cli media player |
| 164 | 41 | [phoenix-tui/phoenix](https://github.com/phoenix-tui/phoenix) |  High-performance TUI framework for Go with DDD + Rich model inspired architecture, perfect Unicode, and Elm-inspired design. Modern alternative to Bubbletea/Lipgloss. |
| 165 | 41 | [lucky7xz/drako](https://github.com/lucky7xz/drako) | A grid-based, customizable TUI-Deck launcher |
| 166 | 40 | [elliot40404/volgo](https://github.com/elliot40404/volgo) | Volgo is a cross-platform CLI app written in Go for controlling system volume from the terminal. Use simple commands or a beautiful interactive TUI—even over SSH! |
| 167 | 40 | [Astro-Han/diffpane](https://github.com/Astro-Han/diffpane) | Real-time TUI diff viewer for AI coding agents |
| 168 | 39 | [samuelstranges/chronos](https://github.com/samuelstranges/chronos) | A Vimlike Calendar TUI |
| 169 | 39 | [ghostwright/specter](https://github.com/ghostwright/specter) | Deploy AI agents to dedicated VMs in 90 seconds. Interactive TUI. Automatic DNS and TLS. You own the infrastructure. |
| 170 | 38 | [daptify14/chezit](https://github.com/daptify14/chezit) | Terminal UI for chezmoi dotfile management |
| 171 | 38 | [dennisbergevin/mash](https://github.com/dennisbergevin/mash) | A customizable command launcher for storing and executing commands. |
| 172 | 38 | [freislot/kiri](https://github.com/freislot/kiri) | Minimalist & smart Go-based TUI assistant for plant care. Tracks watering schedules locally with SQLite, syncs with Open-Meteo API, and supports status bar integrations. |
| 173 | 38 | [programmersd21/zap](https://github.com/programmersd21/zap) | blazing fast file operations with gorgeous progress - modern replacement for cp, mv, rm |
| 174 | 37 | [seunggabi/claude-dashboard](https://github.com/seunggabi/claude-dashboard) | k9s-style TUI for managing Claude Code sessions via tmux |
| 175 | 37 | [qascade/yast](https://github.com/qascade/yast) | Yet Another Streaming Tool |
| 176 | 36 | [szktkfm/sptui](https://github.com/szktkfm/sptui) | Spotify TUI player |
| 177 | 36 | [kts982/wintui](https://github.com/kts982/wintui) | Go TUI frontend for winget (Windows Package Manager) — built with Bubble Tea, Bubbles, and Lip Gloss |
| 178 | 35 | [fullerzz/herdr-plugin-sesh](https://github.com/fullerzz/herdr-plugin-sesh) | Sesh-style workspace picker TUI for Herdr. Integrates with zoxide to create workspaces from commonly used directories. |
| 179 | 35 | [lusingander/gotip](https://github.com/lusingander/gotip) | Go Test Interactive Picker 🧪 |
| 180 | 34 | [Aploide/spettro](https://github.com/Aploide/spettro) | Spettro is a terminal coding assistant built in Go. It automates planning, coding, and testing with multi-agent workflows, model selection, and an intuitive UI. Connects to local and remote AI providers, supports conversation history, and empowers developers to streamline their workflow from the terminal. |
| 181 | 34 | [kpumuk/lazykiq](https://github.com/kpumuk/lazykiq) | rich terminal UI for Sidekiq |
| 182 | 32 | [savannahostrowski/tree-bubble](https://github.com/savannahostrowski/tree-bubble) | 🌳🫧 A TUI tree view for Charm's Bubble Tea framework |
| 183 | 32 | [overhttps/req](https://github.com/overhttps/req) | Req - Test APIs with Terminal Velocity |
| 184 | 32 | [elliot40404/easycron](https://github.com/elliot40404/easycron) | Easycron is a simple cross platform cli app that helps to configure cron jobs. |
| 185 | 31 | [DaltonSW/aocgo](https://github.com/DaltonSW/aocgo) | Go Package + CLI Tool for interacting with Advent of Code workflows |
| 186 | 30 | [ivan-penchev/system-monitor-tui](https://github.com/ivan-penchev/system-monitor-tui) | A simple System Monitor terminal user interface with Go |
| 187 | 30 | [leolorenzato/yoyo](https://github.com/leolorenzato/yoyo) | TUI command launcher 🚀 |
| 188 | 30 | [kldzj/pzmod](https://github.com/kldzj/pzmod) | Keyboard-driven terminal app + CLI for Project Zomboid dedicated-server mods: Steam Workshop search, dependency resolution, load-order, pre-boot validation, and timestamped backups. Byte-exact servertest.ini edits. |
| 189 | 30 | [lenny-ts/caddy-analyzer](https://github.com/lenny-ts/caddy-analyzer) | Fast, zero-dependency access log analyzer, security threat inspector, and TUI dashboard for Caddy v2 |
| 190 | 29 | [dunkbing/kana](https://github.com/dunkbing/kana) | Terminal app to practice typing Kana (Japanese characters) in Romaji |
| 191 | 29 | [colonyops/hive](https://github.com/colonyops/hive) | CLI/TUI for managing multiple AI agent sessions in isolated git environments with real-time status monitoring, task tracking, and inter-agent messaging. |
| 192 | 28 | [dhth/punchout](https://github.com/dhth/punchout) | punchout takes the suck out of logging time on JIRA |
| 193 | 28 | [mistakenelf/bubbletea-starter](https://github.com/mistakenelf/bubbletea-starter) | Starting point for a bubbletea app |
| 194 | 28 | [untemi/unkeyb](https://github.com/untemi/unkeyb) | A simple TUI keyboard typing speed test built using Go and the bubbletea framework. |
| 195 | 28 | [taigrr/bubbleterm](https://github.com/taigrr/bubbleterm) | Headless, embeddable ANSI terminal emulator in Go with a Bubble Tea-compatible rendering layer for building TUIs |
| 196 | 28 | [dhth/schemas](https://github.com/dhth/schemas) | Inspect postgres schemas via a TUI |
| 197 | 28 | [usepier/pier](https://github.com/usepier/pier) | Give every agent session its own VM. One command up, zero burn when idle. |
| 198 | 27 | [collinvandyck/gpterm](https://github.com/collinvandyck/gpterm) | terminal client for openai's GPT completion APIs |
| 199 | 27 | [tfcace/hash](https://github.com/tfcace/hash) | A POSIX shell with AI in the grammar: type ?? to drive Claude Code, Cursor CLI, Gemini CLI, or Ollama over ACP. Local-first, no login, no telemetry. |
| 200 | 26 | [snowtema/drift](https://github.com/snowtema/drift) | Lightweight project tracker for AI-assisted (vibe) coders. Terminal TUI + CLI + open protocol. Track dozens of projects without leaving your terminal. Single Go binary, zero dependencies. |
| 201 | 26 | [kanywst/y509](https://github.com/kanywst/y509) | A terminal user interface (TUI) tool for viewing and analyzing X.509 certificate chains |
| 202 | 25 | [mousems/cc-tmux-menu](https://github.com/mousems/cc-tmux-menu) | TUI session manager for Claude Code with tmux — built with Go + Bubbletea |
| 203 | 25 | [Abi-Liu/TextTunnel](https://github.com/Abi-Liu/TextTunnel) | A real time chat application based in the terminal |
| 204 | 25 | [scaccogatto/nastro](https://github.com/scaccogatto/nastro) | Terminal-native call recorder for macOS: system audio + mic, local files, on-device diarized transcription. No bots, no cloud. |
| 205 | 25 | [floatpane/lattice](https://github.com/floatpane/lattice) | A modular terminal dashboard built with Go and Bubble Tea |
| 206 | 24 | [winder/bubblelayout](https://github.com/winder/bubblelayout) | Declarative layout manager for BubbleTea. |
| 207 | 24 | [NimbleMarkets/ollamatea](https://github.com/NimbleMarkets/ollamatea) | BubbleTea Components for Ollama |
| 208 | 24 | [Domenez-dev/lazy-chron](https://github.com/Domenez-dev/lazy-chron) | A fast, keyboard-driven terminal UI cron job manager for Linux. Built with Go + charmbracelet/bubbletea. |
| 209 | 23 | [slinlee/bubbletea-heatmap](https://github.com/slinlee/bubbletea-heatmap) | A Charm.sh Bubbletea Heatmap component |
| 210 | 23 | [moniquelive/gocheat](https://github.com/moniquelive/gocheat) | Golang terminal client for cht.sh that uses charm.sh's bubbletea project |
| 211 | 22 | [BigJk/bubbletea-in-wasm](https://github.com/BigJk/bubbletea-in-wasm) | POC for bubbletea in WASM |
| 212 | 22 | [tmc/bubbweb](https://github.com/tmc/bubbweb) | bubbweb - bubbletea CLIs in the browser |
| 213 | 22 | [gabe565/cli-of-life](https://github.com/gabe565/cli-of-life) | Play Conway's Game of Life in your terminal. |
| 214 | 22 | [PrathamX595/weaver](https://github.com/PrathamX595/weaver) | made myself a tool to avoid the constant and repetitive boilerplate for each and every project, P.S. I published this package use it for yourself, also I appreciate any feedback    |
| 215 | 22 | [zackproser/teatutor](https://github.com/zackproser/teatutor) | Tea Tutor is a terminal UI (TUI) program that lets you take quizzes over SSH! |
| 216 | 22 | [joaoheitorgarcia/Mezzotone](https://github.com/joaoheitorgarcia/Mezzotone) | Mezzotone is a terminal UI (TUI) application written in Go that converts images and GIFs into ASCII or Unicode art. |
| 217 | 22 | [acifani/formula1-go](https://github.com/acifani/formula1-go) | 🏎 Formula1 CLI |
| 218 | 22 | [pranav718/tsuna](https://github.com/pranav718/tsuna) | peer-to-peer synchronized video watching. no servers. no accounts. just a room code and a udp packet. |
| 219 | 22 | [Sachamama/sacha](https://github.com/Sachamama/sacha) | sacha is a keyboard-first AWS TUI inspired by classic two-pane file managers. Browse, search, and manage CloudWatch Logs, S3, DynamoDB, and Lambda from your terminal. |
| 220 | 21 | [anton-abyzov/ccx-go](https://github.com/anton-abyzov/ccx-go) | Go implementation of an AI coding assistant CLI. Zero-dependency binary, goroutine-based agents, Bubbletea TUI. |
| 221 | 21 | [ntk148v/goignore](https://github.com/ntk148v/goignore) | A .gitignore wizard in your command line written in Golang |
| 222 | 20 | [martinbjeldbak/twitch-chat-cli](https://github.com/martinbjeldbak/twitch-chat-cli) | Chat on Twitch.tv from your CLI |
| 223 | 20 | [leolorenzato/aurevoir](https://github.com/leolorenzato/aurevoir) | TUI power menu 🌙 |
| 224 | 20 | [joaofnds/astro](https://github.com/joaofnds/astro) | a habit tracker for the terminal with a GitHub-style activity graph |
| 225 | 20 | [fpigeonjr/music-for-coding-tui](https://github.com/fpigeonjr/music-for-coding-tui) | A terminal UI for musicforprogramming.net — Go + Bubble Tea |
| 226 | 20 | [DoraleCitrus/gentr](https://github.com/DoraleCitrus/gentr) | The intelligent, interactive project tree generator for developers.  Navigate, filter, annotate, and export your project structure with style.And git status in a readable way. |
| 227 | 20 | [Gaurav-Gosain/scraped](https://github.com/Gaurav-Gosain/scraped) | A fast, parallelized CLI tool that scrapes web pages and converts them to markdown with an interactive TUI browser. |
| 228 | 20 | [saheersk/lazymongo](https://github.com/saheersk/lazymongo) | A fast, keyboard-driven terminal UI for MongoDB - like lazygit, but for Mongodb |
| 229 | 20 | [nnnkkk7/memtui](https://github.com/nnnkkk7/memtui) | A modern TUI client for Memcached with tree-structured key navigation, smart JSON/binary formatting, and Vim keybindings |
| 230 | 19 | [ditpoo/tictactoe-tui](https://github.com/ditpoo/tictactoe-tui) | Tui built using bubbletea for playing tic tac toe in the command line. |
| 231 | 19 | [nicolasparada/go-tea-weather](https://github.com/nicolasparada/go-tea-weather) |  |
| 232 | 19 | [jamesthesken/ipfs-tui](https://github.com/jamesthesken/ipfs-tui) | Terminal User Interface for Go-IPFS |
| 233 | 19 | [lusingander/topi](https://github.com/lusingander/topi) | Terminal OpenAPI documentation viewer 🐐 |
| 234 | 19 | [mikelorant/committed](https://github.com/mikelorant/committed) | :writing_hand: Committed is a WYSIWYG Git commit editor that helps improve the quality of your commits by showing you the layout in the same format as git log. |
| 235 | 19 | [taigrr/ssh-wars](https://github.com/taigrr/ssh-wars) | Serve the Star Wars ASCII animation (asciimation) over SSH, built with Charm Bubble Tea and Wish |
| 236 | 18 | [alexanderbh/bubbleapp](https://github.com/alexanderbh/bubbleapp) | An opinionated App Framework for BubbleTea |
| 237 | 18 | [taciturnaxolotl/charming-slack](https://github.com/taciturnaxolotl/charming-slack) | A cool slack over ssh project with a pretty bubbletea tui |
| 238 | 18 | [anchore/clio](https://github.com/anchore/clio) | An easy way to bootstrap your application with batteries included. |
| 239 | 18 | [eugeniofciuvasile/ssh-x-term](https://github.com/eugeniofciuvasile/ssh-x-term) | TUI to handle multiple SSH connections simultaneously |
| 240 | 18 | [rdalbuquerque/azdoext](https://github.com/rdalbuquerque/azdoext) | A terminal UI, powered by bubbletea framework to help streamline the process of commiting, pushing, creating PRs and following pipelines in Azure DevOps. |
| 241 | 18 | [Davidcreador/herdr-token-dashboard](https://github.com/Davidcreador/herdr-token-dashboard) | Live token spend dashboard and notifications for Herdr agent panes |
| 242 | 18 | [fawni/asunder](https://github.com/fawni/asunder) | 🍡 Command line TOTP authenticator |
| 243 | 18 | [DatCodeMania/discord-delete](https://github.com/DatCodeMania/discord-delete) | Bulk-delete your Discord messages and reactions from your data package |
| 244 | 18 | [pdfrg/rptui](https://github.com/pdfrg/rptui) | The ultimate (terminal) client for Radio Paradise |
| 245 | 18 | [bhavya-dang/pkgui](https://github.com/bhavya-dang/pkgui) | 📦 a terminal dashboard for everything you've installed |
| 246 | 18 | [Cod-e-Codes/parsec](https://github.com/Cod-e-Codes/parsec) | A fast terminal-based file inspector with live preview, fuzzy search, and multi-language support. Built for developers who need rapid code browsing. |
| 247 | 18 | [liu-ethan/golem](https://github.com/liu-ethan/golem) | Go LLM Execution Model — Go 原生 AI 编程 Agent CLI。单二进制、三层记忆、YAML 权限规则，TUI 对齐 Claude Code / Codex。 |
| 248 | 18 | [portuber/portato](https://github.com/portuber/portato) | SSH port-forwarding manager with a TUI — toggle local, remote and SOCKS5 tunnels from one screen. Single Go binary, background daemon, autostart. |
| 249 | 18 | [JimZhang168872/vpnkit](https://github.com/JimZhang168872/vpnkit) | Terminal-native manager for the mihomo proxy core — subscriptions, local nodes, proxy chains, rule-based routing. TUI + CLI, no Electron, no daemon. |
| 250 | 17 | [fulsiram/type-cli](https://github.com/fulsiram/type-cli) | terminal-based typing test built with Go and BubbleTea. |
| 251 | 17 | [metafates/bento](https://github.com/metafates/bento) | 🍱 Go framework for building TUI apps. Based on ratatui and bubbletea |
| 252 | 17 | [nkxxll/byebye](https://github.com/nkxxll/byebye) | byebye is a command line interface to logout, lock, suspend, hibernate, shutdown, reboot, from the command line on every Linux system (with xdg-screensaver on X11 or swaylock on wayland). |
| 253 | 17 | [fawni/def](https://github.com/fawni/def) | ☁️ Comfy terminal dictionary navigator |
| 254 | 17 | [kgoettler/twe](https://github.com/kgoettler/twe) | Timewarrior extension for power users, including an interactive TUI and a Go package for writing your own extensions. |
| 255 | 17 | [lusingander/ghcv-cli](https://github.com/lusingander/ghcv-cli) | view the user-created issues, pull requests, and repositories in the terminal 🧑‍💻 |
| 256 | 17 | [dasvh/go-learn-vim](https://github.com/dasvh/go-learn-vim) | An interactive terminal-based application designed to teach users the basics of Vim motions in a fun and engaging way |
| 257 | 17 | [DaltonSW/stylish](https://github.com/DaltonSW/stylish) | Simple and intuitive lscolors configuration. Put that glam in your term. |
| 258 | 17 | [newbpydev/bubblyui](https://github.com/newbpydev/bubblyui) | Vue-inspired TUI framework for Go with type-safe reactivity, component-based architecture, and 30+ built-in components |
| 259 | 16 | [b-sia/cli_kanban](https://github.com/b-sia/cli_kanban) | command line kanban built using the BubbleTea framework |
| 260 | 16 | [Genekkion/theHermit](https://github.com/Genekkion/theHermit) | A quick fix model for the Charm BubbleTea ecosystem. |
| 261 | 16 | [psychosomat/Clio](https://github.com/psychosomat/Clio) | A lightning-fast, keyboard-driven TUI for taking Markdown notes in the terminal. Powered by Go & Bubble Tea. |
| 262 | 16 | [ramirezfernando/ube](https://github.com/ramirezfernando/ube) | A fun and fast lines of code counter, made with Go! |
| 263 | 15 | [shashanktomar/sprinkles](https://github.com/shashanktomar/sprinkles) | A collection of components for bubbletea |
| 264 | 15 | [brittonhayes/glitter](https://github.com/brittonhayes/glitter) | UI components + themes for lipgloss, glamour, and bubbletea |
| 265 | 15 | [rasjonell/chessh](https://github.com/rasjonell/chessh) | SSH into a chess game |
| 266 | 15 | [shainilps/curlify](https://github.com/shainilps/curlify) | TUI API testing tool. |
| 267 | 15 | [bakayu/lq](https://github.com/bakayu/lq) | A CLI tool to add custom .gitignore and LICENSE templates to your projects right from your terminal |
| 268 | 15 | [ushiradineth/lazytf](https://github.com/ushiradineth/lazytf) | lazygit but for Terraform :D |
| 269 | 15 | [Manzanita-Research/chaparral](https://github.com/Manzanita-Research/chaparral) | 🔗 The connective tissue between your projects. Manages shared Claude Code skills and configuration across sibling repos. |
| 270 | 15 | [kbrdn1/LazyCurl](https://github.com/kbrdn1/LazyCurl) | A powerful Terminal User Interface (TUI) HTTP client |
| 271 | 15 | [sam0uly/spin](https://github.com/sam0uly/spin) | Universal project scaffolder with supreme ease, |
| 272 | 15 | [harumiWeb/eitango](https://github.com/harumiWeb/eitango) | Offline TUI English vocabulary trainer built with Go, Bubble Tea, and SQLite. |
| 273 | 14 | [harryfrzz/re-tui](https://github.com/harryfrzz/re-tui) | Re-TUI is a fast and lightweight terminal-based Reddit client written in Go. It lets you browse subreddits and read discussions directly from your command line with a clean interface. |
| 274 | 14 | [wingkwong/bootstrap-cli](https://github.com/wingkwong/bootstrap-cli) | 💻 A minimalistic CLI to bootstrap projects with different frameworks |
| 275 | 14 | [yuzuy/todo-cli](https://github.com/yuzuy/todo-cli) | A TUI todo app |
| 276 | 14 | [noborus/pgsp](https://github.com/noborus/pgsp) | PostgreSQL Stat Progress (pg_stat_progress) CLI Monitor |
| 277 | 14 | [wanyuqin/live-trading](https://github.com/wanyuqin/live-trading) | 一个终端的盯盘软件,可以添加自己需要的股票和基金代码 |
| 278 | 14 | [nore-dev/what-cli](https://github.com/nore-dev/what-cli) | what-to-code.com CLI client |
| 279 | 14 | [dennisbergevin/pwgo](https://github.com/dennisbergevin/pwgo) | Multi-list interactive cli tool to run your Playwright suite. |
| 280 | 14 | [k9withabone/fluttui](https://github.com/k9withabone/fluttui) | A terminal app for `flutter create`. |
| 281 | 14 | [lmilojevicc/seshagy](https://github.com/lmilojevicc/seshagy) | Agent-aware session manager for tmux & herdr — discover projects, launch sessions, and track your AI agents work. |
| 282 | 14 | [ParendumOU/Nexora-CLI](https://github.com/ParendumOU/Nexora-CLI) | Terminal client for Nexora — streaming chat, agents, tasks & local tool execution from your shell. Single static Go binary (Bubble Tea TUI). Connects to any Nexora/NexoraCloud instance. |
| 283 | 14 | [ggfevans/endorse](https://github.com/ggfevans/endorse) | LinkedIn CLI / TUI. LinkedIn messaging without the feed, the ads, or the "Thoughts?" posts. |
| 284 | 14 | [AyKrimino/SysAct](https://github.com/AyKrimino/SysAct) | SysAct is a terminal‑user‑interface (TUI) application written in Go, built with the Bubble Tea framework and Lip Gloss styling. It provides a quick way to perform common system actions (logout, suspend, reboot, poweroff) directly from your terminal or launcher. |
| 285 | 14 | [artyomsv/quil](https://github.com/artyomsv/quil) | Reboot-proof terminal multiplexer for AI-native devs — a tmux alternative that persists your whole workspace across reboots and auto-resumes Claude Code & OpenCode sessions. Ships an MCP server so AI agents can drive your panes. Go, cross-platform (Linux/macOS/Windows). Apache-2.0. |
| 286 | 14 | [SaneiyanReza/smsir-cli](https://github.com/SaneiyanReza/smsir-cli) | Go-based CLI & TUI for SMS.ir API — send SMS, manage credit, and configure SMS account directly from terminal. |
| 287 | 14 | [Resetnak/cooldeck](https://github.com/Resetnak/cooldeck) | Keyboard-first terminal dashboard for Coolify: deployments, live logs across your whole fleet, and an MCP server for your agent |
| 288 | 14 | [KabosuNeko/Futon](https://github.com/KabosuNeko/Futon) | A blazingly fast, minimalist TUI manga reader written in Go |
| 289 | 14 | [programmersd21/wlocks](https://github.com/programmersd21/wlocks) | 📂 see which processes hold your files open: a smooth tui alternative to lsof/fuser with auto-refresh, fuzzy search, sort modes, and themes |
| 290 | 13 | [aschey/bubbleprompt](https://github.com/aschey/bubbleprompt) | Prompt widget for bubbletea based on go-prompt |
| 291 | 13 | [mJehanno/gtop](https://github.com/mJehanno/gtop) | Gtop is a clearer alternative to gtop written in golang |
| 292 | 13 | [espcaa/slack-tui](https://github.com/espcaa/slack-tui) | a basic bubbletea tui to browse slack! (very broken, do not use) |
| 293 | 13 | [LeoCaprile/pokego](https://github.com/LeoCaprile/pokego) | A TUI that mimics a pokedex, this project was build to learn Bubbletea library and also apply golang knowledge learnt from Boot.dev |
| 294 | 13 | [Astrak00/AGDownloader](https://github.com/Astrak00/AGDownloader) | AGDownloader is a program designed to download the course contents of AulaGlobal, the implementation of moodle used at UC3M. |
| 295 | 13 | [alt-romes/llvm-c-search](https://github.com/alt-romes/llvm-c-search) | Terminal interface to search the LLVM-C API |
| 296 | 13 | [edereagzi/portui](https://github.com/edereagzi/portui) | Port-first process manager TUI for macOS, Linux, and Windows (experimental). |
| 297 | 13 | [Kumneger0/yt-tracks](https://github.com/Kumneger0/yt-tracks) | listen to youtube music in ur terminal |
| 298 | 13 | [IllusionMan1212/gorc](https://github.com/IllusionMan1212/gorc) | Modern terminal IRC client written in golang |
| 299 | 13 | [gbarany/tea-dash](https://github.com/gbarany/tea-dash) | A gh-dash-style terminal dashboard (TUI) for Gitea & Forgejo — PRs, issues, notifications, Actions runs, and branches. Go + Bubble Tea. |
| 300 | 13 | [sadopc/godu](https://github.com/sadopc/godu) | Fast, interactive disk usage analyzer for the terminal. Tree view, treemap, file type breakdown, safe deletion, and JSON export built in Go. |
| 301 | 13 | [dubeyKartikay/peacock](https://github.com/dubeyKartikay/peacock) | Beautiful JSON log viewer for your terminal. |
| 302 | 13 | [lu-zhengda/updater](https://github.com/lu-zhengda/updater) | macOS app update manager — check and update apps from Sparkle, Homebrew, Mac App Store, and GitHub Releases |
| 303 | 13 | [MegaGrindStone/doconvo](https://github.com/MegaGrindStone/doconvo) | A Terminal User Interface (TUI) application that enables interactive conversations with your documents using Large Language Models (LLM) and Retrieval-Augmented Generation (RAG) techniques. |
| 304 | 13 | [Flowtriq/nethawk](https://github.com/Flowtriq/nethawk) | Real-time network traffic analysis in your terminal. |
| 305 | 13 | [Codexia-afk/JeeraType](https://github.com/Codexia-afk/JeeraType) | JeeraType — a 100% offline, cross-platform terminal typing speed test for macOS, Windows, and Linux. Single binary, no internet required. |
| 306 | 13 | [blackwell-systems/shelfctl](https://github.com/blackwell-systems/shelfctl) | PDF/EPUB library CLI/TUI backed by GitHub Release assets (Git LFS alternative). On-demand per-file downloads, metadata catalog, migrate/split bloated repos. |
| 307 | 12 | [eznix86/irc-client](https://github.com/eznix86/irc-client) | Made with Go and BubbleTea |
| 308 | 12 | [athulreji/vault-client](https://github.com/athulreji/vault-client) | Terminal based chat application build using WebSocket and Bubble Tea Framework |
| 309 | 12 | [ryu-ryuk/porterm](https://github.com/ryu-ryuk/porterm) | interactive terminal portfolio and resume viewer |
| 310 | 12 | [kwame-Owusu/lista](https://github.com/kwame-Owusu/lista) | your todo list on the terminal |
| 311 | 12 | [MawCeron/justwrite](https://github.com/MawCeron/justwrite) | A distraction-free terminal text editor for prose — a centred page, one status line, and everything else behind a keystroke. |
| 312 | 12 | [TusharIbtekar/go-typ0](https://github.com/TusharIbtekar/go-typ0) | TUI for typing test |
| 313 | 12 | [lunguini/flatte](https://github.com/lunguini/flatte) | Build TUIs like ordinary Go — one state struct, direct mutation, pure views — the deliberate inverse of Bubble Tea |
| 314 | 12 | [vonglasow/gaia](https://github.com/vonglasow/gaia) | Gaia is a command-line interface (CLI) tool for interacting with language models via a local API. It features a beautiful terminal UI, robust configuration management, and multiple interaction modes (default, describe, code, shell) for versatile assistance with programming, system administration, and more. |
| 315 | 12 | [jasonuc/gignr](https://github.com/jasonuc/gignr) | Effortlessly Manage and Generate .gitignore files |
| 316 | 12 | [lucasmelin/gh-hook](https://github.com/lucasmelin/gh-hook) | 🪝A GitHub CLI extension to easily manage your repository webhooks. |
| 317 | 12 | [Innovology/claude-dispatcher](https://github.com/Innovology/claude-dispatcher) | Terminal cockpit for running a fleet of Claude Code sessions across all your repos — dispatch work, watch what ships, and act on what's blocked, from one screen. |
| 318 | 11 | [NotUnlikeTheWaves/minesweeper](https://github.com/NotUnlikeTheWaves/minesweeper) | A minesweeper in the terminal written in Golang with Bubbletea |
| 319 | 11 | [dhth/outtasync](https://github.com/dhth/outtasync) | Identify Cloudformation stacks that have drifted or gone out of sync |
| 320 | 11 | [camilo-zuluaga/zui](https://github.com/camilo-zuluaga/zui) | A minimalistic terminal user interface for Zotero. |
| 321 | 11 | [Fuwn/faustus](https://github.com/Fuwn/faustus) | 🛎️ A beautiful TUI for managing Claude Code sessions |
| 322 | 11 | [fzdwx/md](https://github.com/fzdwx/md) | ✍  A tui markdown editor |
| 323 | 11 | [gopad-dev/gopad](https://github.com/gopad-dev/gopad) | TUI Editor inspired by nano and powered by Tree-Sitter |
| 324 | 11 | [loog-project/loog](https://github.com/loog-project/loog) | L👀G! \\| Kubernetes Resource History Viewer |
| 325 | 11 | [Fiend3d/mc](https://github.com/Fiend3d/mc) | A TUI file manager built for gigachads |
| 326 | 11 | [Megge06/TermiCam](https://github.com/Megge06/TermiCam) | A real-time ASCII camera for your terminal. |
| 327 | 11 | [k2m30/a9s](https://github.com/k2m30/a9s) | Terminal UI AWS Resource Manager — browse, inspect, and manage 70+ AWS resource types |
| 328 | 11 | [Resetnak/herdr-logbook](https://github.com/Resetnak/herdr-logbook) | Your terminal's working memory — offline, Markdown-first notes, decisions, and an active-task now.md for Herdr |
| 329 | 11 | [adot-7/ncr-on-terminal](https://github.com/adot-7/ncr-on-terminal) | Delhi NCR in your terminal |
| 330 | 11 | [standardbeagle/mcp-tui](https://github.com/standardbeagle/mcp-tui) | Fast terminal UI and CLI for testing, debugging, and automating Model Context Protocol (MCP) servers. STDIO, SSE, HTTP, Streamable HTTP. Built on the official MCP Go SDK. |
| 331 | 11 | [itsdevcoffee/plum](https://github.com/itsdevcoffee/plum) | 🍑 Discover and manage 750+ Claude Code plugins from 12 marketplaces. Fast TUI with fuzzy search, dynamic registry, and zero setup. |
| 332 | 11 | [programmersd21/charm-masterclass](https://github.com/programmersd21/charm-masterclass) | 🎓 God-level Charm TUI masterclass — build production-grade, animated, SSH-powered terminal apps in Go with Bubble Tea, Lip Gloss, Bubbles & Wish 🚀 |
| 333 | 11 | [programmersd21/mint](https://github.com/programmersd21/mint) | 🧑‍🚀 terminal-native client for modrinth - browse, search, download, and install modpacks from your terminal |
| 334 | 11 | [PatrickKoss/keyforge.nvim](https://github.com/PatrickKoss/keyforge.nvim) | Tower defense game to master vim keybindings inside Neovim |
| 335 | 10 | [xaviergodart/bubble-carousel](https://github.com/xaviergodart/bubble-carousel) | 🎠 A carousel component for bubbletea applications |
| 336 | 10 | [shricodev/gophercast](https://github.com/shricodev/gophercast) | 🎵 Stream music to every device on your LAN, all playing at the same time. Blog: https://dev.to/shricodev/my-speaker-broke-so-i-built-a-lan-speaker-19nc |
| 337 | 10 | [its-me-abhishek/tidy](https://github.com/its-me-abhishek/tidy) | A cli tool created to reorganize directories based on their extensions |
| 338 | 10 | [indaco/prompti](https://github.com/indaco/prompti) | Interactive TUI prompts for Go CLI applications, powered by Charm. |
| 339 | 10 | [aohoyd/aku](https://github.com/aohoyd/aku) | Terminal UI for Kubernetes clusters built with Bubble Tea. Resource browsing, YAML/describe views, live log streaming, split panes, exec, debug containers, port forwarding, Helm management, and vim-style navigation. |
| 340 | 10 | [oaklandgit/vizigo](https://github.com/oaklandgit/vizigo) | A terminal spreadsheet that feels like a spreadsheet. |
| 341 | 10 | [gitcoder89431/tuitube](https://github.com/gitcoder89431/tuitube) | YouTube catalog TUI — browse curated channels, stream via mpv, download on demand. |
| 342 | 10 | [stormlightlabs/noteleaf](https://github.com/stormlightlabs/noteleaf) | read, write, & manage your life through the terminal |
| 343 | 10 | [mi-schm/tui-wiki](https://github.com/mi-schm/tui-wiki) | A minimalist, and offline-first personal wiki for the terminal. Built with Go and Bubble Tea, powered by SQLite |
| 344 | 10 | [jalpp/PassDIY](https://github.com/jalpp/PassDIY) | A TUI for password management on cloud |
| 345 | 10 | [StaiLee/Argos](https://github.com/StaiLee/Argos) | 👁️ High-performance concurrent network scanner in Go with a cinematic real-time TUI war-room dashboard — TCP connect/SYN scanning, passive fingerprinting, an Nmap-style recon command center. |
| 346 | 10 | [davidroman0O/vtable](https://github.com/davidroman0O/vtable) | Virtualized Table & List components for Bubble Tea applications in Go. Efficiently handles large datasets with minimal memory usage through virtualized rendering. Seamless integration with the Bubble Tea TUI framework. |
| 347 | 10 | [vulcanshen/kbu](https://github.com/vulcanshen/kbu) | Single-pane Kubernetes workspace — Tab / Space / Enter / Esc drives everything. Zero learning curve, Relatives navigation, YAML compare, and an embedded persistent shell built in. |
| 348 | 10 | [fuad-daoud/pkgmate](https://github.com/fuad-daoud/pkgmate) |  |
| 349 | 10 | [programmersd21/pyproject-tui](https://github.com/programmersd21/pyproject-tui) | ⌨️ keyboard-driven tui for pyproject.toml |
| 350 | 9 | [craigderington/skyterm](https://github.com/craigderington/skyterm) | A terminal-based astronomy application built with Golang and Bubbletea. |
| 351 | 9 | [javascriptizer1/grpc-cli-chat.mono](https://github.com/javascriptizer1/grpc-cli-chat.mono) | CLI chat. [Go, GRPC, Cobra, Bubbletea, Mongo, Postgres]. 2 microservices and 1 tui-client in monorepo |
| 352 | 9 | [wI2L/scrabbler](https://github.com/wI2L/scrabbler) | Pick tiles, but not yourself! |
| 353 | 9 | [donderom/bubblon](https://github.com/donderom/bubblon) | Manage nested Bubble Tea models 🥞 |
| 354 | 9 | [neyham/herdr-paddock](https://github.com/neyham/herdr-paddock) | 🐑 A card-wall feed for your herdr agents — glance over the flock, zoom into one, reply, all over plain SSH |
| 355 | 9 | [willofdaedalus/superluminal](https://github.com/willofdaedalus/superluminal) | share your terminal content to other clients across the internet |
| 356 | 9 | [taigrr/teaqlite](https://github.com/taigrr/teaqlite) | TeaQLite: a colorful, keyboard-driven TUI for browsing and editing SQLite databases, built with Bubble Tea. |
| 357 | 9 | [blacktop/fluxy](https://github.com/blacktop/fluxy) | FLUX image generator TUI |
| 358 | 9 | [overthinker1127/tui-worktree](https://github.com/overthinker1127/tui-worktree) | A terminal UI for reviewing AI-generated Git worktree changes |
| 359 | 9 | [d0d012/claude-trim](https://github.com/d0d012/claude-trim) | Audit your Claude Code token usage — see which skills eat your context budget and find conflicts before they cost you |
| 360 | 9 | [raspbeguy/ncdeck](https://github.com/raspbeguy/ncdeck) | CLI and TUI client for Nextcloud Deck |
| 361 | 9 | [KonstantinGasser/scotty](https://github.com/KonstantinGasser/scotty) | Multiplex log streams and query your logs while developing your applications |
| 362 | 9 | [DuckInAShirt/leetmate](https://github.com/DuckInAShirt/leetmate) | Terminal LeetCode coach — hints, not answers. Spaced review + ACM mode. |
| 363 | 9 | [maleesha-pramud/DevBase](https://github.com/maleesha-pramud/DevBase) | DevBase is a high-performance command-line project manager built with Go, designed for efficient discovery, organization, and management of development projects across all programming languages. Featuring optimistic UI updates, concurrent directory scanning, seamless VS Code integration, and secure GitHub cloud synchronization, it enhances producti |
| 364 | 8 | [Gaurav-Gosain/sshnake](https://github.com/Gaurav-Gosain/sshnake) | SSHnake is a classic single player Snake game built with Bubbletea served over SSH using Wish! |
| 365 | 8 | [rorycl/cexfind](https://github.com/rorycl/cexfind) | Find equipment on cex/webuy.io. Go module with multiplatform cli, webserver and bubbletea console apps. |
| 366 | 8 | [DomBlack/bubble-shell](https://github.com/DomBlack/bubble-shell) | A Bubble Tea library for creating an interactive shell using Cobra commands |
| 367 | 8 | [techygrrrl/timerrr](https://github.com/techygrrrl/timerrr) | ⏳ A CLI-based timer app with a TUI written in Go |
| 368 | 8 | [caarlos0/uhr](https://github.com/caarlos0/uhr) | Zeichenorientierte Benutzerschnittstelle Uhr |
| 369 | 8 | [thetnaingtn/kirin](https://github.com/thetnaingtn/kirin) | Scaffold a Full-stack Go gRPC Application With End-to-end Type Safety. |
| 370 | 8 | [unknwon/kargo-tui](https://github.com/unknwon/kargo-tui) | The missing Kargo console for sane operators at scale |
| 371 | 8 | [tnagatomi/gh-lsmod](https://github.com/tnagatomi/gh-lsmod) | gh-lsmod is a gh extension which allow you to browse a project go.mod's direct dependent packages |
| 372 | 8 | [sa-/wordle-tui](https://github.com/sa-/wordle-tui) | Play wordle in your terminal |
| 373 | 8 | [huseynovvusal/goch](https://github.com/huseynovvusal/goch) | 🌐 A command-line chat application for discovering and communicating with users on your local area network (LAN). |
| 374 | 8 | [pranav718/ikiru](https://github.com/pranav718/ikiru) | live system vitals for the terminal |
| 375 | 8 | [cristovaoolegario/orders-tracker-cli](https://github.com/cristovaoolegario/orders-tracker-cli) | A CLI tool to track your orders. |
| 376 | 8 | [MWhyte/rig](https://github.com/MWhyte/rig) | rig.fm |
| 377 | 8 | [akifzdemir/postgres-tui](https://github.com/akifzdemir/postgres-tui) | This is a terminal user interface (TUI) application built with Bubble Tea that connects to a PostgreSQL database and displays data in a table format. |
| 378 | 8 | [boolproof/docker-tui](https://github.com/boolproof/docker-tui) | Simple docker tui to list, start and stop your containers |
| 379 | 8 | [thetnaingtn/synrk](https://github.com/thetnaingtn/synrk) | Synchronize your forks with ease |
| 380 | 8 | [dropalltables/cdp](https://github.com/dropalltables/cdp) | coolify deploy tool |
| 381 | 8 | [isac322/rkmon](https://github.com/isac322/rkmon) | Real-time hardware monitor TUI for Rockchip RK3588 SBCs. Like htop, but for the GPU, NPU, VPU, RGA, and thermal zones of your Rock 5B+ / RK3588 board. |
| 382 | 8 | [marcoarnulfo/clickup-cli](https://github.com/marcoarnulfo/clickup-cli) | Terminal TUI for ClickUp time tracking & billing: monthly hours reports (self + team), per-list rates, billable amounts, and CSV/JSON/Markdown export. Written in Go. |
| 383 | 8 | [crosleyzack/wndr](https://github.com/crosleyzack/wndr) | A TUI Tree View for Data Files |
| 384 | 8 | [PedroMosquera/squadai](https://github.com/PedroMosquera/squadai) | One config. Every AI agent. Zero drift. Idempotent setup for OpenCode, Claude Code, Cursor, Windsurf, and Copilot. |
| 385 | 8 | [TNT-Likely/killport](https://github.com/TNT-Likely/killport) |   🔫 跨平台端口管理工具 - 一条命令解决端口占用 \\| Cross-platform CLI to kill processes on ports |
| 386 | 8 | [bastio-ai/bast](https://github.com/bastio-ai/bast) | Bast is a free, open-source CLI built to bring security to AI-powered terminal operations. It integrates with Bastio AI Security Gateway to protect your data before it reaches the LLM. |
| 387 | 8 | [eshanized/M31A](https://github.com/eshanized/M31A) | The terminal-native AI coding agent that ships, not just suggests. Six-phase workflow, git rollback chain, zero telemetry. One static binary. |
| 388 | 7 | [Yukaii/chatgpt-tui](https://github.com/Yukaii/chatgpt-tui) | A ChatGPT TUI written in Go with bubbletea |
| 389 | 7 | [tuiphi/soda](https://github.com/tuiphi/soda) | 🥤 Framework for building TUI apps based on the bubbletea |
| 390 | 7 | [akr411/doit](https://github.com/akr411/doit) | A feature-rich terminal-based todo application built with Go and Bubbletea framework. |
| 391 | 7 | [Revyy/kubeui](https://github.com/Revyy/kubeui) | A collection of gui based terminal applications designed to make working with kubernetes on the command line easier. Implemented using the awesome bubbletea(https://github.com/charmbracelet/bubbletea) module. |
| 392 | 7 | [nicolasparada/go-tea-counter](https://github.com/nicolasparada/go-tea-counter) |  |
| 393 | 7 | [frectonz/CoinGopher](https://github.com/frectonz/CoinGopher) | An expense tracker TUI app made with Bubble Tea. |
| 394 | 7 | [GuillaumeMCK/ghostygg](https://github.com/GuillaumeMCK/ghostygg) | A tool for downloading torrents without seeding & increasing download rate. |
| 395 | 7 | [haochend413/ntkpr](https://github.com/haochend413/ntkpr) | Terminal jounal management system |
| 396 | 7 | [ARJ2211/cpgrinder](https://github.com/ARJ2211/cpgrinder) | TUI based competitive programming  |
| 397 | 7 | [Qjs/Quanti-tea](https://github.com/Qjs/Quanti-tea) | Dynamic quantified-self metric prometheus exporter |
| 398 | 7 | [Vyogami/paruz](https://github.com/Vyogami/paruz) | Paruz is a fast, Terminal User Interface (TUI) for Arch Linux that acts as a visual wrapper for the paru AUR helper. |
| 399 | 7 | [yazeed1s/goback](https://github.com/yazeed1s/goback) | A terminal based command history browser |
| 400 | 7 | [sorafujitani/wez-cc-viewer](https://github.com/sorafujitani/wez-cc-viewer) | TUI dashboard for monitoring Claude Code instances across WezTerm workspaces |
| 401 | 7 | [Snappey/MQTT-Explorer](https://github.com/Snappey/MQTT-Explorer) | CLI based tool for monitoring and exploring MQTT brokers, inspired by the original MQTT Explorer |
| 402 | 7 | [harrisoncramer/joke-gpt](https://github.com/harrisoncramer/joke-gpt) | An example TUI with Cobra, Charm, and Viper that uses ChatGPT to tell you jokes |
| 403 | 7 | [notjedi/gotem](https://github.com/notjedi/gotem) | [WIP] A glamorous TUI for transmission. |
| 404 | 7 | [seanmartinsmith/beadstui](https://github.com/seanmartinsmith/beadstui) | graph-aware task management TUI for beads projects |
| 405 | 7 | [CognitoBit/netpulse](https://github.com/CognitoBit/netpulse) | Terminal network toolkit — speed test, LAN speed server, packet loss, MTR. One binary, one-line install. |
| 406 | 7 | [nart38/ollmao](https://github.com/nart38/ollmao) | Simple TUI client for ollama |
| 407 | 7 | [zrougamed/portainer-cli](https://github.com/zrougamed/portainer-cli) | A terminal user interface (TUI) for managing Portainer — browse containers, stacks, images, volumes, and environments without leaving your terminal. |
| 408 | 7 | [Simon-Hostettler/dnc](https://github.com/Simon-Hostettler/dnc) | A terminal-based TTRPG character manager |
| 409 | 7 | [pierinho13/cmdpeek](https://github.com/pierinho13/cmdpeek) | Searchable interactive command palette for discovering, previewing and running reusable terminal workflows from YAML. |
| 410 | 7 | [chalvinwz/lazyssm](https://github.com/chalvinwz/lazyssm) | Persistent-panel terminal UI for AWS SSM Session Manager — browse, filter, pin, and connect to instances with no SSH keys or open ports. |
| 411 | 7 | [Cycl0o0/OpenDeezer](https://github.com/Cycl0o0/OpenDeezer) | OpenDeezer — an open source reimplementation of Deezer. Terminal (TUI) + native cross-platform GUI; browse, stream, decrypt & play locally. |
| 412 | 7 | [furkanalp41/toktop](https://github.com/furkanalp41/toktop) | btop for your AI coding agents — a single-binary local TUI showing live token & dollar cost for Claude Code, with a budget bar that turns red before your bill does. No signup, no network. |
| 413 | 7 | [mirageglobe/scout](https://github.com/mirageglobe/scout) | a fast, read-only terminal file explorer in go with live git status, syntax-highlighted previews, and time-aware themes. inspired by ranger, broot, and yazi. |
| 414 | 7 | [Mateooo93/cortex-cli](https://github.com/Mateooo93/cortex-cli) | Sleek, fast, token-efficient AI coding agent. Multi-provider (Cortex, OpenAI, Anthropic, Ollama) with a polished terminal UI. Fork of vix. |
| 415 | 6 | [go-go-golems/bobatea](https://github.com/go-go-golems/bobatea) | Custom bubbletea bubbles |
| 416 | 6 | [anadale/teaservice](https://github.com/anadale/teaservice) | Components for BubbleTea applications |
| 417 | 6 | [raphamorim/pacman](https://github.com/raphamorim/pacman) | Pacman game for terminals built with Bubbletea |
| 418 | 6 | [rdalbuquerque/viewsearch](https://github.com/rdalbuquerque/viewsearch) | Bubbletea's viewport with vim-like search feature |
| 419 | 6 | [NStefan002/tui-calendar](https://github.com/NStefan002/tui-calendar) | - Terminal-based - Google Calendar sync - Built with BubbleTea - |
| 420 | 6 | [haykh/tuigo](https://github.com/haykh/tuigo) | a terminal UI framework written in Go using the `bubbletea` library |
| 421 | 6 | [ibrokemypie/kwatch](https://github.com/ibrokemypie/kwatch) | a tui media player for remote file servers built with bubbletea |
| 422 | 6 | [XhuyZ/lazysys](https://github.com/XhuyZ/lazysys) | A TUI application for managing systemd services, built with Go and BubbleTea. |
| 423 | 6 | [nanvenomous/pulsetui](https://github.com/nanvenomous/pulsetui) | a nice bubbletea TUI to simplify some pulseaudio actions (similar to pavucontrol) |
| 424 | 6 | [waseem-medhat/gopen](https://github.com/waseem-medhat/gopen) | Simple CLI to quick-start coding projects \\| built with Go and Bubbletea |
| 425 | 6 | [matfire/libsqltui](https://github.com/matfire/libsqltui) |  |
| 426 | 6 | [f01c33/enc-pad](https://github.com/f01c33/enc-pad) | Encryted at rest file editor |
| 427 | 6 | [sivepanda/teabag](https://github.com/sivepanda/teabag) | Go/BubbleTea Application to move your AppImage files to a centralized directory, and create an entry in your app drawer/applications directory. |
| 428 | 6 | [cpaluszek/gh-ci](https://github.com/cpaluszek/gh-ci) |  |
| 429 | 6 | [nutcas3/api-client-tui](https://github.com/nutcas3/api-client-tui) | A powerful terminal-based API client built with [Bubble Tea][bubbletea], providing a fast and efficient alternative to GUI tools like Postman or Insomnia. |
| 430 | 6 | [alswl/go-toodledo](https://github.com/alswl/go-toodledo) | Client(CLI and TUI) and SDK  for toodledo.com |
| 431 | 6 | [matt-riley/newbrew](https://github.com/matt-riley/newbrew) | Terminal UI for discovering recent Homebrew formula additions, with searchable results, cached fetches, and one-key homepage opening. |
| 432 | 6 | [ivanvc/bubble-marquee](https://github.com/ivanvc/bubble-marquee) | Marquee bubble component |
| 433 | 6 | [MawCeron/lazyftp](https://github.com/MawCeron/lazyftp) | A simple, keyboard-driven TUI FTP, FTPS and SFTP client |
| 434 | 6 | [lxwiq/AudioSort](https://github.com/lxwiq/AudioSort) | Organize your audiobook collection with intelligence and automation |
| 435 | 6 | [pdfrg/must](https://github.com/pdfrg/must) | MUSic TUI - local music player and subsonic/navidrome client, with albumart in the terminal |
| 436 | 6 | [ceuk/git-file-history](https://github.com/ceuk/git-file-history) | Browse all changes to a file |
| 437 | 6 | [bnema/gart](https://github.com/bnema/gart) | Gart is a command-line tool for managing dotfiles. |
| 438 | 6 | [AidanFogarty/pitwall](https://github.com/AidanFogarty/pitwall) | Pitwall is a F1 live timing client built for your terminal |
| 439 | 6 | [chriskim06/kubectl-topui](https://github.com/chriskim06/kubectl-topui) | a tui for kubectl top |
| 440 | 6 | [simonyos/Z-CODE](https://github.com/simonyos/Z-CODE) | A beautiful terminal-based AI coding assistant with multi-provider support (Claude, Gemini, OpenAI). |
| 441 | 6 | [Zebbeni/ansipx](https://github.com/Zebbeni/ansipx) | A golang library to convert image files to ansi art |
| 442 | 6 | [taigrr/gico](https://github.com/taigrr/gico) | A collection of tools for visualizing git commit heatmaps |
| 443 | 6 | [bapatchirag/revision](https://github.com/bapatchirag/revision) | Lazygit for SVN - a fast, keyboard-driven TUI for SVN |
| 444 | 6 | [zeozeozeo/teacrush](https://github.com/zeozeozeo/teacrush) | Compresses your videos down to a set size |
| 445 | 6 | [sroberts/plumbline](https://github.com/sroberts/plumbline) | Repo-level AI coding readiness assessment, in Go. Maps a repository to its level on the AI Codebase Maturity Model (ACMM). |
| 446 | 6 | [electroheadfx/efx-ai-skills](https://github.com/electroheadfx/efx-ai-skills) | A beautiful TUI for discovering, previewing, and managing AI agent skills across multiple providers (Claude, Cursor, Opencode, Qoder, and more) |
| 447 | 6 | [Tyooughtul/lume](https://github.com/Tyooughtul/lume) | A macOS cleanup tool that never permanently deletes anything. 55+ developer-focused scan targets,   3-stage SHA-256 duplicate detection, and a beautiful TUI. Move to Trash, never rm -rf. |
| 448 | 6 | [rabeeh-ta/manygit](https://github.com/rabeeh-ta/manygit) | A multi-repository Git CLI manager — see what's ahead, behind, or dirty across all of them at once. |
| 449 | 6 | [aorumbayev/herdr-canvas](https://github.com/aorumbayev/herdr-canvas) | Mouse-driven ASCII diagram canvas for herdr agents - draw in the TUI, share structured JSON, and let AI edit it. |
| 450 | 6 | [scoutme/milk](https://github.com/scoutme/milk) | milk — a CLI/TUI that routes prompts between agents: a primary agent and a configurable escalation agent, with session-aware state management and real-time streaming. Escalation happens automatically (via routing rules) or manually, preserving context across the handoff between agents. |
| 451 | 6 | [nashory/agent-cockpit](https://github.com/nashory/agent-cockpit) | Live terminal cockpit for token usage, cost, and speed across your coding agents (Claude Code, Codex, Gemini). Local-only, no upload. |
| 452 | 6 | [nawodyaishan/universal-mcp-sync](https://github.com/nawodyaishan/universal-mcp-sync) | Universal MCP Sync (usync) is a local-first MCP configuration manager for developers who use multiple AI clients. |
| 453 | 5 | [metafates/soda](https://github.com/metafates/soda) | 🥤 Bubbletea TUI Framework |
| 454 | 5 | [dreamjz/bubbletea-notes](https://github.com/dreamjz/bubbletea-notes) | Notes about BubbleTea TUI framwork |
| 455 | 5 | [motemen/example-go-bubbletea](https://github.com/motemen/example-go-bubbletea) |  |
| 456 | 5 | [FelineStateMachine/puzzletea](https://github.com/FelineStateMachine/puzzletea) | A collection of puzzle bubbles for BubbleTea. |
| 457 | 5 | [chriskim06/bubble-plot](https://github.com/chriskim06/bubble-plot) | a Bubbletea component for displaying graphs using braille characters |
| 458 | 5 | [kirg0/d9c](https://github.com/kirg0/d9c) | Go TUI for monitoring and managing Docker over TCP/SSH (bubbletea + Docker SDK) |
| 459 | 5 | [rajatnai49/mentat](https://github.com/rajatnai49/mentat) | Mentat is a simple tool for managing task via markdown files. |
| 460 | 5 | [NimbleMarkets/ticker_autocomplete](https://github.com/NimbleMarkets/ticker_autocomplete) | Golang library for auto-completion of financial symbols aka tickers. |
| 461 | 5 | [lusingander/edist](https://github.com/lusingander/edist) | Edit mac Stickies in terminal |
| 462 | 5 | [Gylmynnn/goclean](https://github.com/Gylmynnn/goclean) | Archlinux TUI cleaner |
| 463 | 5 | [aitmiloud/ngxtui](https://github.com/aitmiloud/ngxtui) | [experimental] - TUI app that lets you manage Nginx like a2ensite/a2dissite, but prettier. |
| 464 | 5 | [HMZElidrissi/eol-checker](https://github.com/HMZElidrissi/eol-checker) | A TUI for checking container image End-of-Life status. |
| 465 | 5 | [Chakrabortysoura/NvFile](https://github.com/Chakrabortysoura/NvFile) | A Tui Based File explorer that can work with terminal text editors to make a more cohesive experience. Written in Go with the BubbleTea framework. |
| 466 | 5 | [leslieriver/ltv-go](https://github.com/leslieriver/ltv-go) | my attempt at rewriting ltv in go |
| 467 | 5 | [Zarox28/DC-Generator](https://github.com/Zarox28/DC-Generator) | Interactive CLI to generate DevContainer configs with templates for popular languages. |
| 468 | 5 | [tristanisham/grove](https://github.com/tristanisham/grove) | Grove is a modern package manager and software installer |
| 469 | 5 | [AzraelSec/Glock](https://github.com/AzraelSec/Glock) | Simplify multi-repository project management |
| 470 | 5 | [bschimke95/jara](https://github.com/bschimke95/jara) | A terminal UI for Juju -- like k9s, but for your clouds. |
| 471 | 5 | [sorafujitani/wez-kv](https://github.com/sorafujitani/wez-kv) | Fuzzy-searchable TUI viewer for WezTerm keybindings |
| 472 | 5 | [lusingander/enigma](https://github.com/lusingander/enigma) | Terminal Enigma machine simulator ⚙️ |
| 473 | 5 | [SreeAditya-Dev/Cello-TUI](https://github.com/SreeAditya-Dev/Cello-TUI) | A beautiful terminal-based spreadsheet editor with Vim-style keybindings, AI formula generation, and rich visualization. |
| 474 | 5 | [byte2pixel/gh-statline](https://github.com/byte2pixel/gh-statline) | Statline, your GitHub team-stats TUI |
| 475 | 5 | [54L1M/snip](https://github.com/54L1M/snip) | Save long, parameterized commands as named snippets and run them with ease. |
| 476 | 5 | [AbdelilahOu/Bubly-cli](https://github.com/AbdelilahOu/Bubly-cli) | TUI to download video, audio, and transcriptions from YouTube. |
| 477 | 5 | [idan-at/lazygcs](https://github.com/idan-at/lazygcs) | A fast, keyboard-driven Terminal User Interface (TUI) for exploring and managing Google Cloud Storage (GCS). |
| 478 | 5 | [felipeospina21/mrglab](https://github.com/felipeospina21/mrglab) | Gitlab Merge Requests TUI |
| 479 | 5 | [danterolle/loqi](https://github.com/danterolle/loqi) | Loqi is a local-first, hackable, scriptable translation tool for desktop and developer workflows. Translate text, files, docs and structured content using local LLMs |
| 480 | 5 | [Cod-e-Codes/marchat-plugins](https://github.com/Cod-e-Codes/marchat-plugins) | Official plugin registry and release hub for marchat — a terminal-native, offline-first group chat application. |
| 481 | 5 | [maskedsyntax/focusbrew](https://github.com/maskedsyntax/focusbrew) | A lightweight terminal-based Pomodoro Timer built with Go and Bubble Tea |
| 482 | 5 | [stiffis/lazycider](https://github.com/stiffis/lazycider) | A modern terminal UI for Cider (Apple Music), built in Go, with playback controls, playlists, lyrics fallback, and multi-panel search. |
| 483 | 5 | [iluaii/IluskaX](https://github.com/iluaii/IluskaX) |  |
| 484 | 5 | [p-obrthr/all-in-intelligence](https://github.com/p-obrthr/all-in-intelligence) | texas holdem poker per terminal |
| 485 | 5 | [StefanoGuerrini/c9s](https://github.com/StefanoGuerrini/c9s) | Terminal dashboard for Claude Code, zero config, track all sessions, switch context instantly, never lose track of what's running |
| 486 | 5 | [umairabid/lazysql](https://github.com/umairabid/lazysql) | A minimal, Vim-friendly TUI for PostgreSQL — clean three-pane UI, full Vim query editor, and a shell-command credential mode for rotating secrets. |
| 487 | 5 | [0xDarkMatter/conclave](https://github.com/0xDarkMatter/conclave) | Stop juggling six AI CLIs. Conclave is your universal remote for LLMs - query any model with one syntax, or unleash them all in parallel and let a judge synthesize the chaos into a verdict. Why trust one robot when you can have a whole council argue about your code? |
| 488 | 5 | [studiowebux/restcli](https://github.com/studiowebux/restcli) | REST CLI + TUI - A way to test your API |
| 489 | 5 | [jabbott-iii/Munus](https://github.com/jabbott-iii/Munus) | Munus is a lightweight Go-based task management TUI/CLI tool for organizing, tracking, and managing tasks from the terminal. Built in 100% Go, it offers a fast, efficient command-line workflow for productivity-focused users who want a simple, terminal-native task manager with a clean text user interface and developer-friendly CLI experience. |
| 490 | 4 | [andrewsomething/bubbletea-droplet](https://github.com/andrewsomething/bubbletea-droplet) |  |
| 491 | 4 | [xpufx/bubbletea-layout-examples](https://github.com/xpufx/bubbletea-layout-examples) |  |
| 492 | 4 | [zackproser/unicode-cli](https://github.com/zackproser/unicode-cli) | small bubbletea and glamour example program |
| 493 | 4 | [yardbirdsax/bubblewrap](https://github.com/yardbirdsax/bubblewrap) | A wrapper around bubbles and bubbletea from Charmbracelet |
| 494 | 4 | [mieubrisse/vim-bubble](https://github.com/mieubrisse/vim-bubble) | A BubbleTea component for emulating a Vim buffer |
| 495 | 4 | [Gaurav-Gosain/game-of-life](https://github.com/Gaurav-Gosain/game-of-life) | Conway's Game of Life, written in Go + Bubbletea |
| 496 | 4 | [melisapo/pomogo](https://github.com/melisapo/pomogo) | Terminal Pomodoro timer built with Go and Bubbletea :3 |
| 497 | 4 | [mahlburgc/teaterm](https://github.com/mahlburgc/teaterm) | Serial Terminal TUI written in GO using the bubbletea framework   |
| 498 | 4 | [rorycl/lsbookmarks](https://github.com/rorycl/lsbookmarks) | List and interactively search firefox bookmark jsonlz4 files from the terminal using bubbletea |
| 499 | 4 | [Bhargav16exd/nginxctl](https://github.com/Bhargav16exd/nginxctl) | Nginxctl \\| Utility built to interact with Nginx over terminal UI and make config easy. |
| 500 | 4 | [anotherhadi/usbguard-tui](https://github.com/anotherhadi/usbguard-tui) | A  terminal UI for managing USB devices via usbguard, with keybindings & mouse support. TUI built with golang & bubbletea. |
| 501 | 4 | [wangYX657211334/nacos-tui](https://github.com/wangYX657211334/nacos-tui) | 一个nacos可视化终端工具 |
| 502 | 4 | [15jgme/tusk](https://github.com/15jgme/tusk) | A cli tool for updating docker containers while keeping the same port settings |
| 503 | 4 | [dhth/ecsv](https://github.com/dhth/ecsv) | Quickly check the versions of your systems running in ECS tasks across various environments |
| 504 | 4 | [Huseynteymurzade28/bytesizepet](https://github.com/Huseynteymurzade28/bytesizepet) | A byte-sized virtual pet living entirely in your TTY. An asynchronous, state-driven Tamagotchi built with Go and Bubbletea. |
| 505 | 4 | [gunererd/helix-health](https://github.com/gunererd/helix-health) | Overengineered helix --health |
| 506 | 4 | [achill06/git-zen](https://github.com/achill06/git-zen) | A Go CLI extension and interactive TUI for git branch management with fuzzy search and live GitHub PR statuses. |
| 507 | 4 | [L-Michael1/clinote](https://github.com/L-Michael1/clinote) | TUI notes manager to read, edit, and add notes in style |
| 508 | 4 | [aymanhs/jeeves](https://github.com/aymanhs/jeeves) | TUI for working with Systemd  |
| 509 | 4 | [tmustier/economist-tui](https://github.com/tmustier/economist-tui) | Read The Economist from your terminal using your subscription. Unofficial. Comes with CLI and SKILL.md. |
| 510 | 4 | [Jarimus/BibleTUI](https://github.com/Jarimus/BibleTUI) | A text-based user interface for reading the Bible with the scripture.api.bible |
| 511 | 4 | [brennerm/slashmetrics-cli](https://github.com/brennerm/slashmetrics-cli) | A terminal application to visually explore Prometheus metrics endpoints |
| 512 | 4 | [BL19/commands-wiki-cli](https://github.com/BL19/commands-wiki-cli) | A cli for https://commands.wiki |
| 513 | 4 | [fchastanet/shell-command-bookmarker](https://github.com/fchastanet/shell-command-bookmarker) | A terminal-based UI for managing shell commands with bookmarking, categorization and search capabilities. Built with Bubbletea framework, featuring tab navigation, keyboard shortcuts and focus management. |
| 514 | 4 | [simonwhitaker/shellsnake](https://github.com/simonwhitaker/shellsnake) | Play Snake on the command line! 😄🐛🐛🐛 |
| 515 | 4 | [krisvandebroek/opencode-launcher](https://github.com/krisvandebroek/opencode-launcher) | Lightning-fast TUI launcher for OpenCode: pick a project + model, or resume a session. |
| 516 | 4 | [jwc20/svt](https://github.com/jwc20/svt) | SSH TUI game for the entrepreneurially-inclined |
| 517 | 4 | [provsalt/soramail](https://github.com/provsalt/soramail) | Generate forwarded email addresses using Cloudflare email routing in your terminal! |
| 518 | 4 | [bilguun0203/tailscale-tui](https://github.com/bilguun0203/tailscale-tui) | TUI app for viewing your Tailscale devices. |
| 519 | 4 | [nathanielfernandes/wrdl](https://github.com/nathanielfernandes/wrdl) | a wordle solver in your cli |
| 520 | 4 | [arush-sal/bulk-delete-chatgpt-conversations](https://github.com/arush-sal/bulk-delete-chatgpt-conversations) | A Go terminal UI that authenticates with your ChatGPT browser session, load all conversations, let you select multiple entries, and then bulk archive or bulk delete them. |
| 521 | 4 | [nayanjohnston/shanty](https://github.com/nayanjohnston/shanty) | A terminal music player for Navidrome! |
| 522 | 4 | [Denklinie/peer-pressure](https://github.com/Denklinie/peer-pressure) |  |
| 523 | 4 | [shawalli/bubbles](https://github.com/shawalli/bubbles) | TUI components for Bubble Tea 🫧 |
| 524 | 4 | [audreyteles/githp](https://github.com/audreyteles/githp) | A TUI app to help you with git commit process. |
| 525 | 4 | [Bloby22/andtls](https://github.com/Bloby22/andtls) | Interactive terminal dashboard for monitoring and managing Android devices via ADB. |
| 526 | 4 | [divijg19/Grimoire](https://github.com/divijg19/Grimoire) | RPG game created around a data persistence engine used for key-value storage with TUI, and CLI REPL implementation |
| 527 | 4 | [AxBolduc/gomlb](https://github.com/AxBolduc/gomlb) | A TUI application for browsing MLB baseball games and statistics  |
| 528 | 4 | [satya-sudo/editgo](https://github.com/satya-sudo/editgo) | A terminal-based text editor written in Go, designed to explore data structures like Stack and Trie, and utilize Go concurrency patterns (channels, goroutines). The editor provides essential editing features and operates within the terminal using the Bubbletea framework. |
| 529 | 4 | [theognis1002/lightfold-cli](https://github.com/theognis1002/lightfold-cli) | Lightfold CLl - Minimal deployment tool for indie devs. Detects your app framework, builds, and deploys to your own VPS with simple defaults. No YAML, no lock-in. |
| 530 | 4 | [alcb1310/kanban](https://github.com/alcb1310/kanban) | CLI tool to manage Kanban boards |
| 531 | 4 | [Polqt/gitflowtui](https://github.com/Polqt/gitflowtui) | gitflow-tui is a keyboard-first CLI app for developers who use GitFlow-style branching and want fewer context switches than raw Git commands or GUI tools. |
| 532 | 4 | [muzzlol/nomodit](https://github.com/muzzlol/nomodit) | A tui/cli tool for interfacing with a LLM fine-tuned on various language tasks. It emphasizes on making the user see the changes made in order to learn |
| 533 | 4 | [garyblankenship/gist-blog](https://github.com/garyblankenship/gist-blog) | Transform GitHub Gists into a beautiful blog with $0 hosting on Cloudflare Workers. Features H2-only content enhancement system for dual-purpose articles (30-second scan or deep read). |
| 534 | 4 | [paulo-amaral/dockup](https://github.com/paulo-amaral/dockup) | dockup - interactive TUI to install, harden and maintain Docker Engine, Compose v2, NVIDIA Container Toolkit and Apple container. One command, checksum-verified. |
| 535 | 4 | [hiroaqii/bgg-tui](https://github.com/hiroaqii/bgg-tui) | A terminal user interface for browsing BoardGameGeek |
| 536 | 4 | [mubbie/stacksmith](https://github.com/mubbie/stacksmith) | Ultralight Artisan Git Stacking Tool |
| 537 | 4 | [ariguillegp/rivet](https://github.com/ariguillegp/rivet) | workspace manager for humans to steer agents |
| 538 | 4 | [mogglemoss/lazytailscale](https://github.com/mogglemoss/lazytailscale) | A terminal dashboard for your Tailscale network. Two-pane keyboard-driven TUI: peer list on the left, selected-peer detail on the right. Runs entirely from your local Tailscale socket — no API key, no cloud, no opinions about your network topology. |
| 539 | 4 | [qyinm/phtui](https://github.com/qyinm/phtui) |  |
| 540 | 4 | [tirthpatell/mdr](https://github.com/tirthpatell/mdr) | Markdown renderer, editor, and linter for the terminal |
| 541 | 4 | [leonardorifeli/clawtop](https://github.com/leonardorifeli/clawtop) | Multi-host TUI dashboard for Anthropic Claude subscription usage. Daemon + viewer split keeps OAuth credentials on your workstation while a tmux dashboard on any host renders the merged view. |
| 542 | 4 | [Gaurav-Gosain/streamd](https://github.com/Gaurav-Gosain/streamd) | A CLI tool that renders streamed LLM output as beautiful markdown in the terminal |
| 543 | 4 | [ahmadraza100/dotlock](https://github.com/ahmadraza100/dotlock) | Encrypted .env vault manager with interactive TUI — written in Go |
| 544 | 4 | [gadflysu/aps](https://github.com/gadflysu/aps) | Interactive session picker for coding agents (Claude Code, OpenCode) |
| 545 | 4 | [luthermonson/linode-tui](https://github.com/luthermonson/linode-tui) | k9s-inspired terminal UI for the Linode API. Browse and manage Linodes, LKE, NodeBalancers, Volumes, and more with a : command palette, multi-pane layouts, and drill-ins to k9s, lazysql, SSH, and lish. |
| 546 | 4 | [Cod-e-Codes/voicelog](https://github.com/Cod-e-Codes/voicelog) | Terminal-based voice memo app in Go with TUI, audio recording/playback, and memo management. |
| 547 | 4 | [aaryanrwt/chibi](https://github.com/aaryanrwt/chibi) | The AI-native, predictive SRE assistant for Kubernetes. Chibi correlates live Prometheus metrics with cluster state to provide instant, terminal-based LLM diagnostics. |
| 548 | 4 | [psychedelicdevx/bosun](https://github.com/psychedelicdevx/bosun) | A fast terminal UI for Docker and Podman. Containers, images, volumes, networks, compose stacks, live logs and stats, local or over SSH. |
| 549 | 4 | [Cod-e-Codes/tuitar](https://github.com/Cod-e-Codes/tuitar) | Terminal-based guitar tablature editor with modal Vim-style editing, real-time visual feedback, and MIDI playback — built in Go with Bubble Tea. |
| 550 | 4 | [lumipallolabs/diskdive](https://github.com/lumipallolabs/diskdive) |  A fast, terminal-based disk space analyzer for macOS, Windows and Linux |
| 551 | 4 | [Phundahl/tailtui](https://github.com/Phundahl/tailtui) | The terminal-based control room for Tailscale. A Go-built TUI to monitor nodes, validate subnet routes, toggle advanced settings, and manage profiles securely. |
| 552 | 4 | [Cod-e-Codes/gophetch](https://github.com/Cod-e-Codes/gophetch) | Terminal system monitor with animated ASCII rain clouds, interactive multi-tab interface, real-time system metrics, customizable frames, and cross-platform support built in Go using Bubble Tea. |
| 553 | 4 | [stxkxs/mkt](https://github.com/stxkxs/mkt) | Real-time stock & crypto market dashboard for the terminal — live prices, charts, portfolio P&L, and alerts. No API keys. |
| 554 | 4 | [gowtham012/pair-live](https://github.com/gowtham012/pair-live) | Real-time multiplayer coding from the terminal. Any editor. Any AI. Built-in chat, file sync, AI context sharing. |
| 555 | 4 | [dsk1ra/snipe-cli](https://github.com/dsk1ra/snipe-cli) | Terminal job-search pipeline on local Ollama. Self-hosted models score a posting against your CV, then tailor a one-page PDF using only bullets your CV already contains. It fabricates nothing and uploads nothing. |
| 556 | 4 | [halpworld/halpradio](https://github.com/halpworld/halpradio) | 📻 LazyVim-inspired Terminal Internet Radio Streamer built with Go & Bubble Tea. 30,000+ stations, beat-reactive animal DJ visualizers, live ICY metadata, and zero-dependency audio. |
| 557 | 4 | [airiclenz/apogee](https://github.com/airiclenz/apogee) | Terminal coding agent for local LLMs (llama.cpp, Ollama, vLLM) and any OpenAI-compatible API. OS-sandboxed autonomy, MCP, sessions. Go. |
| 558 | 3 | [ankkyprasad/compass](https://github.com/ankkyprasad/compass) | Navigator for bubbletea |
| 559 | 3 | [anhoder/foxful-cli](https://github.com/anhoder/foxful-cli) |  |
| 560 | 3 | [ignoxx/bubbles](https://github.com/ignoxx/bubbles) | bubbletea compatible component collection |
| 561 | 3 | [vknabel/go-bubblenav](https://github.com/vknabel/go-bubblenav) | Navigation helpers for charm bubbletea |
| 562 | 3 | [willgorman/teash](https://github.com/willgorman/teash) | Put some bubbletea in your tsh |
| 563 | 3 | [srihari93/bubble-list](https://github.com/srihari93/bubble-list) | A fully featured scrollable bubbletea list component. |
| 564 | 3 | [xevrion/Chronapse](https://github.com/xevrion/Chronapse) | A minimalist Linux timelapse recorder with a Bubbletea TUI and Python backend — capture time, frame by frame, the smart way. |
| 565 | 3 | [srrathi/go-basic-tui](https://github.com/srrathi/go-basic-tui) | Basic Weather TUI made using bubbletea in Golang |
| 566 | 3 | [ary82/pacman](https://github.com/ary82/pacman) | single-level tui pacman game made with bubbletea |
| 567 | 3 | [thisguycodes/modelstack](https://github.com/thisguycodes/modelstack) | Cleanly nested Modes for https://github.com/charmbracelet/bubbletea |
| 568 | 3 | [nathanaelcunningham/tmuxSessions](https://github.com/nathanaelcunningham/tmuxSessions) | tmux session manager written in go, uses charm.sh bubbletea |
| 569 | 3 | [Polqt/clifolio](https://github.com/Polqt/clifolio) | Interactive terminal portfolio accessible via SSH. Built with Go and Bubbletea. |
| 570 | 3 | [williavs/charm-dev-skill-marketplace](https://github.com/williavs/charm-dev-skill-marketplace) | Skills for building terminal UIs with Bubbletea and the Charm ecosystem |
| 571 | 3 | [kevmul/clockify-tui](https://github.com/kevmul/clockify-tui) | Version 2 of a Clockify TUI created in Go with Bubbletea |
| 572 | 3 | [Morphclue/ygo-bubble-tea](https://github.com/Morphclue/ygo-bubble-tea) | CLI for Yu-Gi-Oh! cards made with Bubble Tea |
| 573 | 3 | [halsten-dev/orvyn](https://github.com/halsten-dev/orvyn) | Orvyn is built on top of BubbleTea, helping developping complexe TUI applications. |
| 574 | 3 | [jeremyengland/peekctx](https://github.com/jeremyengland/peekctx) | A TUI for scraping and inspecting Claude Code's context built with Bubbletea + Lipgloss |
| 575 | 3 | [espcaa/spaceship-tui](https://github.com/espcaa/spaceship-tui) | a simple go bubbletea tui to see/edit dns records on spaceship.com domains |
| 576 | 3 | [pgavlin/tea-grid](https://github.com/pgavlin/tea-grid) | An AG Grid-inspired data grid component for [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| 577 | 3 | [nick-popovic/custom-bubbles](https://github.com/nick-popovic/custom-bubbles) | Custom Bubbles I wrote that I found to be reuseable and helpful :) |
| 578 | 3 | [aashishpanchal/chatgpt_go](https://github.com/aashishpanchal/chatgpt_go) | chatgpt_go is a powerful terminal-based client for chatgpt, built in Go. |
| 579 | 3 | [HalxDocs/lazydb](https://github.com/HalxDocs/lazydb) | A keyboard-driven terminal UI for exploring Postgres, MySQL, and SQLite databases — built with Go, bubbletea, and lipgloss. |
| 580 | 3 | [leonardobiffi/wave](https://github.com/leonardobiffi/wave) | Terminal radio player written in Golang |
| 581 | 3 | [aziis98/menu](https://github.com/aziis98/menu) | A small dmenu like tool for the terminal with integrated fuzzy search functionality, easily extensible with shell scripts |
| 582 | 3 | [irskep/localci](https://github.com/irskep/localci) | CI that runs locally, integrated with Mise, with web UI and TUI |
| 583 | 3 | [Josehpequeno/lumus](https://github.com/Josehpequeno/lumus) | Lumus is a command-line tool written in Go that allows you to read PDF files  |
| 584 | 3 | [yassernasc/td](https://github.com/yassernasc/td) | todo app for the command-line |
| 585 | 3 | [moaqz/news](https://github.com/moaqz/news) | 📺 CLI to watch developer news |
| 586 | 3 | [n-g-u-y-e-n/docker-status](https://github.com/n-g-u-y-e-n/docker-status) | A K9s-inspired terminal user interface (TUI) for Docker management. |
| 587 | 3 | [connerohnesorge/spectr](https://github.com/connerohnesorge/spectr) | validatable spec driven development (inspired by openspec and kiro) |
| 588 | 3 | [Kumneger0/cligram](https://github.com/Kumneger0/cligram) | CLI-based Telegram client |
| 589 | 3 | [divijg19/Peony](https://github.com/divijg19/Peony) | CLI-first Cognitive holding space for unfinished thoughts |
| 590 | 3 | [admcpr/gh-reponark](https://github.com/admcpr/gh-reponark) | Explore and audit settings across your GitHub repositories with a delightful tui 🔬 |
| 591 | 3 | [marcelblijleven/bubbles-hlist](https://github.com/marcelblijleven/bubbles-hlist) | Horizontal version of the bubbles list component |
| 592 | 3 | [IwnuplyNotTyan/Hera](https://github.com/IwnuplyNotTyan/Hera) | 🐙 ~ Turn Based RogueLike? |
| 593 | 3 | [angel-git/doings](https://github.com/angel-git/doings) | TUI application for managing tasks in kanban-ish way using Markdown files |
| 594 | 3 | [techquestsdev/crontab-guru](https://github.com/techquestsdev/crontab-guru) | Interactive terminal-based cron expression editor built with Go and Bubble Tea |
| 595 | 3 | [hopefulTex/rainbownya](https://github.com/hopefulTex/rainbownya) | a lolcat clone |
| 596 | 3 | [ismaelosuna7824/herdr-file-viewer](https://github.com/ismaelosuna7824/herdr-file-viewer) | A keyboard-driven file explorer, code viewer and git client in a single Herdr pane — Go + Bubble Tea. |
| 597 | 3 | [ernesto27/dcli](https://github.com/ernesto27/dcli) | A  docker tui  managment  |
| 598 | 3 | [Jaybee18/todo](https://github.com/Jaybee18/todo) | a minimalistic todo app for the terminal |
| 599 | 3 | [itisbryan/herdr-gh-checks](https://github.com/itisbryan/herdr-gh-checks) | Herdr plugin: watch & review the current PR's CI in a pane; CI/merge status on sidebar rows. Go + Bubble Tea. |
| 600 | 3 | [khuynh22/racellm](https://github.com/khuynh22/racellm) | Race multiple LLMs simultaneously from your terminal — fan-out one prompt to OpenAI, Anthropic, Gemini & Ollama, stream results in parallel, and see who wins first. Live BubbleTea TUI with timing stats. |
| 601 | 3 | [m4rii0/rdtui](https://github.com/m4rii0/rdtui) | Real Debrid terminal UI client |
| 602 | 3 | [WilsonNet/mvave-chocolate-tui](https://github.com/WilsonNet/mvave-chocolate-tui) | Terminal UI for configuring the M-Vave Chocolate MIDI footswitch controller on Linux |
| 603 | 3 | [alexaldearroyo/catselector](https://github.com/alexaldearroyo/catselector) | Interactive file selector for concatenating and exporting text files |
| 604 | 3 | [Prathamesh0901/journal-tui](https://github.com/Prathamesh0901/journal-tui) | A terminal UI for tailing and filtering journald/systemd logs in real time, built with Go, Cobra, and Bubble Tea. |
| 605 | 3 | [kurojs/ovpngate](https://github.com/kurojs/ovpngate) | Terminal-based OpenVPN client for VPN Gate with server list, filtering, and one-key connect |
| 606 | 3 | [BvChung/prtls](https://github.com/BvChung/prtls) | Terminal UI For File Directory Traversal and Visual Representation |
| 607 | 3 | [marcantoineg/ls-projects](https://github.com/marcantoineg/ls-projects) | A simple shell app to list projects and open them in a new window of VS Code. 📚 |
| 608 | 3 | [TheLustriVA/fontlet](https://github.com/TheLustriVA/fontlet) | A TUI for figlet written in Go |
| 609 | 3 | [blacktop/lifx](https://github.com/blacktop/lifx) | LIFX Light TUI and MCP Server |
| 610 | 3 | [cloudboy-jh/bentotui](https://github.com/cloudboy-jh/bentotui) | The App framework for Bubble Tea 🍱 |
| 611 | 3 | [divijg19/Trellis](https://github.com/divijg19/Trellis) | A concurrent task lifecycle manager & runtime for executing, scheduling, and observing background jobs with a webUI task-tracker in GoTH, Dart, Jaspr |
| 612 | 3 | [ohSystemmm/Ani-Track](https://github.com/ohSystemmm/Ani-Track) | Your simple tool for tracking anime and manga progress. |
| 613 | 3 | [jakmaz/arcade](https://github.com/jakmaz/arcade) | Classic games for your terminal |
| 614 | 3 | [yanmoyy/go-go-go](https://github.com/yanmoyy/go-go-go) | A Terminal-based Multiplayer Stone Shooting Game written in Go |
| 615 | 3 | [camilin7483/cine-cli](https://github.com/camilin7483/cine-cli) | cine-cli — Watch movies and TV shows from your terminal. A modern CLI tool with TUI, multiple providers, downloads, fuzzy search, i18n support (6 languages), TMDB integration, and auto-update. Built with Go + Bubble Tea. |
| 616 | 3 | [divijg19/Nightshade](https://github.com/divijg19/Nightshade) | Live terminal-native multiplayer fog of war RPG with humans and machine learning agents |
| 617 | 3 | [pol-cova/termkit-go](https://github.com/pol-cova/termkit-go) | Composable charts, motion, and polished terminal UI components for Go CLIs and TUIs |
| 618 | 3 | [rootlyhq/rootly-tui](https://github.com/rootlyhq/rootly-tui) | Terminal UI for viewing Rootly incidents and alerts |
| 619 | 3 | [jstreitb/baa](https://github.com/jstreitb/baa) | BAA: The cozy TUI that herds your package managers. Unified updates for apt, pacman, snap & flatpak... 🐑 |
| 620 | 3 | [truffle-dev/glyph](https://github.com/truffle-dev/glyph) | Beautifully designed components for the terminal. Yours to copy, paste, own. |
| 621 | 3 | [ollykeran/sshush](https://github.com/ollykeran/sshush) | Go SSH Agent and utilities  |
| 622 | 3 | [nicolito128/tempo](https://github.com/nicolito128/tempo) | TUI music player written in Go. |
| 623 | 3 | [blackcoderx/falcon](https://github.com/blackcoderx/falcon) | A terminal-based AI agent for API developers. |
| 624 | 3 | [thsnkhn/bluff](https://github.com/thsnkhn/bluff) | A terminal-first poker bank and game ledger. |
| 625 | 3 | [hluaguo/commity](https://github.com/hluaguo/commity) | AI-powered git commit message generator with TUI. Select files, generate conventional commits, split changes, and more. |
| 626 | 3 | [mrkayhyun/gito](https://github.com/mrkayhyun/gito) | A fast, subcommand-based Git TUI for people who live in the terminal. Go + Bubble Tea. |
| 627 | 3 | [lestex/torrnado](https://github.com/lestex/torrnado) | A terminal BitTorrent client with a vim-like TUI and a daemon that keeps running without it. |
| 628 | 3 | [georgebuilds/carlos](https://github.com/georgebuilds/carlos) | BYOK coding and research agent. Obsidian-compatible markdown wiki as a baked-in memory system. Isolate work from life using frames. |
| 629 | 3 | [subhadeeproy3902/pong-ball](https://github.com/subhadeeproy3902/pong-ball) | A minimalist, physics-based paddleball game for the terminal — sub-stepped collisions, a spring-driven paddle (keys or mouse), five restrained themes, four modes, and persistent score history. One binary, any OS. |
| 630 | 3 | [luizfx22/essh](https://github.com/luizfx22/essh) | A fast CLI and TUI tool to visually manage your SSH connections and auto-fix private key permission errors on the fly. |
| 631 | 3 | [getkonfi/konfi](https://github.com/getkonfi/konfi) | TUI for exploring your favourite tools' configurations 🔍 you can also edit ✨ |
| 632 | 3 | [armstrongl/nd](https://github.com/armstrongl/nd) | Napoleon dynamite - a coding agent asset management CLI tool |
| 633 | 3 | [flyingnobita/llml](https://github.com/flyingnobita/llml) | Terminal UI for discovering local GGUF and safetensors models, detecting llama.cpp or vLLM runtimes, and launching them with saved presets. |
| 634 | 3 | [icortesb/cupstui](https://github.com/icortesb/cupstui) | A terminal interface for CUPS: queue, printers, printing, quotas, history and logs |
| 635 | 3 | [janekbaraniewski/synoctl](https://github.com/janekbaraniewski/synoctl) | A management terminal tool for Synology DSM. Browse and manage files, backups, packages, containers, VMs and many many more. |
| 636 | 3 | [JoshuaAFerguson/terminal-velocity](https://github.com/JoshuaAFerguson/terminal-velocity) | A multiplayer space trading and combat game inspired by Escape Velocity, playable entirely through SSH |
| 637 | 3 | [brunobrise/tf-drift](https://github.com/brunobrise/tf-drift) | Terraform/OpenTofu drift scanner for layered workspaces. |
| 638 | 3 | [benstroud/lazygaze](https://github.com/benstroud/lazygaze) | Split-pane TUI for AI code review. Pipes git diffs to Claude CLI or GitHub Copilot CLI with streaming output, prompt library, and persona system. |
| 639 | 3 | [crossben/orchestra-code](https://github.com/crossben/orchestra-code) | The operating system for AI coding agents — one supervised interface to run Claude Code, OpenCode, Mimo & more. |
| 640 | 3 | [qrxnz/gopuntes](https://github.com/qrxnz/gopuntes) | A flexible note-browsing tool  💖🇪🇸 |
| 641 | 3 | [prabalesh/croptop](https://github.com/prabalesh/croptop) | A beautiful, real-time terminal-based system monitor for Linux, macOS, and Windows. Built with Go and Bubble Tea, CropTop provides an interactive interface to monitor CPU, memory, processes, network, disk usage, and battery status with smooth progress bars and keyboard navigation. |
| 642 | 3 | [fezcode/atlas.websearch](https://github.com/fezcode/atlas.websearch) | A blazing fast, interactive CLI search tool for the terminal. Search DuckDuckGo, Wikipedia, Hacker News, and Reddit with a beautiful TUI. Zero configuration required. |
| 643 | 3 | [index-null/cmus-lyric](https://github.com/index-null/cmus-lyric) | Modern synced lyrics viewer for cmus with auto-fetching from LRCLIB/Netease, supports cover. Built with Bubble Tea TUI framework. One-line install via Homebrew. |
| 644 | 3 | [subhashraveendran/aero-shutter](https://github.com/subhashraveendran/aero-shutter) | Fast, cable-free photo transfer from Wi-Fi Nikon cameras over their native PTP/IP protocol — a click-first Go terminal app and an iOS & Android companion. |
| 645 | 3 | [KamilSupera/github-pullrequests-checkecker](https://github.com/KamilSupera/github-pullrequests-checkecker) | Terminal UI for reviewing GitHub pull requests with Claude Code or Cursor — drafts pending reviews you submit yourself |
| 646 | 3 | [Cod-e-Codes/dungeon-dash](https://github.com/Cod-e-Codes/dungeon-dash) | A terminal-based dungeon crawler game built with Go and Bubble Tea. Navigate levels, collect treasures, avoid traps and enemies. |
| 647 | 3 | [XenomorphingTV/burrow](https://github.com/XenomorphingTV/burrow) | Task catalogue, runner and scheduler |
| 648 | 3 | [pageton/dbview](https://github.com/pageton/dbview) | Terminal TUI database viewer for SQLite, MySQL, MariaDB, PostgreSQL, CockroachDB, MSSQL, MongoDB, Redis, and Cassandra |
| 649 | 3 | [CaseyBlackburn/glsms](https://github.com/CaseyBlackburn/glsms) | Go library, CLI, and TUI for reading and sending SMS through a GL.iNet cellular router's JSON-RPC API (firmware 4.x, e.g. GL-X3000 / Spitz AX). |
| 650 | 3 | [Yyyangshenghao/goani-cli](https://github.com/Yyyangshenghao/goani-cli) | Go 语言编写的命令行动漫播放器，支持多中文动漫源搜索与播放。 |
| 651 | 3 | [floatpane/bubble-overlay](https://github.com/floatpane/bubble-overlay) | ANSI-aware overlay painter for Bubble Tea / lipgloss views |
| 652 | 3 | [ibnaleem/vtscan](https://github.com/ibnaleem/vtscan) | 🛡️ VirusTotal for the terminal |
| 653 | 3 | [sonquer/opendba](https://github.com/sonquer/opendba) | A terminal workbench for PostgreSQL, SQL Server and SQLite. Health dashboard, tabbed SQL editor, schema browser, and an assistant that reads your schema. Every statement is parsed against the real grammar of your database and refused before it is sent if a read-only connection would be violated. |
| 654 | 3 | [holo-q/ratatui-go](https://github.com/holo-q/ratatui-go) |  |
| 655 | 2 | [profullstack/agentbbs](https://github.com/profullstack/agentbbs) | A modern BBS over SSH for humans and AI agents — arcade (DOOM), agent game ladders, personal Linux pods ($1/mo via CoinPay), and the AgentAd marketplace. Go + wish + bubbletea. |
| 656 | 2 | [northwood-labs/cli-helpers](https://github.com/northwood-labs/cli-helpers) | Helpers for Cobra and Bubbletea which follow Northwood Labs standards for CLI apps. |
| 657 | 2 | [joeel561/golang-pomodoro](https://github.com/joeel561/golang-pomodoro) | with bubbletea charm |
| 658 | 2 | [vertexE/task](https://github.com/vertexE/task) | Task management in bubbletea |
| 659 | 2 | [Raphexion/minibubbletea](https://github.com/Raphexion/minibubbletea) | A minimal bubbletea program. |
| 660 | 2 | [jhowrez/tui-hotreload](https://github.com/jhowrez/tui-hotreload) | Simple hotreload for TUI application (bubbletea) |
| 661 | 2 | [opd-ai/mtox](https://github.com/opd-ai/mtox) | TUI tox client using BubbleTea |
| 662 | 2 | [ProggerX/go-power-menu](https://github.com/ProggerX/go-power-menu) | Simple powermenu in go using bubbletea |
| 663 | 2 | [bilalyazicioglu/Currency-converter](https://github.com/bilalyazicioglu/Currency-converter) | CLI app made with BubbleTea, Golang |
| 664 | 2 | [xorsense/bubbletea_router_example](https://github.com/xorsense/bubbletea_router_example) | An example of a router working in BubbleTea |
| 665 | 2 | [rigerc/bubbletea-v2-scaffold](https://github.com/rigerc/bubbletea-v2-scaffold) |  |
| 666 | 2 | [topfunky/tiny-timer](https://github.com/topfunky/tiny-timer) | CLI task timer using BubbleTea TUI framework  |
| 667 | 2 | [melisapo/kancli](https://github.com/melisapo/kancli) | basic kanban cli built with go and bubbletea |
| 668 | 2 | [seletz/odoo-work-cli](https://github.com/seletz/odoo-work-cli) |  |
| 669 | 2 | [nitrictech/pearls](https://github.com/nitrictech/pearls) | TUI helper library and components for bubbletea programs |
| 670 | 2 | [olomix/bubblesnake](https://github.com/olomix/bubblesnake) | silly tty snake game to play with bubbletea library |
| 671 | 2 | [GianniBYoung/ghost-ship](https://github.com/GianniBYoung/ghost-ship) | Transmission TUI Client Using Writting in Go Using BubbleTea |
| 672 | 2 | [cweinberger/tmux-connect](https://github.com/cweinberger/tmux-connect) | Futuristic TUI for managing remote tmux sessions. Built with Bubbletea. |
| 673 | 2 | [go-i2p/i2ptui](https://github.com/go-i2p/i2ptui) | Embeddable TUI and freestanding CLI for I2P using BubbleTea |
| 674 | 2 | [mikeschinkel/go-tealeaves](https://github.com/mikeschinkel/go-tealeaves) | Collection of Go modules for use in BubbleTea apps |
| 675 | 2 | [Icodextrin/pomodoro-timer](https://github.com/Icodextrin/pomodoro-timer) | A TUI Pomodoro timer written in Go using Bubbletea / Lipgloss |
| 676 | 2 | [cfbender/panopticon](https://github.com/cfbender/panopticon) | a customizable command runner on file change powered by bubbletea |
| 677 | 2 | [LinPr/lazys3](https://github.com/LinPr/lazys3) | A S3 TUI for navigating objects, developed with bubbletea framework |
| 678 | 2 | [melisapo/go-shopping-list](https://github.com/melisapo/go-shopping-list) | simple shopping list tui for learning go and bubbletea framework |
| 679 | 2 | [marslo/mtui](https://github.com/marslo/mtui) | Terminal UI prompts with customizable placeholder styling, powered by Bubbletea |
| 680 | 2 | [halsten-dev/bubblehelp](https://github.com/halsten-dev/bubblehelp) | A manager to render, contextualize and manage BubbleTea keybinds globally |
| 681 | 2 | [jin-ttao/resumer](https://github.com/jin-ttao/resumer) | Browse & resume your Claude Code / Codex sessions — one picker, zero dependencies |
| 682 | 2 | [dinzz005/gitease](https://github.com/dinzz005/gitease) | A terminal-based interactive Git assistant built with Go and Bubbletea. |
| 683 | 2 | [iowarp/gact-tui](https://github.com/iowarp/gact-tui) | Generic Agentic TUI, a Bubbletea frontend that drives any compliant backend |
| 684 | 2 | [christopher-kleine/bubble-games](https://github.com/christopher-kleine/bubble-games) | Play games over an SSH connection using the Bubbletea and Wish libraries |
| 685 | 2 | [TheIncredibleMulk/go-tui-todo](https://github.com/TheIncredibleMulk/go-tui-todo) | learning about bubbletea, lipgloss, and all the other prettyness from charm.sh |
| 686 | 2 | [NLaundry/zboxes](https://github.com/NLaundry/zboxes) | Learning some Go by building a ZFS tui using bubbletea and lipgloss |
| 687 | 2 | [plutonium-239/listExtensions](https://github.com/plutonium-239/listExtensions) | bubbletea list extension bubbles - scrolling list and basic list w/o pagination |
| 688 | 2 | [406-mot-acceptable/lmtm](https://github.com/406-mot-acceptable/lmtm) | Interactive SSH tunnel builder for MikroTik and Ubiquiti gateways (Go, Bubbletea TUI) |
| 689 | 2 | [RajeshkannanRamakrishnan/lv](https://github.com/RajeshkannanRamakrishnan/lv) | A fast and interactive command-line log viewer built with Go and Bubbletea. |
| 690 | 2 | [tedilabs/ota](https://github.com/tedilabs/ota) | ♥️  k9s for Okta — a Go/Bubbletea TUI to inspect Okta admin resources |
| 691 | 2 | [FullFran/claudeops-tui](https://github.com/FullFran/claudeops-tui) | Local TUI to track Claude Code usage, costs, and tasks (Go + Bubbletea + SQLite) |
| 692 | 2 | [Grubba27/todo-list](https://github.com/Grubba27/todo-list) | Todo-list made in a cli with Go cli framework Charm.sh (bubbletea) |
| 693 | 2 | [pabloduke/naviterm](https://github.com/pabloduke/naviterm) | TUI that is simple, handles menu navigation, keeps you in control.  The Anti-BubbleTea |
| 694 | 2 | [sarkarshuvojit/bubblebook](https://github.com/sarkarshuvojit/bubblebook) | A Storybook-inspired development tool for building and testing Bubble Tea components in isolation. |
| 695 | 2 | [Meenachinmay/dev-pay-client](https://github.com/Meenachinmay/dev-pay-client) | I am building a instant payment application for developers using golang, bubbleTea, postgres and Tigerbeetle.  |
| 696 | 2 | [rhajizada/donezo](https://github.com/rhajizada/donezo) | Simple TUI to-do app in written Go using Bubble Tea and SQLite. |
| 697 | 2 | [justincordova/seshr](https://github.com/justincordova/seshr) | Replay, inspect, and prune AI agent conversation sessions in the terminal. |
| 698 | 2 | [arafatamim/bitwarden-tui](https://github.com/arafatamim/bitwarden-tui) | Bitwarden TUI with password viewing implemented (still WIP) |
| 699 | 2 | [NDOY3M4N/gopass](https://github.com/NDOY3M4N/gopass) | A super basic CLI app for generating password 😁. |
| 700 | 2 | [anotherhadi/ilovetui](https://github.com/anotherhadi/ilovetui) | A minimal Go library that provides a shared Base16color theme for terminal UIs built with bubbletea and lipgloss. |
| 701 | 2 | [DonAlexandro/Gote](https://github.com/DonAlexandro/Gote) | A terminal-based note taking app built with Go |
| 702 | 2 | [miniaturebase/pearls](https://github.com/miniaturebase/pearls) | A set of reusable terminal UI components for the Bubbletea framework and a companion to the Bubbles component library! |
| 703 | 2 | [betofigueiredo/BJJ-Instructional](https://github.com/betofigueiredo/BJJ-Instructional) | One attack and one defense technique to inspire your training! |
| 704 | 2 | [Nybkox/lazyopenconnect](https://github.com/Nybkox/lazyopenconnect) | A TUI for managing OpenConnect VPN connections. Inspired by LazyGit. |
| 705 | 2 | [minhajul/docker-cleaner](https://github.com/minhajul/docker-cleaner) | This is a terminal user interface (TUI) application built with Go and BubbleTea for cleaning up Docker images and containers. |
| 706 | 2 | [huhndev/godmarc](https://github.com/huhndev/godmarc) | TUI DMARC report analyzer |
| 707 | 2 | [sectore/fit-activities-tui](https://github.com/sectore/fit-activities-tui) | `FIT` activity data in your terminal. |
| 708 | 2 | [marcantoineg/planetary-observation](https://github.com/marcantoineg/planetary-observation) | A simple TUI to visualise CSV data. 🔭🪐 |
| 709 | 2 | [thatstoasty/sheets](https://github.com/thatstoasty/sheets) | D&D character sheet builder TUI and Web UI built using Golang, Echo, and Bubble Tea! |
| 710 | 2 | [gxstxxv/Schmierblatt](https://github.com/gxstxxv/Schmierblatt) | Welcome to Schmierblatt, a lightweight editor designed for quick, short, and easily reachable notes. As the name "Schmierblatt" (a German word meaning "scrap paper") suggests, this app provides a simple and efficient way to jot down your thoughts and ideas. |
| 711 | 2 | [richardnascimento18/devdock](https://github.com/richardnascimento18/devdock) | DevDock scans your project roots, lets you navigate your entire workspace in a keyboard-driven TUI, and drops you straight into a fully configured tmux session (windows, panes, and commands) ready to go. No config files to edit, no scripts to maintain. |
| 712 | 2 | [TesfaAsmara/tar-cli](https://github.com/TesfaAsmara/tar-cli) | A TUI for the tar command: generate tar commands with the bat of an eye |
| 713 | 2 | [sivepanda/carya](https://github.com/sivepanda/carya) | wip: Go TUI to supercharge your Git experience. |
| 714 | 2 | [TahaTesser/go-cli-template](https://github.com/TahaTesser/go-cli-template) | A modern command-line interface template built with Bubble Tea for interactive TUI functionality and Lipgloss for beautiful styling. |
| 715 | 2 | [nexuls/mcworker-go](https://github.com/nexuls/mcworker-go) | 🌍 Create, manage, and monitor Minecraft servers from your terminal — supports Paper, Fabric, Forge, and more. |
| 716 | 2 | [rubiojr/eyez](https://github.com/rubiojr/eyez) | glossy-ly experimental proxy 💄☣️🖧 |
| 717 | 2 | [moKshagna-p/cinder](https://github.com/moKshagna-p/cinder) | music visualization from your terminal |
| 718 | 2 | [Karim-W/toastman](https://github.com/Karim-W/toastman) | TUI Http Client To make HTTP Requests  |
| 719 | 2 | [jwr1/pixelstream](https://github.com/jwr1/pixelstream) | Stream videos to your awtrix clock with ease |
| 720 | 2 | [craigderington/lazystack](https://github.com/craigderington/lazystack) | A unified Terminal User Interface (TUI) for managing Kubernetes resources. Built with Go and Bubbletea.  Supports microk8s, k3s, minikube or any cluster with a KUBECONFIG. |
| 721 | 2 | [aadithpm/bookraid](https://github.com/aadithpm/bookraid) | Download fiction from Libgen with a CLI |
| 722 | 2 | [Eduardo79Silva/taskarena](https://github.com/Eduardo79Silva/taskarena) | A pull-based architecture for task management |
| 723 | 2 | [corbinmemo/portal-svc](https://github.com/corbinmemo/portal-svc) | A cross-platform daemon for sing-box |
| 724 | 2 | [mtyurt/ecstui](https://github.com/mtyurt/ecstui) | Opinionated ECS TUI |
| 725 | 2 | [Borderliner/emacs-installer](https://github.com/Borderliner/emacs-installer) | A good-looking Bubble Tea TUI that compiles GNU Emacs from source and installs it across mainstream Linux distros (and macOS). |
| 726 | 2 | [AlNaheyan/termchat](https://github.com/AlNaheyan/termchat) | Lightweight terminal-based chat |
| 727 | 2 | [fontainecoutino/dont-sleep](https://github.com/fontainecoutino/dont-sleep) | Prevent system from sleeping |
| 728 | 2 | [David-mwas/WindowsTempCleanify](https://github.com/David-mwas/WindowsTempCleanify) | A command-line tool and script written in Go that automates the cleanup of temporary files on Windows. It cleans your user TEMP directory, Windows system temporary directory (%temp%) and Prefetch directory. |
| 729 | 2 | [bpicode/tmus](https://github.com/bpicode/tmus) | tmus is a terminal music player written in Go. |
| 730 | 2 | [alexferl/myrient-browser](https://github.com/alexferl/myrient-browser) | A terminal-based file browser and downloader for Myrient, featuring concurrent downloads, resume support, automatic extraction, and a clean TUI interface built with Bubble Tea. |
| 731 | 2 | [crper/tqrx](https://github.com/crper/tqrx) | Terminal-first QR generator with CLI/TUI workflows, live preview, and PNG/SVG export. |
| 732 | 2 | [senisia/mpgo](https://github.com/senisia/mpgo) | simple mp3 player with go, plays in terminal, built with beep and bubbletea, music files should be in ~/music, i dont know the equivalent of tilde in windows |
| 733 | 2 | [howardhenrystephen/TerminalReader](https://github.com/howardhenrystephen/TerminalReader) | A terminal-based novel reader TUI app built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Search, download, manage and read web novels comfortably from the command line. |
| 734 | 2 | [Jayesh-Dev21/CalledCode](https://github.com/Jayesh-Dev21/CalledCode) | A professional AI coding assistant with a minimal TUI interface and seamless VS Code integration and context window |
| 735 | 2 | [sangdth/randomport](https://github.com/sangdth/randomport) | Learning project: Generating random ports within range |
| 736 | 2 | [alkanoidev/journalismus-cli](https://github.com/alkanoidev/journalismus-cli) | Capture thoughts effortlessly in the command line. |
| 737 | 2 | [kaandesu/bubble-chat](https://github.com/kaandesu/bubble-chat) | Basic TCP chat rooms from terminal |
| 738 | 2 | [acrobatstick/gisting](https://github.com/acrobatstick/gisting) | Interactive gists management inside the terminal |
| 739 | 2 | [Kumneger0/cheztui](https://github.com/Kumneger0/cheztui) | A tiny chezmoi wrapper that makes super easy to manage your dotfiles |
| 740 | 2 | [svenrisse/gomodoro](https://github.com/svenrisse/gomodoro) | Pomodoro Timer TUI |
| 741 | 2 | [YeyoM/better-commits](https://github.com/YeyoM/better-commits) | My CLI tool for creating commits using the conventional commits standard |
| 742 | 2 | [arthursfares/skulls](https://github.com/arthursfares/skulls) | Terminal-based voice/text chat over WebRTC. A Bubble Tea TUI client connects peer-to-peer for audio (Opus) and messaging, using a lightweight Go signaling server to broker room handshakes. |
| 743 | 2 | [kiraxbt/btcx](https://github.com/kiraxbt/btcx) | Terminal UI for managing Bitcoin wallets — built in Go |
| 744 | 2 | [jasonuc/usermakertui](https://github.com/jasonuc/usermakertui) | interactive bubble tea form |
| 745 | 2 | [twtrubiks/ptt-spider-go](https://github.com/twtrubiks/ptt-spider-go) | ptt-spider-go 圖片爬蟲(表特板  八卦板) |
| 746 | 2 | [TheNeovimmer/devsentinel](https://github.com/TheNeovimmer/devsentinel) | A modern, all-in-one Terminal User Interface (TUI) for project intelligence - combining architecture analysis, runtime monitoring, git insights, code quality metrics, and real-time log streaming. |
| 747 | 2 | [iv4n-ga6l/Go-Computer-system-resources-monitoring](https://github.com/iv4n-ga6l/Go-Computer-system-resources-monitoring) | Computer system resources monitoring in Go |
| 748 | 2 | [mwiater/agon](https://github.com/mwiater/agon) | A terminal-first companion for Ollama users. Agon helps you manage multiple hosts, compare models, chat with models, and synchronize your model inventory across machines. |
| 749 | 2 | [vnsonvo/pomodoro-cli](https://github.com/vnsonvo/pomodoro-cli) | Pomodoro CLI |
| 750 | 2 | [damienp199/bagent](https://github.com/damienp199/bagent) | Lanceur de workspaces en terminal (TUI) pour ouvrir un dossier dans VSCode, Claude Code ou Codex |
| 751 | 2 | [leschuster/deepl-cli](https://github.com/leschuster/deepl-cli) | A terminal translator utilizing the power of DeepL |
| 752 | 2 | [michaelhass/gitglance](https://github.com/michaelhass/gitglance) | Git terminal UI written in Go |
| 753 | 2 | [bloodynite/lazyredis](https://github.com/bloodynite/lazyredis) | Terminal UI for browsing and editing Redis keys |
| 754 | 2 | [arunaiDeepan/rss-genie](https://github.com/arunaiDeepan/rss-genie) | A simple terminal-based RSS reader that lets you configure your favorite feeds and read them directly in your terminal |
| 755 | 2 | [xRiErOS/beans-tui](https://github.com/xRiErOS/beans-tui) | PO-cockpit TUI for beans repos — an independent client on top of hmans/beans |
| 756 | 2 | [kurojs/EchoRoutine](https://github.com/kurojs/EchoRoutine) | AI-powered daily schedule announcer. Define time blocks in the TUI editor (Bubbletea/Go), get each block announced via ElevenLabs TTS with custom AI motivation. Runs on a systemd timer — your schedule, spoken aloud. |
| 757 | 2 | [huoyijie/GoChat](https://github.com/huoyijie/GoChat) | a chat program with golang and protobuf |
| 758 | 2 | [jrniemiec/lore](https://github.com/jrniemiec/lore) | lore is a terminal TUI chat app in Go replacing ask-repl. Built on Bubbletea, it supports Anthropic, OpenAI, and Ollama providers with streaming.   Conversations live in ~/.lore/, compatible with ask's storage format.  |
| 759 | 2 | [indrasvat/gh-ghent](https://github.com/indrasvat/gh-ghent) | A GitHub CLI extension for agents monitoring PRs |
| 760 | 2 | [stefanvanburen/buftui](https://github.com/stefanvanburen/buftui) | TUI for the Buf Schema Registry |
| 761 | 2 | [prit342/kns](https://github.com/prit342/kns) | CLI tool to interactively switch kubernetes namespaces |
| 762 | 2 | [jordangarrison/cloud-build-cli](https://github.com/jordangarrison/cloud-build-cli) | A TUI to select and stream Google Cloud Build logs in the terminal. |
| 763 | 2 | [thaddeusrhatcher/jirate](https://github.com/thaddeusrhatcher/jirate) | Jirate is a CLI tool for managing Jira comments. |
| 764 | 2 | [Namacha411/scpick](https://github.com/Namacha411/scpick) | Interactive cross-platform SCP/SFTP file transfer TUI (dual-pane, vim-style keys) |
| 765 | 2 | [fdddf/pixa](https://github.com/fdddf/pixa) | Search and download Pixabay images and videos from the terminal — TUI browser, batch mode, in-terminal previews, dedupe |
| 766 | 2 | [ch1lam/aice-cli](https://github.com/ch1lam/aice-cli) | A small, batteries-included open-source AI coding agent for your terminal—one Go binary with multi-provider LLM support, guarded execution, and recoverable branching sessions.一个小而完整的终端 AI coding agent：以单个 Go 二进制交付，支持多 Provider，内置安全门控与可恢复、可分支的 Session。 |
| 767 | 2 | [shonenm/live-pr](https://github.com/shonenm/live-pr) | Living pull request for LLM-assisted development: agent-fed decision timeline + local PR-style TUI + timeline-reflecting PR export |
| 768 | 2 | [lemantorus/RepoScope](https://github.com/lemantorus/RepoScope) |     🔭 RepoScope — A high-performance TUI project discovery tool. Instantly scan your drive to find, sort, and audit your coding projects. Tracks Git status, project types, and disk usage in a sleek terminal interface. Built with Go. |
| 769 | 2 | [hrvadl/gowatchsql](https://github.com/hrvadl/gowatchsql) | SQL TUI |
| 770 | 2 | [TECHedgehog/boltx](https://github.com/TECHedgehog/boltx) | A TUI tool for setting up and managing Linux systems. Built with Go + Bubbletea, it guides you through system configuration, security hardening, and software installation — whether you're spinning up a VPS, a homelab, or a dev VM. |
| 771 | 2 | [m1ck6x/PwdMan](https://github.com/m1ck6x/PwdMan) | A simplistic password manager with a terminal UI supported on both the Windows and Linux operating systems. |
| 772 | 2 | [kanywst/rapg](https://github.com/kanywst/rapg) | Local-first secret manager for the AI-agent era — keep API keys out of .env files and out of your agent transcripts. |
| 773 | 2 | [xiguayiqiu/termd](https://github.com/xiguayiqiu/termd) | termd 是一个运行于终端的轻量级 Markdown 编辑器，基于 bubbletea 框架构建，提供编辑与实时富文本预览双模式。 |
| 774 | 2 | [stussysenik/vfx](https://github.com/stussysenik/vfx) | Terminal-native 2D graphics engine — 15 GPU-quality animations at 60fps using braille characters + true color |
| 775 | 2 | [fezcode/atlas.conquistador](https://github.com/fezcode/atlas.conquistador) | A beautiful terminal-based file explorer for the Atlas Suite with multi-selection, file operations, internal text viewer, scrolling, system-default file opening, and a modern, boxed UI. |
| 776 | 2 | [vanderheijden86/b9s](https://github.com/vanderheijden86/b9s) | TUI for beads tasks. Forked from beads_viewer, inspired by k9s. Tree-view, editing and project switcher added. |
| 777 | 2 | [Emmyme/file-organizer](https://github.com/Emmyme/file-organizer) | A colourful, interactive TUI app for organizing files by type, date, or size. Built with Go, Bubble Tea, and Lipgloss. |
| 778 | 2 | [zigai/zgod](https://github.com/zigai/zgod) | Interactive shell history search with fuzzy, regex, and glob matching |
| 779 | 2 | [aikoncwd/ytcli](https://github.com/aikoncwd/ytcli) | Reproductor de música TUI para YouTube (solo audio), escrito en Go con mpv y yt-dlp |
| 780 | 2 | [burgr033/flashMe](https://github.com/burgr033/flashMe) | terminal based flash card application written in Go utilizing the awesome bubble tea framework |
| 781 | 2 | [fntune/turbobook](https://github.com/fntune/turbobook) | Terminal notebook editor with vim keybindings and a live Python kernel |
| 782 | 2 | [Soapalin/IdleReader](https://github.com/Soapalin/IdleReader) | IdleReader: A Text-Based Terminal Game |
| 783 | 2 | [binRick/go-extrace-parser](https://github.com/binRick/go-extrace-parser) | Go CLI to parse extrace process-execution logs into an interactive Bubble Tea TUI table |
| 784 | 2 | [mkrowiarz/ccmonitor](https://github.com/mkrowiarz/ccmonitor) | Terminal dashboard for monitoring Claude Code usage, sessions, and rate limits |
| 785 | 2 | [rahji/speachy](https://github.com/rahji/speachy) | A command-line program for returning words, from a text file, that match specific parts of speech tags |
| 786 | 2 | [pfrederiksen/rivian-ls](https://github.com/pfrederiksen/rivian-ls) | Production-quality TUI + headless CLI for monitoring Rivian vehicle telemetry via unofficial API |
| 787 | 2 | [lunemis/skimd](https://github.com/lunemis/skimd) | TUI markdown browser for skimming docs without leaving the terminal |
| 788 | 2 | [StaiLee/Ikelos](https://github.com/StaiLee/Ikelos) | 🧬 Ikelos — tactical website cloner & web-recon engine in Go with a slick Bubble Tea TUI, powered by goquery. |
| 789 | 2 | [qrxnz/gk](https://github.com/qrxnz/gk) | A terminal-based task and habit manager written in Go |
| 790 | 2 | [cmj0121/baton](https://github.com/cmj0121/baton) | A terminal multiplexer for AI coding agents: one keyboard cockpit for the whole fleet, with work items, a task queue, and per-panel resource caps. |
| 791 | 2 | [lakerszhy/ght](https://github.com/lakerszhy/ght) | A TUI for browsing GitHub trending repositories |
| 792 | 2 | [chardoncs/downjack](https://github.com/chardoncs/downjack) | Set up your gitignore and license files like using a lumberjack (or a down jacket) |
| 793 | 2 | [dwgx/SKIBOX](https://github.com/dwgx/SKIBOX) | Local API key manager and config switcher for AI developer tools |
| 794 | 2 | [marcelblijleven/gh-hookshot](https://github.com/marcelblijleven/gh-hookshot) | A TUI extension for the GitHub CLI to inspect webhook deliveries |
| 795 | 2 | [fmo/skytui](https://github.com/fmo/skytui) | A focused terminal Pomodoro timer for macOS, built with Go and Bubble Tea. |
| 796 | 2 | [dev-core-busy/winq](https://github.com/dev-core-busy/winq) | Minimalist TUI Windows agent: translate natural language to PowerShell commands via a local LLM |
| 797 | 2 | [he8um/daryaft](https://github.com/he8um/daryaft) | A fast, polished TUI/CLI downloader written in Go |
| 798 | 2 | [gianni-labs/count-focus](https://github.com/gianni-labs/count-focus) | A terminal-native focus timer for deep work. Pomodoro, countdowns, goals, and automation hooks. |
| 799 | 2 | [ninanung/csm](https://github.com/ninanung/csm) | A small TUI for browsing and resuming Claude Code sessions — rich identification, auto cwd/branch alignment, i18n (en/ko). |
| 800 | 2 | [mojomast/geoffrussy](https://github.com/mojomast/geoffrussy) | Go 1.24 AI-driven software delivery orchestrator that guides projects through a staged pipeline (interview → design → plan → review → develop). It uses SQLite (via mattn/go-sqlite3, requires CGO), Cobra for CLI, Bubbletea for TUI, and supports 16 AI providers via an OpenAI-compatible interface pattern. |
| 801 | 2 | [binbandit/goshed](https://github.com/binbandit/goshed) | A powerful, interactive Go playground manager with a beautiful TUI - create, manage, and organize your Go experiments with style 🚀 |
| 802 | 2 | [cicbyte/answer-cli](https://github.com/cicbyte/answer-cli) | Apache Answer Q&A community CLI tool — TUI browser, CLI operations, AI chat, and MCP integration in your terminal |
| 803 | 2 | [blakemister/qs](https://github.com/blakemister/qs) | Quickly open any project with your AI coding tool of choice. Terminal launcher TUI for Claude Code, Codex, Gemini, and more. |
| 804 | 2 | [erik-adelbert/flame](https://github.com/erik-adelbert/flame) | A high-performance DOOM fire animation for the terminal, written in Go with Bubble Tea and Lip Gloss. |
| 805 | 2 | [googlesky/sstop](https://github.com/googlesky/sstop) | Real-time per-process network bandwidth monitor for the terminal |
| 806 | 2 | [Zayan-Mohamed/nova](https://github.com/Zayan-Mohamed/nova) | A terminal-based WiFi and LAN security assessment tool with active risk scoring. |
| 807 | 2 | [erik-adelbert/firework](https://github.com/erik-adelbert/firework) | A high-performance firework animation for the terminal, written in Go with Bubble Tea and Lip Gloss. |
| 808 | 2 | [phlx0/ghscope](https://github.com/phlx0/ghscope) | Analyse any GitHub repository from your terminal — stars, contributors, churn, deps, PR velocity and more |
| 809 | 2 | [olgeni/facl](https://github.com/olgeni/facl) | Terminal UI and CLI for FreeBSD NFSv4 ACLs, modelled after the Windows Advanced Security Settings dialog — inheritance scopes, effective access, ZFS aclinherit/aclmode awareness, chmod preview; verified against the kernel. |
| 810 | 2 | [bachtiarpanjaitan/ihand-tui](https://github.com/bachtiarpanjaitan/ihand-tui) | Terminal User Interface AI agent using ihandai-go |
| 811 | 2 | [juanhuttemann/monkey-cli](https://github.com/juanhuttemann/monkey-cli) |  Chat with Claude from your terminal — interactive TUI, one-shot pipes, or agentic file operations. 🐒 |
| 812 | 2 | [JeiKeiLim/claude-code-log-viewer-cli](https://github.com/JeiKeiLim/claude-code-log-viewer-cli) | Terminal-based viewer for Claude Code conversation logs. Browse projects, view conversations with vim-style navigation, collapsible thinking blocks, and CJK character support. |
| 813 | 2 | [dohwi/tmux-manager](https://github.com/dohwi/tmux-manager) | TUI session manager for tmux — define workspaces in YAML, restore on reboot |
| 814 | 2 | [fezcode/atlas.cam](https://github.com/fezcode/atlas.cam) | A terminal-based webcam viewer and ASCII camera with real-time edge detection, image filters, and GIF recording. Part of the Atlas Suite. |
| 815 | 2 | [erik-adelbert/donut](https://github.com/erik-adelbert/donut) | A high-performance donut animation for the terminal, written in Go with Bubble Tea and Lip Gloss. |
| 816 | 2 | [fezcode/atlas.clock](https://github.com/fezcode/atlas.clock) | A high-visibility, multi-timezone world clock dashboard with millisecond precision for the terminal. Part of the Atlas Suite. |
| 817 | 2 | [Laurent00TT/slimproxy](https://github.com/Laurent00TT/slimproxy) | A slim LLM reverse proxy with a terminal dashboard: OpenAI / Anthropic / Gemini compatible endpoints backed by your own OAuth credentials, plus Cloudflare Tunnel lifecycle and honest deployment diagnostics. Built on CLIProxyAPI's SDK. |
| 818 | 2 | [DementevVV/commitsum](https://github.com/DementevVV/commitsum) | An interactive terminal TUI for exploring and summarizing GitHub commits. Built in Go with Bubble Tea and GitHub CLI, featuring multi-repo support, filtering, statistics, and export to text, Markdown, or JSON. |
| 819 | 2 | [fezcode/atlas.grave](https://github.com/fezcode/atlas.grave) | A high-fidelity, interactive process reaper for the Mojave wasteland. Hunt down resource-heavy processes ("Restless Souls") and bury them in an Onyx graveyard. Reclaim your system's memory with a Pip-Boy inspired interface and a heavy-metal aesthetic. |
| 820 | 2 | [agentcarto/agentcarto](https://github.com/agentcarto/agentcarto) | Terminal UI that brings the local sessions of all your AI coding agents — Claude Code, Codex, Grok, and GitHub Copilot Chat — into one searchable place to browse, search, resume, and fork. |
| 821 | 2 | [olekgolus11/SliceDiff](https://github.com/olekgolus11/SliceDiff) | AI-assisted terminal UI for slicing large GitHub pull request diffs into reviewable chunks. |
| 822 | 2 | [fezcode/atlas.facade](https://github.com/fezcode/atlas.facade) | An instant, high-fidelity mock API server for the Mojave wasteland. Forge your backend infrastructure from simple blueprints (PIML) while the real world burns. Features real-time request logging, latency simulation, and a rugged Pip-Boy inspired TUI. |
| 823 | 2 | [fezcode/atlas.todo](https://github.com/fezcode/atlas.todo) | A fast, minimalist, local-first TUI task manager with Vim-bindings, smart grouping, and CLI quick-add support. |
| 824 | 2 | [cicbyte/anki-cli](https://github.com/cicbyte/anki-cli) | A CLI tool for Anki — manage cards, notes, and decks from the terminal, with MCP Server and Skill integration for AI clients. / Anki 命令行工具 — 在终端管理卡片、笔记和牌组，支持 MCP Server 和 Skill 集成，可对接 AI 客户端。 |
| 825 | 2 | [hoppxi/bpv](https://github.com/hoppxi/bpv) | local music player |
| 826 | 2 | [muhammadSaad63/ConnectFourGO](https://github.com/muhammadSaad63/ConnectFourGO) | A high-performance Connect Four engine and TUI built in Go. Refactored from C to leverage concurrency (Goroutines) for AI and the Bubble Tea framework for a reactive terminal experience. :D |
| 827 | 2 | [kyungw00k/upbit](https://github.com/kyungw00k/upbit) | 📈 AI-native CLI for Upbit cryptocurrency exchange — market data, trading, real-time TUI, candle cache, i18n (ko/en) |
| 828 | 2 | [Ilya-sss/ToDo-Tui-App](https://github.com/Ilya-sss/ToDo-Tui-App) | 🎨 A beautiful terminal To-Do app built with Go, Bubble Tea & Lip Gloss. Keyboard-driven, distraction-free task management right in your terminal! |
| 829 | 2 | [tom-draper/typing-speed](https://github.com/tom-draper/typing-speed) | A command-line typing speed tester. |
| 830 | 2 | [tidefly-oss/tidefly-tui](https://github.com/tidefly-oss/tidefly-tui) | Interactive terminal installer for Tidefly — sets up the full stack on any Linux server in minutes. Zero config, guided TUI built with Bubble Tea. Supports rollback, health checks and rootless install. |
| 831 | 2 | [fezcode/atlas.stats](https://github.com/fezcode/atlas.stats) | A fast, cross-platform terminal user interface (TUI) for real-time system monitoring. built with Go, Bubble Tea, and Lip Gloss. Monitor CPU, Memory, multi-disk usage, network throughput, and top processes in a beautiful, responsive layout. |
| 832 | 2 | [fezcode/atlas.screensaver](https://github.com/fezcode/atlas.screensaver) | A collection of nostalgic and aesthetic terminal screensavers including vintage Pipes and Starfield animations. Built with Go and Bubble Tea. Part of the Atlas Suite. |
| 833 | 2 | [lakerszhy/rssx](https://github.com/lakerszhy/rssx) | A RSS feed reader in terminal. |
| 834 | 2 | [fezcode/atlas.games](https://github.com/fezcode/atlas.games) | Wilson's Revenge 🐔 - A high-octane(!?), horizontal terminal runner built with Go. Dodge cars, blast through human enemies with a shotgun, and manage game-altering power-ups in a beautiful ASCII-art world. Fast, challenging, and built for the terminal. |
| 835 | 2 | [Joncik91/inflate](https://github.com/Joncik91/inflate) | Context-aware CLI prompt inflater for Claude Code. Type a fragment, get a context-loaded prompt. |
| 836 | 2 | [fezcode/atlas.otp](https://github.com/fezcode/atlas.otp) | 🔐 Minimalist, high-visibility terminal TOTP (2FA) manager with a signature Onyx & Gold aesthetic. Part of the Atlas Suite. |
| 837 | 2 | [fezcode/atlas.compass](https://github.com/fezcode/atlas.compass) | A secure, local-first terminal password manager. Zero-knowledge AES-256-GCM encryption, Argon2id key derivation, and a high-visibility TUI. Part of the Atlas Suite. |
| 838 | 2 | [aafeher/nftui](https://github.com/aafeher/nftui) | Terminal UI for the Linux nftables firewall — browse and edit rules, sets, maps and named objects without touching the nft CLI |
| 839 | 2 | [JasonLovesDoggo/roastme](https://github.com/JasonLovesDoggo/roastme) | A witty CLI tool that analyzes your shell history and brutally roasts your command-line habits using AI to keep your ego in check.  |
| 840 | 2 | [fezcode/atlas.deck](https://github.com/fezcode/atlas.deck) | An interactive TUI "Stream Deck" for your terminal. Organize, trigger, and monitor complex project workflows through a customizable grid of pads mapped to shell commands. Part of the Atlas Suite. |
| 841 | 2 | [DementevVV/sshwatch](https://github.com/DementevVV/sshwatch) | A modern CLI and TUI tool for monitoring SSH logins on remote servers with instant Telegram alerts. Built with Go, featuring secure agent deployment, PAM integration, and a beautiful terminal UI. |
| 842 | 2 | [rmpato/pogo](https://github.com/rmpato/pogo) | curl, but it remembers. One binary: it makes the request and keeps the answer — then gives you a terminal UI over everything you have run, grouped by API, with staging and production as environments of the same thing. |
| 843 | 2 | [yxxbc/spark-store-tui](https://github.com/yxxbc/spark-store-tui) | Spark Store TUI：Go + Bubble Tea 原生终端应用商店，支持浏览与搜索、图片预览、断点续传、安装/卸载以及 Debian、RPM、AUR 多架构发布。 |
| 844 | 2 | [turkprogrammer/sql-top](https://github.com/turkprogrammer/sql-top) | A real-time SQL profiler and monitor TUI for PostgreSQL. High-performance "htop" for your database with instant EXPLAIN, wait event analysis, and delta highlighting. |
| 845 | 2 | [fezcode/atlas.bench](https://github.com/fezcode/atlas.bench) | 🏎️ High-precision, multi-command TUI benchmarking tool with statistical rigor and side-by-side comparisons. Part of the Atlas Suite. |
| 846 | 2 | [SirSobhan0/gotodo](https://github.com/SirSobhan0/gotodo) | Simple project tracking app with golang |
| 847 | 2 | [svyatov/oz](https://github.com/svyatov/oz) | Config-driven CLI wizard framework |
| 848 | 2 | [KevinTCoughlin/lazydeck](https://github.com/KevinTCoughlin/lazydeck) | A lazydocker-style terminal UI for managing a fleet of Steam Deck and Steam Machine devkits. |
| 849 | 2 | [fmilioni/envault](https://github.com/fmilioni/envault) | Save and restore .env files with a single command — like git stash for your environment variables. Small cross-platform CLI (Go + Bubble Tea) with an interactive vault browser and encrypted, portable export/import. |
| 850 | 2 | [alexander-posztos/ccdash](https://github.com/alexander-posztos/ccdash) | Answers "where did I leave off?" across all your Claude Code projects - a per-project recap, git state, and one-key resume, in one terminal dashboard. |
| 851 | 2 | [RStephanH/sentrymesh-gateway](https://github.com/RStephanH/sentrymesh-gateway) | Simulated IoT security gateway (Go) - MQTT ingestion with validation, replay detection, and flood/rate-limit protection, backed by SQLite and a live Bubble Tea dashboard. |
| 852 | 2 | [jim-ww/kage](https://github.com/jim-ww/kage) | TUI XMPP client  |
| 853 | 2 | [mogglemoss/pelorus](https://github.com/mogglemoss/pelorus) | Dual-pane TUI file manager. Local, SFTP, and Tailscale panes behave identically. Archives as directories. Fuzzy everywhere. One static Go binary, no runtime dependencies. |
| 854 | 2 | [Luv-Goel/contextflow](https://github.com/Luv-Goel/contextflow) | Shell history that understands your workflows, not just your commands. Fuzzy search TUI, workflow detection, session replay, and project-aware history management. |
| 855 | 2 | [N-Erickson/termidar](https://github.com/N-Erickson/termidar) | Real-time weather radar in your terminal - ssh termidar.app |
| 856 | 2 | [Riximus/SBBuddy](https://github.com/Riximus/SBBuddy) | 🧭 A friendly Swiss public transport terminal app. Search timetables, plan connections, and get the connection to your phone via QR. All from your console. |
| 857 | 2 | [andrew528i/SoLock](https://github.com/andrew528i/SoLock) | Decentralized password manager on Solana. Single binary, single password, no servers |
| 858 | 2 | [shahadulhaider/restless](https://github.com/shahadulhaider/restless) | A terminal-native HTTP client with TUI. Uses .http files — the same format as JetBrains and VS Code REST Client. Import from Postman, Insomnia, Bruno, curl, OpenAPI. Response assertions for CI/CD. Code generation in 8 languages. |
| 859 | 2 | [agenticpoa/sshsign](https://github.com/agenticpoa/sshsign) | SSH signing service for AI agents. Scoped authorization, co-sign approval with handwritten signatures, immutable audit trail. |
| 860 | 2 | [8tp/netmap](https://github.com/8tp/netmap) | Visual network topology mapper and scanner TUI — discover devices, scan ports, measure latency, all in the terminal |
| 861 | 2 | [clarabennettdev/logpilot](https://github.com/clarabennettdev/logpilot) | 🪵 A fast, multi-source structured log viewer for the terminal — stream, search, and correlate JSON, logfmt, and plain text logs with color-coded output. Built with Go and Bubble Tea. |
| 862 | 2 | [dzaurov/claude-sessions](https://github.com/dzaurov/claude-sessions) | Browse and resume every Claude Code chat across every project, from any terminal. Fuzzy-search by topic, pin favorites, hide noise — Enter and the resumed session takes over your terminal. |
| 863 | 2 | [madstone-tech/ason](https://github.com/madstone-tech/ason) | 🚀 A powerful Go-based project scaffolding tool that transforms templates into fully-formed projects using Jinja2-style templating, with CLI automation and interactive prompts. |
| 864 | 2 | [s-johri/sshush](https://github.com/s-johri/sshush) | Interactive terminal UI (TUI) for SSH key management — control ssh-agent, browse & edit ~/.ssh/config, connect to hosts, generate keys, and audit permissions. Written in Go. |
| 865 | 2 | [jefrnc/sekd](https://github.com/jefrnc/sekd) | Day-trading due diligence in your terminal. Decodes SEC filings, scores dilution risk, and uses AI to surface warrants and shelf capacity from S-3s and 10-Qs. |
| 866 | 1 | [skatkov/bubbleteaTwoScreens](https://github.com/skatkov/bubbleteaTwoScreens) | Experimenting with multiple screen TUI application and how to better work with those. |
| 867 | 1 | [abulujayn/bubble-tea-test](https://github.com/abulujayn/bubble-tea-test) | Just testing out bubbletea. |
| 868 | 1 | [jhowrez/bubbletea-menus](https://github.com/jhowrez/bubbletea-menus) | Bubbletea menu/submenu propotype |
| 869 | 1 | [justwasm/boba](https://github.com/justwasm/boba) | browser-oriented bubbletea adapter |
| 870 | 1 | [alimsk/list](https://github.com/alimsk/list) | list component for bubbletea |
| 871 | 1 | [rusq/rbubbles](https://github.com/rusq/rbubbles) | TUI Components built with bubbletea |
| 872 | 1 | [code-yeongyu/bubbletea-wm](https://github.com/code-yeongyu/bubbletea-wm) | A floating window manager built with Bubbletea v2 - vibecoded while learning Go TUIs |
| 873 | 1 | [KenMwaura1/bubble](https://github.com/KenMwaura1/bubble) | Termianl apps using bubbletea package |
| 874 | 1 | [zamoosh/bubbletea_practice](https://github.com/zamoosh/bubbletea_practice) |  |
| 875 | 1 | [vedrankolka/bubbletea-tictactoe](https://github.com/vedrankolka/bubbletea-tictactoe) | Command line demo app built using Bubbletea for playing tic tac toe in the command line! |
| 876 | 1 | [katelynn620/bubbletea-playground](https://github.com/katelynn620/bubbletea-playground) |  |
| 877 | 1 | [NimbleMarkets/ntcharts-pdf](https://github.com/NimbleMarkets/ntcharts-pdf) | PDF Viewer widget for BubbleTea |
| 878 | 1 | [nmin11/bubbletea-practice](https://github.com/nmin11/bubbletea-practice) | bubbletea of Charm CLI exercise |
| 879 | 1 | [jejacks0n/bubbletea-menubar](https://github.com/jejacks0n/bubbletea-menubar) | I wanted a menubar that I could share across projects. This is that library. |
| 880 | 1 | [puttehi/tui-games](https://github.com/puttehi/tui-games) | Exploring Charms [Wish](https://github.com/charmbracelet/wish) & [Bubbletea](https://github.com/charmbracelet/bubbletea) |
| 881 | 1 | [pgavlin/bubbletea-nav](https://github.com/pgavlin/bubbletea-nav) | Screen navigation and focus management for Bubble Tea applications. Authored using Claude Code. |
| 882 | 1 | [blackwell-systems/bubbletea-commandpalette](https://github.com/blackwell-systems/bubbletea-commandpalette) | Fuzzy-search command palette overlay for Bubble Tea (VS Code Ctrl+P style) |
| 883 | 1 | [rohitrgupta/bubbletea-practice](https://github.com/rohitrgupta/bubbletea-practice) |  |
| 884 | 1 | [swinton/example-bubbletea](https://github.com/swinton/example-bubbletea) |  |
| 885 | 1 | [Bachmann1234/chatbotCli](https://github.com/Bachmann1234/chatbotCli) | Interface to chat GBT using Bubbletea |
| 886 | 1 | [leg100/bubbletea-dev-example](https://github.com/leg100/bubbletea-dev-example) |  |
| 887 | 1 | [mhb8898/bubbletea-downloader-progressbar](https://github.com/mhb8898/bubbletea-downloader-progressbar) | Bubble Tea TUI progress bar for concurrent / multi-connection downloads (IDM-style and per-worker views) |
| 888 | 1 | [WilmerLeonCh/TaskManager](https://github.com/WilmerLeonCh/TaskManager) | CLI task manager stylized with bubbletea |
| 889 | 1 | [m3talsmith/btpolling](https://github.com/m3talsmith/btpolling) | An example of polling in BubbleTea |
| 890 | 1 | [notjedi/tabs](https://github.com/notjedi/tabs) | [WIP] tabs bubble for bubbletea apps. |
| 891 | 1 | [LowZaar/Semana_Academica_UNIDEAU_2025](https://github.com/LowZaar/Semana_Academica_UNIDEAU_2025) | Introdução ao golang e o bubbletea |
| 892 | 1 | [itsanji/learn-charm-bubbletea](https://github.com/itsanji/learn-charm-bubbletea) | Learning Charmbracelet/bubbletea. Differrent branch contain differrent app |
| 893 | 1 | [bushyFUPA/bookbookbook](https://github.com/bushyFUPA/bookbookbook) | Audiobook player TUI using Bubbletea and FFmpeg |
| 894 | 1 | [zakwanhisham/go-time](https://github.com/zakwanhisham/go-time) | Simple clock using golang, bubbletea, and lipgloss |
| 895 | 1 | [dffrs/read-it](https://github.com/dffrs/read-it) | read-it - CLI tool to browse, summarise and ticketify Cypress tests |
| 896 | 1 | [timmattison/golang-bubbletea-tool-template](https://github.com/timmattison/golang-bubbletea-tool-template) | A template I use to create Bubbletea based tools |
| 897 | 1 | [darkhz/treeview-bubbletea-v2](https://github.com/darkhz/treeview-bubbletea-v2) |  |
| 898 | 1 | [Davidca089/KingsTerminal](https://github.com/Davidca089/KingsTerminal) | Chess CLI written in GO with Bubbletea  |
| 899 | 1 | [Grubba27/fetch](https://github.com/Grubba27/fetch) | Simple fetch implementeation made with @charmbracelet/bubbletea |
| 900 | 1 | [ViniZap4/lumi-tui](https://github.com/ViniZap4/lumi-tui) | Terminal UI client for lumi — Go + Bubbletea |
| 901 | 1 | [emarifer/go-cli-bubbletea-todoapp](https://github.com/emarifer/go-cli-bubbletea-todoapp) | Command line application (CLI Todo App) made with Bubble Tea (to create a TUI) & Cobra frameworks, and CRUD to a SQLite database |
| 902 | 1 | [SuperInstance/flux-tui](https://github.com/SuperInstance/flux-tui) | FLUX VM Debugger & Conformance Dashboard — Go + bubbletea |
| 903 | 1 | [yuuohg/musger](https://github.com/yuuohg/musger) | tui wrapper around mpv, made with bubbletea |
| 904 | 1 | [Aayushstha03/hypr-breaktimer](https://github.com/Aayushstha03/hypr-breaktimer) | Go + BubbleTea TUI break reminder for Hyprland. |
| 905 | 1 | [mieubrisse/bubble-bath](https://github.com/mieubrisse/bubble-bath) | A component system for Charm's BubbleTea framework |
| 906 | 1 | [codePriyanshuRajAnand/cliNoteTakingApp](https://github.com/codePriyanshuRajAnand/cliNoteTakingApp) | Created in GoLang using Bubbletea and supporting libraries |
| 907 | 1 | [ichbinbekir/tearouter](https://github.com/ichbinbekir/tearouter) | Model routing for golang bubbletea TUI framework 🔀 |
| 908 | 1 | [darkliquid/charmingui](https://github.com/darkliquid/charmingui) | Render a bubbletea application to an image.Image |
| 909 | 1 | [EthanEFung/ttt](https://github.com/EthanEFung/ttt) | cli tic-tac-toe to learn charmbracelet/bubbletea |
| 910 | 1 | [BYT0723/bilichat](https://github.com/BYT0723/bilichat) | A Bilibili Live Chat TUI based on bubbletea |
| 911 | 1 | [dylanpruitt/residentsleeper](https://github.com/dylanpruitt/residentsleeper) | API client TUI; Insomnia-inspired clone created with Bubbletea |
| 912 | 1 | [vjanelle/mcp-proxy](https://github.com/vjanelle/mcp-proxy) | Go based bubbletea/lipgloss TUI proxy for MCP development |
| 913 | 1 | [MattAMonroe/ConwayTerm](https://github.com/MattAMonroe/ConwayTerm) | A TUI Conways game of life built with BubbleTea |
| 914 | 1 | [rudifa/newpro](https://github.com/rudifa/newpro) | A bubbletea cli project, generates new go or astro projects |
| 915 | 1 | [christophervistal25/repository-cleaner](https://github.com/christophervistal25/repository-cleaner) | Terminal-based GitHub repository cleaner built with Go and Bubbletea |
| 916 | 1 | [lucashancock/term-tree](https://github.com/lucashancock/term-tree) | tree explorer for the terminal, built with bubbletea and lipgloss. |
| 917 | 1 | [boomyao/clash-cli](https://github.com/boomyao/clash-cli) | TUI client for mihomo (Clash Meta) — built with Go + bubbletea |
| 918 | 1 | [visrosa/tasklineUI](https://github.com/visrosa/tasklineUI) | Bubbletea UI for taskline using kitty's text sizing protocol |
| 919 | 1 | [undefinedopcode/sid-synth](https://github.com/undefinedopcode/sid-synth) | A multi chip sid synthesizer in golang, with bubbletea interface |
| 920 | 1 | [dmorn/lit](https://github.com/dmorn/lit) | Literature crawler. Supports Elsevier's Scopus. Powered by charmbracelet/bubbletea |
| 921 | 1 | [ghthor/webtea](https://github.com/ghthor/webtea) | Run bubbletea programs as webapps using gotty & ssh apps using wish |
| 922 | 1 | [hopefulTex/slider](https://github.com/hopefulTex/slider) | A Bubbletea Bubble used to select from a range of inputs |
| 923 | 1 | [TheSquake/charm-otp](https://github.com/TheSquake/charm-otp) | TUI client for OTPClient made in go using bubbletea by charmbracelet |
| 924 | 1 | [dgroomes/bubble-tea-playground](https://github.com/dgroomes/bubble-tea-playground) | 📚 Learning and exploring the Go-based TUI framework: Bubble Tea |
| 925 | 1 | [VitexSoftware/multiflexi-tui](https://github.com/VitexSoftware/multiflexi-tui) | Terminal User Interface for MultiFlexi CLI built with Charmbracelet Bubbletea framework |
| 926 | 1 | [zKurisu/web-tree](https://github.com/zKurisu/web-tree) | A web folder in command line with UI, developing with go bubbletea |
| 927 | 1 | [itsmandrew/scoreboard-tui](https://github.com/itsmandrew/scoreboard-tui) | A scoreboard in my terminal for many sports using Go and BubbleTea |
| 928 | 1 | [Anacardo89/kanboards](https://github.com/Anacardo89/kanboards) | kanban CLI app based on Trello, built with bubbletea (elm-based) framework |
| 929 | 1 | [Nathan-ma/hubstaff-tui](https://github.com/Nathan-ma/hubstaff-tui) | Fast Hubstaff time tracking TUI for tmux floating popups — Go + Bubbletea v2 |
| 930 | 1 | [gregriff/ducky](https://github.com/gregriff/ducky) | GPT-CLI tool rewritten in Go for better, faster UX with bubbletea |
| 931 | 1 | [brucevanhorn2/exfil](https://github.com/brucevanhorn2/exfil) | Cyberpunk TUI SCP/SFTP client for Linux, built with Go, Bubbletea, and Lipgloss. |
| 932 | 1 | [jorgenho/go-tur](https://github.com/jorgenho/go-tur) | bubbletea app som lar meg se bussavanger til og fra jobb i terminalen |
| 933 | 1 | [PPabloMunoz/go-do](https://github.com/PPabloMunoz/go-do) | Go-Do is a very basic todo list using golang and bubbletea framework |
| 934 | 1 | [justinh-rahb/moontop](https://github.com/justinh-rahb/moontop) | A terminal UI for Moonraker / Klipper 3D printers, written in Go using Bubbletea. |
| 935 | 1 | [Trao95/TUI-To-do](https://github.com/Trao95/TUI-To-do) | A To Do list/tracker built for the terminal in GO, using BubbleTea. |
| 936 | 1 | [0xspector/kuroanime](https://github.com/0xspector/kuroanime) | anime in your terminal — bubbletea TUI, AniList search, MP4 downloads, Discord Rich Presence |
| 937 | 1 | [lnardon/GoT21](https://github.com/lnardon/GoT21) | A game of 21 (Blackjack) on the terminal using Golang, Bubbletea and Lipgloss. |
| 938 | 1 | [sarkarshuvojit/firebase-claims-explorer](https://github.com/sarkarshuvojit/firebase-claims-explorer) | A TUI Tool built with BubbleTea to explore custom claims attached to firebase users. |
| 939 | 1 | [n1h41/apk_builder_v3](https://github.com/n1h41/apk_builder_v3) | Commandline application to automate flutter apk building and sharing, using golang and bubbletea framework |
| 940 | 1 | [thousandflowers/qr-multi-imgs](https://github.com/thousandflowers/qr-multi-imgs) | Go BubbleTea TUI - scan qr codes from images folder, organize and export results |
| 941 | 1 | [TheComputerM/lazycph](https://github.com/TheComputerM/lazycph) | Competitive Programming Helper in your terminal |
| 942 | 1 | [CDamianS/cli-rubiks-timer](https://github.com/CDamianS/cli-rubiks-timer) | A CLI rubiks cube timer made with the Bubbletea Framework in Go (WIP) |
| 943 | 1 | [sidhyaashu/xnote-go-cli](https://github.com/sidhyaashu/xnote-go-cli) | Xnote is a fast, minimal, and elegant terminal-based note-taking application built with Go and Bubble Tea |
| 944 | 1 | [MihailKashintsev/meshdiag](https://github.com/MihailKashintsev/meshdiag) | Live terminal dashboard for your local network — concurrent ping sweep, port scan, Bubbletea TUI |
| 945 | 1 | [specCon18/bubblewand](https://github.com/specCon18/bubblewand) | A tool for generating GO code repoisitories that leverage Cobra + Viper + Log + Bubbletea + Nix |
| 946 | 1 | [altugbakan/card-logger](https://github.com/altugbakan/card-logger) | A console application for managing collectible cards |
| 947 | 1 | [mohamedbeat/togo](https://github.com/mohamedbeat/togo) | Togo is a modern, interactive terminal-based todo list application built with Bubble Tea in Go. It features a beautiful TUI, modal dialogs, keyboard navigation, and persistent storage. |
| 948 | 1 | [lingdongomg/gossh](https://github.com/lingdongomg/gossh) | A TUI (Terminal User Interface) SSH connection manager built with Go and Bubbletea. 基于 Go 和 Bubbletea构建的 TUI（终端用户界面）SSH 连接管理器。 |
| 949 | 1 | [Mohammad-Alipour/Gonsole](https://github.com/Mohammad-Alipour/Gonsole) | A fast and hackable terminal code editor built with Go, powered by Bubbletea and Chroma. |
| 950 | 1 | [tunaman/scaleway-tui](https://github.com/tunaman/scaleway-tui) | TUI for managing Scaleway resources (Object Storage, Kubernetes, Billing) — built with Bubbletea and the Dracula theme. |
| 951 | 1 | [OppaiHacker/esp-flasher](https://github.com/OppaiHacker/esp-flasher) | ⚡ Lightning-fast TUI flasher for ESP32, ESP8266, and Raspberry Pi. Built with Go and Bubbletea. |
| 952 | 1 | [coah80/coahgpt](https://github.com/coah80/coahgpt) | Self-hosted AI chat with Go backend, SvelteKit web UI, and bubbletea CLI — powered by Ollama |
| 953 | 1 | [JaredReisinger/committed](https://github.com/JaredReisinger/committed) | A bubbletea-powered text UI that integrates as a proper `commit-msg` hook for conventional commits. |
| 954 | 1 | [bntrtm/structly](https://github.com/bntrtm/structly) | Powerful bubbletea component leveraging reflection to expose struct types as menus for input by CLI users. |
| 955 | 1 | [Desgue/Tasker-CLI](https://github.com/Desgue/Tasker-CLI) | A CLI Terminal User Interface application to manage your tasks in a Kanban view. Built with BubbleTea. |
| 956 | 1 | [steelthedev/gin-mvc-cli](https://github.com/steelthedev/gin-mvc-cli) | A project created using bubbletea  to enhance code structures for beginners using gin as their API framework |
| 957 | 1 | [heliostatic/tui-clock](https://github.com/heliostatic/tui-clock) | A terminal-based world clock for tracking colleague availability across multiple timezones. Built with Go and Bubbletea. |
| 958 | 1 | [MagmaBlock/netmon](https://github.com/MagmaBlock/netmon) | 跨平台网络监测 TUI —— Go + bubbletea，单二进制零依赖 |
| 959 | 1 | [Xitonight/xidots-cli](https://github.com/Xitonight/xidots-cli) | TUI app to manage my dotfiles. Written in Go with the fantastic BubbleTea library from Charm™ |
| 960 | 1 | [Tahsin005/termnote](https://github.com/Tahsin005/termnote) | A note taking application, directly in your terminal! |
| 961 | 1 | [otaviosoaresp/dbtui](https://github.com/otaviosoaresp/dbtui) | Terminal database client with FK navigation, visual mode, row operations, and vim keybindings. Built with Go + BubbleTea. |
| 962 | 1 | [Birduo/sizzlestar](https://github.com/Birduo/sizzlestar) | A simple command line idle game written in Go using Bubbletea! Make a successful shrimp family business! |
| 963 | 1 | [apainintheneck/wpedia](https://github.com/apainintheneck/wpedia) | a simple TUI Wikipedia client |
| 964 | 1 | [T117m/MusicCatalog](https://github.com/T117m/MusicCatalog) | Курсач (3 семестр) |
| 965 | 1 | [Ruben9922/reversi](https://github.com/Ruben9922/reversi) | Command-line version of the classic Reversi / Othello game. |
| 966 | 1 | [allisonhere/ripple](https://github.com/allisonhere/ripple) | A keyboard-first, soft-wrapping multi-line text editor component for [Bubble Tea](https://github.com/charmbracelet/bubbletea). |
| 967 | 1 | [Thiti-Dev/ksc](https://github.com/Thiti-Dev/ksc) | Personal TUI for managing known commands and also make your life easier when you forget all those long args tail of certain commands |
| 968 | 1 | [sebogh/promtui](https://github.com/sebogh/promtui) | tailing Prometheus endpoints |
| 969 | 1 | [kevinliao852/dbterm](https://github.com/kevinliao852/dbterm) | A tui for db |
| 970 | 1 | [zachbonham/go-halo-cli](https://github.com/zachbonham/go-halo-cli) | A terminal user interface (TUI) that provides a developer achievement experience, inspired by Halo, powered by github.com/charmbracelet/bubbletea. |
| 971 | 1 | [bevis-hp/glyphfall](https://github.com/bevis-hp/glyphfall) | A bouncy place to pipe to! A terminal physics toy built with Go and the brilliant TUI framework from Charm - BubbleTea |
| 972 | 1 | [Bpazg97/ansii-terminal-generator](https://github.com/Bpazg97/ansii-terminal-generator) | A terminal-based ANSI art editor written in Go with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss). Draw anime characters, game sprites, or anything you like using Unicode block characters and 256-color ANSI, then set it as your terminal splash screen. |
| 973 | 1 | [davidmahbubi/trecli](https://github.com/davidmahbubi/trecli) | Manage your trello cards without ever leaving your terminal! |
| 974 | 1 | [FKouhai/urban-tui](https://github.com/FKouhai/urban-tui) |  |
| 975 | 1 | [tylerolson/tictacgo](https://github.com/tylerolson/tictacgo) | A terminal based multiplayer tic-tac-toe game. |
| 976 | 1 | [alex27riva/ytsearch](https://github.com/alex27riva/ytsearch) | Tiny TUI to search YouTube and play results with mpv. |
| 977 | 1 | [diegorezm/wallpapercl](https://github.com/diegorezm/wallpapercl) |  |
| 978 | 1 | [nbr23/soma](https://github.com/nbr23/soma) | Golang SomaFM tuner |
| 979 | 1 | [itsrainingmani/tailnode](https://github.com/itsrainingmani/tailnode) | TUI to manage Tailscale/Mullvad Exit Nodes |
| 980 | 1 | [hamadrehman/dawn-news-cli](https://github.com/hamadrehman/dawn-news-cli) | Read Dawn News From The Terminal |
| 981 | 1 | [N-Ignacio-Bouffanais/golang_tui](https://github.com/N-Ignacio-Bouffanais/golang_tui) | Terminal app written in GO that allows automating processes via ssh connection to remote servers. |
| 982 | 1 | [Julia-Marcal/Valorant-cmd](https://github.com/Julia-Marcal/Valorant-cmd) | BubbleTea-ValorantCMD is a simple tool that lets Valorant players check their stats straight from their computer's terminal. It combines the BubbleTea library for a user-friendly interface with the Unofficial Valorant API to fetch game data. Use it to stay updated on your performance without leaving the command line. |
| 983 | 1 | [phin-tech/herdr-phin-util](https://github.com/phin-tech/herdr-phin-util) | Personal Herdr Utils |
| 984 | 1 | [MaximilianSoerenPollak/GoChron](https://github.com/MaximilianSoerenPollak/GoChron) | Used to be a fork of 'zeit', but has since diverged too greatly. New UI with BubbleTea incoming, as well as new features. |
| 985 | 1 | [ayoub3bidi/sun](https://github.com/ayoub3bidi/sun) | Generate production-ready backends and documentation sites from a single binary. |
| 986 | 1 | [zigzter/league-predictions](https://github.com/zigzter/league-predictions) | Automate Twitch Predictions using the LoL and Twitch APIs |
| 987 | 1 | [stepan41k/p-manager](https://github.com/stepan41k/p-manager) | TUI Password Manager with Cloud Storage |
| 988 | 1 | [YuriBrunetto/go-repositories](https://github.com/YuriBrunetto/go-repositories) | CLI application designed to seamlessly retrieve a user's repositories from GitHub |
| 989 | 1 | [alex27riva/shto](https://github.com/alex27riva/shto) | Connect faster to SSH hosts |
| 990 | 1 | [kypkk/screensaverX](https://github.com/kypkk/screensaverX) | A terminal screensaver that's intentionally hard to dismiss. Built with Bubble Tea. |
| 991 | 1 | [elentok/gx](https://github.com/elentok/gx) | gx is a TUI for git worktree management, staging, history inspection, and other everyday repo workflows. |
| 992 | 1 | [abilun/keybon](https://github.com/abilun/keybon) | CLI app to practice typing |
| 993 | 1 | [imran-vz/gosqlit](https://github.com/imran-vz/gosqlit) | A TUI for SQL databases. Written in Golang. Currently supporting postgres |
| 994 | 1 | [importre/geeknews](https://github.com/importre/geeknews) | 긱뉴스를 터미널에서 :shipit: |
| 995 | 1 | [ZhehaoTetsuhiro/Codex-Manager](https://github.com/ZhehaoTetsuhiro/Codex-Manager) | Codex CLI 配置可视化管理器 · 现代 TUI · 模型/供应商/API Key/档案 · Go + bubbletea |
| 996 | 1 | [NaNomicon/quotta](https://github.com/NaNomicon/quotta) | Go CLI quota monitoring tool with Bubble Tea TUI, SQLite, multi-provider adapter architecture |
| 997 | 1 | [NoOPeEKS/kanbancli](https://github.com/NoOPeEKS/kanbancli) | A CLI application to manage projects and tasks with the commonly used Kanban methodology |
| 998 | 1 | [ramirezDg/lampp-tui](https://github.com/ramirezDg/lampp-tui) | A keyboard-driven terminal dashboard for managing and monitoring XAMPP/LAMPP services in real time. |
| 999 | 1 | [thewizardshell/octohook](https://github.com/thewizardshell/octohook) |  Git hooks made fast and simple. Zero dependencies, intelligent caching. |
| 1000 | 1 | [manzil-infinity180/aflock-tui](https://github.com/manzil-infinity180/aflock-tui) | Terminal UI for inspecting aflock sessions — browse sessions, decode DSSE attestations & JWTs, replay Claude Code sessions against policies, live-watch agent sessions. Built with   Bubbletea |
