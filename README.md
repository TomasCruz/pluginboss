# pluginboss
pluginboss is a small CLI app (could be a web service) calling plugins using go-plugin

## plugininfo
plugininfo contains code skeleton to be used for and by plugins. Also proto file and generated code. Ideally, this package should live in a library, possibly a separate one

## plugin-*
For demonstration purposes only, these plugins are part of the repo. However, they only depend on plugininfo package (library)

Plugins 3 && 4 are made broken to demonstrate pluginboss' being independant of plugins, usable ones remaining usable
