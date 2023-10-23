##
# Entrypoint for the Makefile
#
# It is composed at mk/includes.mk by including
# small make files which provides all the necessary
# rules.
#
# Some considerations:
#
# - Variables customization can be
#   stored at '.env', 'mk/private.mk' files.
# - By default the 'help' rule is executed.
##

include scripts/mk/includes.mk

# Set the default rule
.DEFAULT_GOAL := help
