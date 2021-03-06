#!/usr/bin/python

import os

class LinkError(Exception):
	"""Raised when forcing a symlink fails for a non-OS reason."""
	pass

def Force(source, linkName):
	"""Force forces linkName to be a symlink to source, as long as its not a dir.

		Creates a symlink from linkName to source, clobbering linkName as long as its not a directory.
	"""
	if not os.path.lexists(linkName):
		os.symlink(source, linkName)
		return

	if os.path.islink(linkName) or os.path.isfile(linkName):
		os.remove(linkName)
		os.symlink(source, linkName)
		return

	raise LinkError("Refusing to clobber " + linkName)
