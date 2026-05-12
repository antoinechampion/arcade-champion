---
name: task
description: >-
  Structured implementation workflow that enforces plan-first development.
  Use when starting any feature, bug fix, or refactor task.
argument-hint: Description of the task to implement
---

# Task Workflow

Enforce the project's 5-step workflow from CLAUDE.md. Do NOT skip steps or start writing code before the plan is confirmed.

## Input

Task description: $ARGUMENTS

## Step 1: Plan

1. Read CLAUDE.md to refresh conventions.
2. Identify which files will be created or modified.
3. List the steps and their dependencies.
4. Present the plan to the user and **wait for confirmation** before proceeding.

## Step 2: Refine

Before implementing, verify:
- Does the plan follow existing patterns in the codebase?
- Are naming conventions consistent with neighboring files?
- Is this the simplest approach (KISS)?

State any adjustments. If none, say so explicitly and proceed.

## Step 3: Implement

Write the code. Follow all conventions from CLAUDE.md:
- Vue Composition API with `<script setup lang="ts">`
- Strict TypeScript, no `any`
- Minimal comments, only when the *why* is non-obvious
- No defensive programming for internal code

## Step 4: Simplify

After implementation, review the change and surrounding code:
- Is there dead code to remove?
- Can anything be flattened or deduplicated?
- Did the change introduce unnecessary abstractions?

State what was simplified, or explicitly confirm nothing needs simplification.

## Step 5: Test

Add or update tests for the new behavior:
- Few but meaningful tests covering realistic scenarios
- No coverage targets, no testing every branch
- Run tests to confirm they pass

## Completion

Summarize what was done across all 5 steps in 2-3 sentences.
