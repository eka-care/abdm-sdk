#!/usr/bin/env python3
"""Regenerate docs/abdm-api-reference.md from developer.eka.care.

The published docs are the source of truth for this SDK. This script enumerates every
page under the ABDM Connect section via llms.txt, parses the embedded OpenAPI blocks and
webhook payload samples, and writes a single reference file.

    python3 scripts/extract-api-docs.py

Requires pyyaml. Network access required. Re-run and diff when the docs change.
"""
import collections, datetime, json, os, pathlib, re, subprocess, sys
from concurrent.futures import ThreadPoolExecutor
from urllib.request import urlopen

INDEX = "https://developer.eka.care/llms.txt"
SECTION = "https://developer.eka.care/api-reference/user-app/abdm-connect/"
ROOT = pathlib.Path(__file__).resolve().parent.parent
OUT = ROOT / "docs" / "abdm-api-reference.md"
SKIP_SCHEMAS = {"GenericError", "SourceErr", "ModelsError"}

try:
    import yaml
except ImportError:
    sys.exit("pyyaml required: pip install pyyaml")


def get(url):
    with urlopen(url, timeout=30) as r:
        return r.read().decode("utf-8", "replace")


def implemented_paths():
    """Endpoint paths this SDK already calls, for the coverage column."""
    out = subprocess.run(["grep", "-rho", r'Path:  *"[^"]*"', "--include=*.go", str(ROOT)],
                         capture_output=True, text=True).stdout
    return {m.group(1) for m in re.finditer(r'Path:\s*"([^"]*)"', out)}


def yaml_blocks(text):
    for m in re.finditer(r"^(`{3,})yaml\s+(\w+)\s+(\S+)\s*\n(.*?)^\1\s*$", text, re.S | re.M):
        yield m.group(2).upper(), m.group(3), m.group(4)


def type_of(v):
    if "$ref" in v:
        return v["$ref"].split("/")[-1]
    t, nullable = v.get("type", ""), False
    if isinstance(t, list):
        nullable = "null" in t
        t = next((x for x in t if x != "null"), "")
    if t == "array":
        it = v.get("items") or {}
        t = "[]" + (it["$ref"].split("/")[-1] if "$ref" in it else it.get("type", ""))
    return t + ("?" if nullable else "")


def refs(node, acc):
    if isinstance(node, dict):
        if "$ref" in node:
            acc.add(node["$ref"].split("/")[-1])
        for v in node.values():
            refs(v, acc)
    elif isinstance(node, list):
        for v in node:
            refs(v, acc)


def props(schema, comps):
    if not isinstance(schema, dict):
        return []
    if "$ref" in schema:
        schema = comps.get(schema["$ref"].split("/")[-1], {})
    required = schema.get("required") or []
    return [(k + ("*" if k in required else ""), type_of(v))
            for k, v in (schema.get("properties") or {}).items()]


def parse_endpoints(fname, text):
    title = (re.search(r"^# (.+)$", text, re.M) or [None, ""])[1]
    for method, path, body in yaml_blocks(text):
        spec = yaml.safe_load(body)
        comps = ((spec.get("components") or {}).get("schemas") or {})
        op = ((spec.get("paths") or {}).get(path) or {}).get(method.lower(), {}) or {}
        params = op.get("parameters") or []
        req = (((op.get("requestBody") or {}).get("content") or {})
               .get("application/json") or {}).get("schema")
        ok = next((c for c in sorted((op.get("responses") or {}), key=str) if str(c).startswith("2")), None)
        resp = None
        if ok is not None:
            resp = ((op["responses"][ok].get("content") or {}).get("application/json") or {}).get("schema")

        used, seen = set(), set()
        refs(req, used)
        refs(resp, used)
        while used - seen:
            n = (used - seen).pop()
            seen.add(n)
            more = set()
            refs(comps.get(n, {}), more)
            used |= more

        yield dict(
            file=fname, title=title, method=method, path=path, ok=str(ok),
            headers=[p["name"] for p in params if p.get("in") == "header"],
            path_params=[p["name"] for p in params if p.get("in") == "path"],
            query=[p["name"] for p in params if p.get("in") == "query"],
            req=props(req, comps), resp=props(resp, comps),
            nested={n: props(comps[n], comps) for n in sorted(seen)
                    if n in comps and n not in SKIP_SCHEMAS})


def shape(v):
    if isinstance(v, list):
        inner = "{" + ", ".join(sorted(v[0])) + "}" if v and isinstance(v[0], dict) else "string"
        return "[]" + inner
    if isinstance(v, dict):
        return "{" + ", ".join(sorted(v)) + "}"
    return type(v).__name__.replace("str", "string").replace("int", "integer")


