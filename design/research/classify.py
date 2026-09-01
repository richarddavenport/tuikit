#!/usr/bin/env python3
"""Bucket bubbletea-ecosystem.json into libraries vs applications, then by kind.

Two stages:

  1. Is this a thing you BUILD WITH (library/component/framework), or a thing you
     RUN (an app)? The discriminator is a "for-builders" phrase — "component for
     bubbletea", "library for TUIs" — not the mere presence of the word
     "framework", which nearly every app description contains as "built with the
     Bubble Tea framework".
  2. Subcategorise each side: libraries by what they give you, apps by domain.

Heuristics over repo name + description, so expect noise; the point is a rough
map of the ecosystem, not a taxonomy. Re-run after refreshing the JSON:

    python3 classify.py
"""
import json, re, pathlib, collections

HERE = pathlib.Path(__file__).parent
repos = json.loads((HERE / "bubbletea-ecosystem.json").read_text())

def rx(*words):
    return re.compile("|".join(words), re.I)

# Phrases that mean "this is for people building TUIs", not "this is a TUI".
FOR_BUILDERS = rx(
    r"\bfor (charm'?s? )?(bubble ?tea|bubbles|lipgloss|glamour|wish)\b",
    r"\bbubble ?tea[- ](component|components|library|libraries|widget|widgets|bubble|bubbles|compatible|programs|apps|applications|models|adapter|helpers?)\b",
    r"\b(component|components|widget|widgets|bubble|bubbles)\b.*\b(bubble ?tea|lipgloss|tui|terminal ui)\b",
    r"\b(library|libraries|framework|toolkit|sdk|helpers?|utilities|primitives|adapter|wrapper|extensions?)\b.*\b(for|to) (build|building|creat|writ|mak)",
    r"\b(library|framework|toolkit|helpers?|utilities|components?)\b.*\bfor (go|golang)? ?(clis?|tuis?|terminal uis?)\b",
    r"\breusable\b.*\bcomponents?\b",
    r"\bcollection of (bubbles|components|widgets)\b",
    r"\b(tui|terminal ui)[- ](framework|toolkit|library|components?|kit)\b",
    r"\bapp framework\b", r"\bcomponent system\b", r"\blayout manager\b",
    r"\bui components\b", r"\bcomponent library\b",
    # A bare "component"/"widget"/"bubble" (the Charm word for a component) is
    # itself the tell — "Marquee bubble component", "Heatmap component".
    r"\bcomponents?\b", r"\bwidgets?\b", r"\bbubbles\b(?!\s*tea)",
    r"\bfor (the )?(golang |go )?bubble ?tea\b",
    r"\byour (bubble ?tea|tui|terminal) (app|application)s?\b",
    r"\band your tuis\b",
)
# Strong "this is an application" signal that overrides a weak library match.
IS_APP = rx(r"\b(built|made|written|created|developed|powered|coded)\s+(with|using|in|by)\b",
            r"\bmy (first|own)\b", r"\blearning\b", r"\bto learn\b")

