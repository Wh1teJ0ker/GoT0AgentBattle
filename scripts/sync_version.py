#!/usr/bin/env python3
"""同步仓库版本到前端 package.json 等元数据文件。"""

from __future__ import annotations

import json
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
VERSION_FILE = ROOT / ".version"
FRONTEND_PACKAGE = ROOT / "frontend" / "package.json"


def read_version() -> str:
    raw = VERSION_FILE.read_text(encoding="utf-8").strip()
    if not raw:
        raise ValueError(".version 不能为空")
    if not raw.startswith("v"):
        raise ValueError(".version 必须使用 v 前缀，例如 v0.1.0")
    return raw


def sync_frontend(version: str) -> None:
    payload = json.loads(FRONTEND_PACKAGE.read_text(encoding="utf-8"))
    payload["version"] = version.removeprefix("v")
    FRONTEND_PACKAGE.write_text(
        json.dumps(payload, indent=2, ensure_ascii=False) + "\n",
        encoding="utf-8",
    )


def main() -> None:
    version = read_version()
    sync_frontend(version)
    print(f"sync version -> {version}")


if __name__ == "__main__":
    main()