def parse_webhook(fname, text):
    title = (re.search(r"^# (.+)$", text, re.M) or [None, fname])[1]
    for m in re.finditer(r"```json[^\n]*\n(.*?)```", text, re.S):
        try:
            payload = json.loads(m.group(1))
        except ValueError:
            continue
        if payload.get("event"):
            return payload["event"], title, sorted((k, shape(v)) for k, v in (payload.get("data") or {}).items())
    return None


def main():
    urls = sorted(set(re.findall(re.escape(SECTION) + r"[^)\s]*\.md", get(INDEX))))
    print(f"fetching {len(urls)} pages...")
    with ThreadPoolExecutor(max_workers=8) as pool:
        pages = list(zip(urls, pool.map(get, urls)))

    endpoints, hooks = [], {}
    for url, text in pages:
        name = url[len(SECTION):-3].replace("/", "__")
        endpoints.extend(parse_endpoints(name, text))
        if name.startswith("webhooks__"):
            hook = parse_webhook(name, text)
            if hook:
                hooks.setdefault(hook[0], hook)

    impl = implemented_paths()
    by = collections.defaultdict(list)
    for e in endpoints:
        by["/".join(e["file"].split("__")[:-1]) or "root"].append(e)
    done = sum(1 for e in endpoints if e["path"] in impl)

    L = ["# ABDM Connect API reference (extracted)", "",
         f"Generated {datetime.date.today()} by `scripts/extract-api-docs.py` from all",
         f"{len(endpoints)} endpoints across {len(urls)} pages under",
         "`developer.eka.care/api-reference/user-app/abdm-connect`, enumerated via `llms.txt`.", "",
         "The published docs are the source of truth. Regenerate this file rather than",
         "hand-editing it.", "",
         "`*` = required. `?` = nullable. `✅` = already implemented in this SDK.", "",
         f"Coverage: **{done} of {len(endpoints)}** endpoints implemented.", "",
         "| Section | Endpoints | Implemented |", "|---|---:|---:|"]
    for s in sorted(by):
        L.append(f"| {s} | {len(by[s])} | {sum(1 for e in by[s] if e['path'] in impl)} |")

    for s in sorted(by):
        L += ["", f"## {s}"]
        for e in sorted(by[s], key=lambda e: (e["path"], e["method"])):
            L += ["", f"### `{e['method']} {e['path']}`" + (" ✅" if e["path"] in impl else ""), ""]
            if e["title"]:
                L += [e["title"], ""]
            meta = [f"{lbl}: " + ", ".join(f"`{x}`" for x in e[k])
                    for lbl, k in (("headers", "headers"), ("path", "path_params"), ("query", "query")) if e[k]]
            L.append("- " + " · ".join(meta + [f"success: `{e['ok']}`"]))
            for lbl, k in (("request", "req"), ("response", "resp")):
                L.append(f"- {lbl}: " + (", ".join(f"`{n}` {t}".strip() for n, t in e[k]) if e[k] else "_none_"))
            if e["nested"]:
                L += ["", "<details><summary>schemas</summary>", ""]
                for n, f in e["nested"].items():
                    L.append(f"- **{n}** — " + (", ".join(f"`{a}` {b}".strip() for a, b in f) if f else "_empty_"))
                L += ["", "</details>"]

    L += ["", "## webhooks", "",
          "All webhooks POST to the integrator's registered endpoint with header",
          "`Eka-Webhook-Signature: t=<unix>,v1=<hex>` — HMAC-SHA256 over `\"{t}.{rawBody}\"` using the",
          "subscription's signing key. Reject when `|now - t|` exceeds ~3 minutes.", "",
          "Common envelope: `service`, `event`, `event_time`, `transaction_id`, `timestamp`,",
          "`business_id`, `client_id`, `data`. Fields below are those of `data`.", "",
          "| Event | Description | `data` fields |", "|---|---|---|"]
    for ev, title, data in sorted(hooks.values()):
        L.append(f"| `{ev}` | {title} | " + (", ".join(f"`{k}`" for k, _ in data) or "_none_") + " |")
    L.append("")
    for ev, title, data in sorted(hooks.values()):
        L += [f"### `{ev}`", "", title, ""] + [f"- `{k}` {t}" for k, t in data] + [""]

    OUT.write_text("\n".join(L))
    print(f"wrote {OUT.relative_to(ROOT)}: {len(endpoints)} endpoints, {len(hooks)} webhooks, {done} implemented")


if __name__ == "__main__":
    main()
