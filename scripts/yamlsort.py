#!.venv/bin/python3
"""yamlfix with custom sorting for OpenAPI"""
import sys
import warnings

from yamlfix import fix_files
from yamlfix.adapters import YamlfixRepresenter
from yamlfix.model import YamlfixConfig, YamlNodeStyle

warnings.filterwarnings("ignore", category=UserWarning)

PRIORITIES = {
    # top-level block
    "openapi": -1000,
    "info": -900,
    "servers": -800,
    "paths": -700,
    # block name/title
    "name": -200,
    "title": -100,
    # additional human text
    "summary": -90,
    "description": -90,
    # schema core elements
    "type": -80,
    "required": -70,
    "additionalProperties": -60,
    "properties": -50,
    "items": -50,
    # response schema
    "schema": -50,
    # sort after all other keys
    "example": 1000,
}


def _customsort(kv):
    key = kv[0]
    return PRIORITIES.get(key, 0), key


INDENT_BASE = 4

config = YamlfixConfig(
    indent_mapping=INDENT_BASE,
    indent_offset=INDENT_BASE,
    indent_sequence=INDENT_BASE + 2,
    line_length=1024,
    sequence_style=YamlNodeStyle.BLOCK_STYLE,
)

orig_represet_mapping = YamlfixRepresenter.represent_mapping


def patched_represent_mapping(self, tag, mapping, flow_style=None):
    if hasattr(mapping, "items"):
        # sort with custom key sorting
        sorted_list = sorted(mapping.items(), key=_customsort)
        # re-use mapping, ruyaml stores comments in additional attributes.
        mapping.clear()
        mapping.update(sorted_list)
        # sort required and enum list
        for key in ("required", "enum"):
            value = mapping.get(key)
            if value and isinstance(value, list):
                value.sort()

    return orig_represet_mapping(self, tag, mapping, flow_style=flow_style)


YamlfixRepresenter.represent_mapping = patched_represent_mapping

if __name__ == "__main__":
    fix_files(sys.argv[1:], config=config)
