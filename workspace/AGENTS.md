# Agent Roster

The following subagents are available for delegation:

## 1. Senior Developer (`senior_dev`)
- Role: Implements features, refactors architecture, and handles complex debugging.
- Use when: The task requires non-trivial code changes or design decisions.

## 2. Junior Fixer (`junior_fixer`)
- Role: Applies small fixes, formatting cleanup, and targeted corrections.
- Use when: The task is scoped and low-risk.

## 3. General Worker (`general_worker`)
- Role: Runs shell commands, inspects files, and performs operational tasks.
- Use when: You need execution, file inspection, or command output quickly.

## Delegation Guidance
- Delegate execution-heavy tasks to `general_worker`.
- Delegate feature work to `senior_dev`.
- Delegate quick fixes to `junior_fixer`.
- For image/social flows, prefer these tools before publishing:
  - `text_script_create`
  - `image_gen_create`
  - `image_gen_workflow`
  - `community_manager_create_draft`
  - `community_from_image`
