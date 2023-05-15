# install Python tools in a virtual environment

PYTHON_VENV := .venv
BONFIRE := $(PYTHON_VENV)/bin/bonfire
PRE_COMMIT := $(PYTHON_VENV)/bin/pre-commit
JSON2YAML := $(PYTHON_VENV)/bin/json2yaml

$(PYTHON_VENV):
	python3 -m venv $(PYTHON_VENV)
	$(PYTHON_VENV)/bin/pip install -U pip setuptools

$(BONFIRE) $(PRE_COMMIT) $(JSON2YAML): $(PYTHON_VENV)
	$(PYTHON_VENV)/bin/pip3 install -r requirements-dev.txt
	touch $(BONFIRE) $(PRE_COMMIT)

.PHONY: install-python-tools
install-python-tools:
	$(MAKE) $(BONFIRE)
	$(MAKE) $(PRE_COMMIT)
	$(MAKE) $(JSON2YAML)
