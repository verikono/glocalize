#!/usr/bin/env python3
"""Import a filesystem skill tree into the riot registry.

A skill on disk may carry sibling markdown files under `assets/` and
`references/`. riot has no filesystem import and no field for sibling files, so
each one becomes a Reference named `<skill>-<stem>` and the skill body's relative
path pointers are rewritten to name those references. The skill then stays small
and the detail is fetched with `get-reference` only when a task needs it.

Usage:
    python3 import-skill-tree.py aws-database aws-iam aws-containers
    python3 import-skill-tree.py --apply aws-database

Without --apply it prints the plan and writes nothing. Re-running is safe: the
previous run's references and skill records are cleared first.
"""

import json
import os
import re
import subprocess
import sys
import urllib.request

GRAPH = os.environ.get("RIOT_GRAPH", "http://localhost:6661/graph")
DB = os.environ.get("RIOT_DB", os.path.expanduser("~/.config/riot/data.sqlite"))
SKILLS_ROOT = os.path.expanduser("~/.agents/skills")
PROJECT = "all"

# Descriptions over the 128-word limit are trimmed to their first sentences,
# which is where the what-it-does and when-to-load-it live.
DESCRIPTION_WORD_LIMIT = 128


def gql(query, variables=None):
    payload = json.dumps({"query": query, "variables": variables or {}}).encode()
    req = urllib.request.Request(
        GRAPH, data=payload, headers={"Content-Type": "application/json"}
    )
    with urllib.request.urlopen(req) as resp:
        body = json.loads(resp.read())
    if "errors" in body:
        raise RuntimeError(json.dumps(body["errors"], indent=2))
    return body["data"]


def parse_frontmatter(text):
    m = re.match(r"^---\n(.*?)\n---\n", text, re.S)
    if not m:
        raise RuntimeError("no frontmatter")
    fm = m.group(1)
    body = text[m.end():]
    desc = re.search(r"^description:\s*(.*?)(?=^version:|\Z)", fm, re.S | re.M).group(1)
    desc = re.sub(r">-?", "", desc)
    return " ".join(desc.split()), body


def trim_description(desc):
    if len(desc.split()) <= DESCRIPTION_WORD_LIMIT:
        return desc
    out = ""
    for sentence in re.split(r"(?<=\.)\s+", desc):
        candidate = (out + " " + sentence).strip()
        if len(candidate.split()) > DESCRIPTION_WORD_LIMIT:
            break
        out = candidate
    return out or desc


def first_heading(text):
    m = re.search(r"^#\s+(.+)$", text, re.M)
    return m.group(1).strip() if m else ""


def bullet(text, label):
    """The value of a `- **Label**: value` bullet, or "" when absent."""
    m = re.search(rf"^-\s+\*\*{re.escape(label)}\*\*:\s*(.+)$", text, re.M)
    return m.group(1).strip() if m else ""


def first_paragraph(text):
    """The first prose paragraph after the H1."""
    out, started = [], False
    for line in text.split("\n"):
        s = line.strip()
        if s.startswith("#"):
            if started:
                break
            continue
        if not s:
            if started:
                break
            continue
        started = True
        out.append(s)
    return " ".join(out)


def describe(text):
    """A one-line description of a sibling file.

    The knowledge cards are structured bullet lists rather than prose, so the
    data model and the workload they suit say more than the first paragraph
    would. A reference file is prose, so its opening paragraph is the summary.
    """
    heading = first_heading(text)
    model = bullet(text, "Data model")
    best = bullet(text, "Best for")
    if model or best:
        parts = [p for p in (model, best) if p]
        return f"{heading} — " + "; ".join(parts) if heading else "; ".join(parts)
    para = first_paragraph(text)
    return f"{heading} — {para}" if heading else para


def collect_files(skill):
    root = os.path.join(SKILLS_ROOT, skill)
    found = []
    for sub in ("assets", "references"):
        d = os.path.join(root, sub)
        if not os.path.isdir(d):
            continue
        for name in sorted(os.listdir(d)):
            if name.endswith(".md"):
                found.append((sub, name, os.path.join(d, name)))
    return found


