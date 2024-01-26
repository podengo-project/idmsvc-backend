#!/usr/bin/python3
"""Generate app-secret for ephemeral environment
"""
import argparse
import os
import secrets

TEMPLATE = """\
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
stringData:
  app_secret: {secret}
"""

parser = argparse.ArgumentParser()
parser.add_argument("file", type=argparse.FileType("w"))


def main():
    args = parser.parse_args()
    print(f"Writing app-secret to {args.file.name}")
    content = TEMPLATE.format(secret=secrets.token_urlsafe())
    with args.file as f:
        f.write(content)


if __name__ == "__main__":
    main()
