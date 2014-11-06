Project templates
=================

This directory contains the project templates which are returned by the
**gondola new** command.

The *_common* directory contains files which are included in any template, while
*_appengine* holds files which are only included when the user wants a Google
App Engine-enabled template. Each one of the remaining directories contain a single
template.

A template is built using the following logic.

  * Add all files from _common.
  * If -gae is enabled, add files from _appengine (overwriting files from _common).
  * Add all the files from template itself (overwrites all the previous ones).