LIB_KINDS = [
    ("notification/overlay", rx(r"\boverlay", r"\bmodal\b", r"\bdialog\b", r"\bpopup",
                                r"\bnotification", r"\btoast", r"\balert")),
    ("input/editor",         rx(r"\binput\b", r"\btext ?area", r"\beditor\b", r"\bline editor",
                                r"\bprompt", r"\bform\b", r"\bautocomplete", r"\bvim\b", r"\bslider")),
    ("picker/list/table",    rx(r"picker", r"\blist\b", r"\btable\b", r"\bgrid\b", r"\btree\b",
                                r"\bmenu\b", r"\bmenubar", r"\bselect", r"\btabs?\b", r"\bcarousel")),
    ("chart/graph",          rx(r"\bterminal charts?\b", r"\bchart", r"\bgraph", r"\bplot", r"\bheatmap", r"\bsparkline",
                                r"\bdiagram", r"\bansi art", r"\banimation")),
    ("layout/navigation",    rx(r"\blayout", r"\bnav(igation|stack)?\b", r"\brout(er|ing)\b",
                                r"\bbreadcrumb", r"\bviewport", r"\bpager\b", r"\bmulti-?tab")),
    ("styling/theming",      rx(r"\btheme", r"\bstyl(e|es|ing)\b", r"\bcolou?r", r"\blipgloss\b")),
    ("testing/devtools",     rx(r"\btest", r"\bstorybook", r"\bdebug", r"\bdev tool", r"\bharness")),
    ("framework/app-shell",  rx(r"\bframework", r"\btoolkit", r"\bapp framework", r"\bcomponent system",
                                r"\bsdk\b", r"\barchitecture")),
]
APP_DOMAINS = [
    ("ai/llm agents",     rx(r"\bllm\b", r"\bai\b", r"\bclaude\b", r"\bopenai\b", r"\bgpt\b", r"\bgemini",
                             r"\bollama", r"\banthropic", r"\bagents?\b", r"\bmcp\b", r"\bcopilot",
                             r"\bcodex\b", r"\bopencode\b", r"\bprompt engineering", r"\bchatbot")),
    ("git/github/forge",  rx(r"\bgit\b", r"\bgithub\b", r"\bgitlab", r"\bgitea", r"\bforgejo",
                             r"\bcommit", r"\bpull request", r"\bpr\b", r"\bjujutsu\b", r"\bjj\b",
                             r"\bdiff\b", r"\bmerge\b", r"\brepo(sitor(y|ies))?\b", r"\bworktree")),
    ("kubernetes/cloud",  rx(r"\bkubernetes\b", r"\bk8s\b", r"\bkubectl", r"\bdocker", r"\bcontainer",
                             r"\baws\b", r"\bs3\b", r"\bazure", r"\bgcp\b", r"\bterraform",
                             r"\bcloud\b", r"\bdeploy", r"\bhelm\b", r"\bnomad\b", r"\bansible")),
    ("db/data",           rx(r"\bdatabase", r"\bsql\b", r"\bpostgres", r"\bmysql", r"\bsqlite",
                             r"\bredis\b", r"\bmongo", r"\bkafka", r"\bclickhouse", r"\bdynamodb",
                             r"\bcsv\b", r"\bjson\b", r"\bparquet")),
    ("http/api/network",  rx(r"\bhttp\b", r"\bapi client", r"\bapi testing", r"\brest\b", r"\bgraphql",
                             r"\bpostman", r"\bcurl\b", r"\bnetwork", r"\bdns\b", r"\bssh\b",
                             r"\bport\b", r"\bpacket", r"\bwireguard", r"\bvpn\b", r"\bproxy\b",
                             r"\bwebsocket", r"\bgrpc\b")),
    ("system/monitoring", rx(r"\bsystem monitor", r"\bprocess", r"\bcpu\b", r"\bmemory\b", r"\bdisk\b",
                             r"\btop\b", r"\bhtop\b", r"\bresource", r"\bmetrics", r"\bobservab",
                             r"\blogs?\b", r"\bjournal", r"\bsystemd", r"\bpackage manager",
                             r"\bapt\b", r"\bpacman\b", r"\bhomebrew", r"\bbrew\b", r"\baur\b")),
    ("files/dotfiles",    rx(r"\bfile manager", r"\bfile explorer", r"\bfile browser", r"\bfiles\b",
                             r"\bdirector(y|ies)\b", r"\bdotfiles", r"\bchezmoi", r"\bbackup",
                             r"\bsync\b", r"\bftp\b", r"\bdisk usage")),
    ("productivity/notes",rx(r"\btodo\b", r"\btask", r"\bkanban", r"\bnote", r"\bjournal", r"\bhabit",
                             r"\bpomodoro", r"\btimer\b", r"\bcalendar", r"\btime track", r"\bflash ?card",
                             r"\bbookmark", r"\bpassword", r"\bsecrets?\b", r"\bvault\b")),
    ("media/music",       rx(r"\bmusic", r"\bspotify", r"\bplayer\b", r"\bmpd\b", r"\bmpv\b", r"\baudio",
                             r"\bpodcast", r"\bradio\b", r"\byoutube", r"\bvideo", r"\bmovie", r"\banime",
                             r"\bimage", r"\bphoto", r"\bebook", r"\bepub", r"\bpdf\b", r"\breader\b",
                             r"\brss\b", r"\bfeed\b", r"\bnews\b")),
    ("comms/social",      rx(r"\bchat\b", r"\bemail\b", r"\bmail\b", r"\bslack\b", r"\bdiscord",
                             r"\bmatrix\b", r"\btelegram", r"\bmastodon", r"\birc\b", r"\bsms\b",
                             r"\bhacker ?news", r"\breddit", r"\bforum")),
    ("finance/crypto",    rx(r"\bstock", r"\bcrypto", r"\bbitcoin", r"\btrading", r"\bportfolio",
                             r"\bmarket", r"\bfinance", r"\bledger", r"\bhledger", r"\bbudget",
                             r"\bexpense", r"\binvoice", r"\bticker")),
    ("games/toys",        rx(r"\bgame\b", r"\bgames\b", r"\btetris", r"\bsnake\b", r"\bsudoku", r"\bchess",
                             r"\bwordle", r"\b2048\b", r"\bminesweeper", r"\brogue", r"\bsolitaire",
                             r"\bpoker\b", r"\bpokemon", r"\bpokedex", r"\btyping (speed|test)",
                             r"\bpuzzle", r"\btoy\b", r"\bfun\b", r"\bascii art")),
    ("learning/template", rx(r"\btemplate", r"\bboilerplate", r"\bstarter", r"\bscaffold", r"\bexample",
                             r"\bdemo\b", r"\btutorial", r"\blearning\b", r"\bto learn\b", r"\bplayground",
                             r"\bpractice\b", r"\bexperiment", r"\bsandbox\b", r"\bmy first\b", r"\bwip\b")),
]

def first(rules, hay, default):
    for name, pat in rules:
        if pat.search(hay):
            return name
    return default

libs, apps = [], []
for r in repos:
    hay = f"{r['fullName'].split('/')[-1]} {r.get('description') or ''}"
    is_lib = bool(FOR_BUILDERS.search(hay)) and not (
        IS_APP.search(hay) and not re.search(r"\bcomponent|\bwidget|\bfor bubble ?tea", hay, re.I))
    (libs if is_lib else apps).append((r, hay))

buckets = collections.defaultdict(list)
for r, hay in libs:
    buckets["lib:" + first(LIB_KINDS, hay, "general/misc")].append(r)
for r, hay in apps:
    d = "no description" if not (r.get("description") or "").strip() else first(APP_DOMAINS, hay, "other apps")
    buckets["app:" + d].append(r)

for k in sorted(buckets, key=lambda k: (not k.startswith("lib:"), -len(buckets[k]))):
    print(f"{len(buckets[k]):5d}  {k}")
print(f"{sum(len(v) for k, v in buckets.items() if k.startswith('lib:')):5d}  TOTAL libraries")
print(f"{sum(len(v) for k, v in buckets.items() if k.startswith('app:')):5d}  TOTAL apps")
(HERE / "buckets.json").write_text(json.dumps(buckets, indent=1))