def rewrite_pointers(body, skill, refs):
    """Point every relative path at the reference that now holds that file."""
    out = body
    for r in refs:
        stem = r["file_stem"]
        for pat in (f"`assets/{stem}.md`", f"`references/{stem}.md`",
                    f"[{stem}.md](assets/{stem}.md)",
                    f"[{stem}.md](references/{stem}.md)",
                    f"assets/{stem}.md", f"references/{stem}.md"):
            out = out.replace(pat, f"`{r['name']}`")

    leftover = re.findall(r"(?:assets|references)/[a-z0-9-]+\.md", out)
    if leftover:
        raise RuntimeError(f"{skill}: unrewritten pointers: {sorted(set(leftover))}")

    # a rewritten markdown link leaves the label and the target both naming the
    # reference: [`select`](`aws-database-select`) -> `aws-database-select`
    out = re.sub(r"\[`([^`]+)`\]\(`([^`]+)`\)", r"`\2`", out)

    # the generic placeholders point at the directory, not a file, so no
    # per-file rewrite reaches them. They become the fetch instruction.
    out = out.replace(
        "read `references/{sub-skill-id}.md` and follow its procedure",
        "call `get-reference` with that sub-skill's reference name and follow its procedure",
    )
    out = out.replace(
        "answer from the knowledge cards in `assets/`",
        "answer from the knowledge cards, which are the references listed below",
    )
    out = out.replace(
        "Read `assets/{filename}` for the relevant service(s)",
        "call `get-reference` with the relevant service's reference name",
    )
    out = out.replace(
        "fetch each reference via `aws___retrieve_skill` with the `file` parameter "
        "(e.g. `file=\"references/ecs.md\"` or `file=\"references/action-logs.md\"`)",
        "fetch each reference with `get-reference` "
        "(e.g. `aws-containers-ecs` or `aws-containers-action-logs`)",
    )
    out = out.replace(
        "Do NOT `file_read` these paths locally — they do not exist on disk.",
        "Do NOT `file_read` these paths locally — they do not exist on disk; the registry holds them.",
    )

    # the guardrail section describes how the skill's own files were resolved
    # before the move into the registry; in riot there is exactly one answer.
    out = re.sub(
        r"## Guardrail — where this skill's own files live \(MCP vs local install\).*?(?=\n## )",
        "## Guardrail — where this skill's own files live\n\n"
        "This skill's reference files are held in the riot registry, not on disk. "
        "Fetch one with `get-reference` and its name (e.g. `aws-containers-ecs`); "
        "the names are listed under References at the end of this skill. "
        "Do NOT `file_read` a path — no such file exists.\n\n",
        out,
        flags=re.S,
    )
    return out


def build_plan(skills):
    plan = []
    for skill in skills:
        with open(os.path.join(SKILLS_ROOT, skill, "SKILL.md")) as f:
            raw = f.read()
        description, body = parse_frontmatter(raw)
        description = trim_description(description)

        refs = []
        for sub, name, path in collect_files(skill):
            file_stem = name[:-3]
            # a file already carrying the skill's prefix would otherwise read
            # aws-iam-aws-iam-policy-generation
            stem = file_stem
            if stem.startswith(skill + "-"):
                stem = stem[len(skill) + 1:]
            with open(path) as f:
                content = f.read()
            desc = describe(content)
            if len(desc) > 300:
                desc = desc[:297].rsplit(" ", 1)[0] + "…"
            refs.append({"name": f"{skill}-{stem}", "description": desc,
                         "content": content, "source": f"{sub}/{name}",
                         "file_stem": file_stem})

        plan.append({"skill": skill, "description": description,
                     "body": rewrite_pointers(body, skill, refs), "refs": refs})
    return plan


def main():
    args = [a for a in sys.argv[1:] if not a.startswith("--")]
    apply = "--apply" in sys.argv
    if not args:
        print(__doc__)
        return 1

    plan = build_plan(args)
    total = sum(len(p["refs"]) for p in plan)
    print(f"{'APPLYING' if apply else 'DRY RUN'}: {len(plan)} skills, {total} references")
    for p in plan:
        print(f"  {p['skill']}: body {len(p['body'])} bytes, {len(p['refs'])} refs")
        for r in p["refs"]:
            print(f"    - {r['name']}  ({r['source']}, {len(r['content'])} bytes)")

    if not apply:
        print("\nre-run with --apply to write")
        return 0

    # SaveReference rejects an existing name and SaveSkill without an id tries to
    # create a duplicate, so a re-run must clear what a previous run wrote. The
    # skill's reference set is replaced wholesale below, so the join rows go too.
    # There is no DropReference mutation, so the clear goes through sqlite.
    names = "','".join(r["name"] for p in plan for r in p["refs"])
    subprocess.run(
        ["sqlite3", DB,
         f"delete from SkillReference where referenceName in ('{names}');"
         f"delete from Reference where name in ('{names}');"],
        check=True,
    )
    for p in plan:
        prior = gql("query($n: String!) { Skill(name: $n) { id } }", {"n": p["skill"]})
        p["id"] = (prior.get("Skill") or {}).get("id") or ""

    for p in plan:
        for r in p["refs"]:
            gql(
                "mutation($input: ReferenceInput!) { SaveReference(input: $input) { name } }",
                {"input": {"name": r["name"], "description": r["description"],
                           "accessType": "static", "resourceType": "literature",
                           "content": r["content"]}},
            )
        gql(
            "mutation($input: SkillInput!) { SaveSkill(input: $input) { name } }",
            {"input": {"id": p["id"], "name": p["skill"],
                       "description": p["description"],
                       "url": f"skill://{PROJECT}/{p['skill']}", "content": p["body"],
                       "references": [r["name"] for r in p["refs"]],
                       "projects": [PROJECT]}},
        )
        print(f"wrote {p['skill']} + {len(p['refs'])} references")
    return 0


if __name__ == "__main__":
    sys.exit(main())
