---
name: A tool hit something comp could not do
about: You were building a tool, comp came up short, and you worked around it
title: ''
labels: from-a-tool
---

<!--
Read this first: the useful report is not "tuikit should have X". It is
"here is what I wrote by hand, and here is why comp could not do it".

Decision 31: no use case is a reason to WAIT; not knowing the shape is a reason
to REFUSE. Only the second is a principle. Which of those this is depends
entirely on the evidence below — so please fill it in even when the answer feels
obvious, because "obvious" is how comp.Tree nearly got built, and it would have
been wrong.
-->

## The tool, and what it was doing

<!-- the cloud tool, drawing the resource tree; the database tool, mid-apply;
and so on. -->

## What you wrote by hand

<!-- The actual code, or a link to it. This is the most important section.

A component gets built when two tools have hand-rolled the SAME thing, and
"the same" is a judgement nobody can make from a description. app.Toggles
exists because two tools' hand-rolled versions were shown side by side and
turned out to be identical. comp.Tree does NOT exist because two tools' versions
were shown side by side and turned out to differ in a way that mattered. -->

```go

```

## Why comp could not do it

<!-- Which component you reached for, and what stopped you. "There is no
component" is an answer; so is "comp.List has one but it assumes X". -->

## Did you compute a coordinate?

<!-- If the workaround involved c.Text/c.Fill/c.Set, or arithmetic on a Rect,
say so plainly. That is the strongest possible signal — the whole point of the
canvas is that a tool never does this, so a tool that had to is a tool the
framework failed. -->

## What you would want it to look like

<!-- A guess is fine and a wrong guess is useful. Say if you have none. -->
