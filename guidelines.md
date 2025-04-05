# Guidelines

## Workflow

- Define variables such as paths, urls and similar at the top of the file.
- Define functions for each step of workflow execution, which shall contain logic for that set of tasks and config definitions.
- Keep each step separate. Think of them as something that may fail and should not affect others.
- If your workflow requires set of steps to be synchronised(dependant), then define that logic as "wrapper", outside of steps definition.
- In similar fashion, you can define asynchronous execution using task groups and goroutines. Currently, must be implemented by you.

## Task

- Think of task as unit of work.
- If your task can be reused by other tasks, consider moving it to TaskHelpers.
- If your task seems to be too broad/big, consider splitting into subset of tasks and combining them on workflow level.
- Each task must have Validate and Run functions. The Update function is neither defined nor implemented as of yet.

## TaskHelper

- Reusable accross Tasks.

## ValidationHelper

- Reusable accross Tasks(Validate function).

## Function

- Neither Task nor Workflow bound.

