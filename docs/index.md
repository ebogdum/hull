---
title: Home
nav_order: 1
description: "Hull is a Kubernetes package manager built around expression-based templating, layered composition, and dependency-aware orchestration."
---

<div class="hull-hero">
  <p class="hull-hero__eyebrow">Kubernetes packaging, reimagined</p>
  <h1>Hull</h1>
  <p class="lead">A Kubernetes package manager built around <strong>expression-based templating</strong>, <strong>layered composition</strong>, and <strong>dependency-aware orchestration</strong> — with a plan/diff/drift workflow you can actually trust.</p>
  <a href="{{ '/guides/quickstart.html' | relative_url }}" class="btn btn-hero">Get started →</a>
  <a href="{{ '/cli/README.html' | relative_url }}" class="btn btn-ghost">CLI reference</a>
  <a href="https://github.com/ebogdum/hull" class="btn btn-ghost">GitHub</a>
</div>

<div class="hull-terminal">
  <div class="hull-terminal__dots"><span></span><span></span><span></span></div>
  <pre><span class="c-prompt">$</span> <span class="c-cmd">hull install web ./web --values prod.yaml</span>
<span class="c-out">Planned 14 resources · 0 drift · 0 warnings</span>
<span class="c-ok">✔ release "web" installed (revision 1)</span>

<span class="c-prompt">$</span> <span class="c-cmd">hull diff web ./web --values prod.yaml</span>
<span class="c-out">~ Deployment/web            replicas 3 → 5</span>
<span class="c-out">~ ConfigMap/web-config      data.LOG_LEVEL info → debug</span>
<span class="c-ok">2 changes, 12 unchanged</span></pre>
</div>

## Start here

New to Hull? The **Quickstart** takes you from an empty directory to a deployed, upgraded, and rolled-back release in a few minutes.

<div class="hull-cards">
  <a class="hull-card" href="{{ '/guides/quickstart.html' | relative_url }}">
    <div class="hull-card__icon">🚀</div>
    <div class="hull-card__title">Quickstart</div>
    <div class="hull-card__desc">Install, upgrade, and roll back your first release.</div>
  </a>
  <a class="hull-card" href="{{ '/guides/packages.html' | relative_url }}">
    <div class="hull-card__icon">📦</div>
    <div class="hull-card__title">Packages</div>
    <div class="hull-card__desc">Structure a Hull package and its templates.</div>
  </a>
  <a class="hull-card" href="{{ '/guides/values.html' | relative_url }}">
    <div class="hull-card__icon">🎛️</div>
    <div class="hull-card__title">Values</div>
    <div class="hull-card__desc">Layer and override configuration cleanly.</div>
  </a>
  <a class="hull-card" href="{{ '/guides/migration.html' | relative_url }}">
    <div class="hull-card__icon">🔀</div>
    <div class="hull-card__title">Migrate from Helm</div>
    <div class="hull-card__desc">Convert an existing Helm chart to a Hull package.</div>
  </a>
</div>

## Explore the docs

<div class="hull-cards">
  <a class="hull-card" href="{{ '/guides/index.html' | relative_url }}">
    <div class="hull-card__icon">📘</div>
    <div class="hull-card__title">Guides</div>
    <div class="hull-card__desc">Task-focused walkthroughs for real work.</div>
  </a>
  <a class="hull-card" href="{{ '/cli/README.html' | relative_url }}">
    <div class="hull-card__icon">⌨️</div>
    <div class="hull-card__title">CLI reference</div>
    <div class="hull-card__desc">Every command, flag, and option — with worked input→output examples.</div>
  </a>
  <a class="hull-card" href="{{ '/templates/index.html' | relative_url }}">
    <div class="hull-card__icon">🧩</div>
    <div class="hull-card__title">Templates</div>
    <div class="hull-card__desc">The <code>${...}</code> expression language, control flow, and ~200 functions.</div>
  </a>
  <a class="hull-card" href="{{ '/reference/index.html' | relative_url }}">
    <div class="hull-card__icon">📑</div>
    <div class="hull-card__title">Reference</div>
    <div class="hull-card__desc">Manifest and values files, field by field.</div>
  </a>
  <a class="hull-card" href="{{ '/comparison.html' | relative_url }}">
    <div class="hull-card__icon">⚖️</div>
    <div class="hull-card__title">Comparison</div>
    <div class="hull-card__desc">Hull vs Helm, Kustomize, kapp, and kpt.</div>
  </a>
  <a class="hull-card" href="{{ '/faq.html' | relative_url }}">
    <div class="hull-card__icon">❓</div>
    <div class="hull-card__title">FAQ</div>
    <div class="hull-card__desc">Frequent questions, answered.</div>
  </a>
</div>

## How these docs are written

Every reference example shows **input → output**: the file you write or the command you run, and the exact result it produces. When a command reads hidden state — a stored release, the live cluster — the docs show that state and trace each output line back to its cause. See [`hull drift`]({{ '/cli/drift.html' | relative_url }}) for the fullest example of this style.
