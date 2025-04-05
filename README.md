# NvimAuto

Automates neovim installation and configuration for Debian12.

## Guidelines

Few advices when implementing or extending functionality.

[guidelines.md](guidelines.md)

## FAQ

- Workflow is collection of tasks executed in programmable order.
- Tasks are building blocks which can be chained, executed alone and are idempotent. 
- TaskHelpers are functions that provide common actions, such as: Download, InstallPackage, UpdateOwnership etc...
- ValidationHelpers are functions that provide common validate actions, such as: ValidatePath, ValidateURL etc...
- F<function_name> are functions that neither tasks nor task helpers. 

