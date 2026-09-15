<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>KubeBridge — Kubernetes-native Integration Orchestration</title>
<style>
  :root {
    --bg: #080d14;
    --bg1: #0d1520;
    --bg2: #111d2e;
    --bg3: #162236;
    --surface: #1a2840;
    --surface2: #1f2f47;
    --border: rgba(56,120,210,0.18);
    --border2: rgba(56,120,210,0.32);
    --text: #e2eaf5;
    --text2: #90a8c8;
    --text3: #5a7a9a;
    --blue: #3b82f6;
    --cyan: #06b6d4;
    --teal: #14b8a6;
    --green: #22c55e;
    --purple: #8b5cf6;
    --amber: #f59e0b;
    --red: #ef4444;
    --glow-blue: rgba(59,130,246,0.15);
    --glow-cyan: rgba(6,182,212,0.12);
    --font: 'Inter', 'SF Pro Display', system-ui, sans-serif;
    --mono: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  }
  * { margin:0; padding:0; box-sizing:border-box; }
  html { scroll-behavior: smooth; }
  body {
    background: var(--bg);
    color: var(--text);
    font-family: var(--font);
    font-size: 15px;
    line-height: 1.7;
    overflow-x: hidden;
  }
  a { color: var(--cyan); text-decoration: none; }
  a:hover { color: var(--blue); }
  code, pre {
    font-family: var(--mono);
    font-size: 13px;
  }
  pre {
    background: var(--bg1);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 20px 24px;
    overflow-x: auto;
    position: relative;
  }
  .tok-k { color: var(--cyan); }
  .tok-s { color: var(--green); }
  .tok-c { color: var(--text3); font-style: italic; }
  .tok-v { color: var(--amber); }
  .tok-p { color: var(--text2); }

  /* NAV */
  nav {
    position: fixed; top:0; left:0; right:0; z-index:100;
    background: rgba(8,13,20,0.88);
    backdrop-filter: blur(20px);
    border-bottom: 1px solid var(--border);
    display: flex; align-items: center; justify-content: space-between;
    padding: 0 40px; height: 56px;
  }
  .nav-logo {
    display: flex; align-items: center; gap: 10px;
    font-weight: 600; font-size: 15px; letter-spacing: -0.01em;
  }
  .nav-logo svg { flex-shrink: 0; }
  .nav-links { display: flex; gap: 28px; list-style: none; }
  .nav-links a { color: var(--text2); font-size: 13.5px; transition: color .2s; }
  .nav-links a:hover { color: var(--text); }
  .nav-star {
    background: linear-gradient(135deg, var(--blue), var(--cyan));
    color: #fff; font-size: 12px; font-weight: 600;
    padding: 6px 14px; border-radius: 6px; letter-spacing: 0.02em;
    transition: opacity .2s;
  }
  .nav-star:hover { opacity: .85; color: #fff; }

  /* SECTIONS */
  section { position: relative; }
  .container { max-width: 1100px; margin: 0 auto; padding: 0 40px; }
  .section-pad { padding: 100px 0; }

  /* HERO */
  #hero {
    min-height: 100vh;
    display: flex; flex-direction: column; align-items: center; justify-content: center;
    text-align: center; padding: 80px 40px 60px;
    overflow: hidden;
    position: relative;
  }
  .hero-grid-bg {
    position: absolute; inset: 0; overflow: hidden; pointer-events: none;
  }
  .hero-grid-bg svg { width: 100%; height: 100%; position: absolute; inset: 0; }
  .hero-eyebrow {
    display: inline-flex; align-items: center; gap: 8px;
    background: rgba(59,130,246,0.1);
    border: 1px solid rgba(59,130,246,0.25);
    border-radius: 100px;
    padding: 5px 14px; font-size: 12px;
    color: var(--cyan); font-weight: 500; letter-spacing: 0.06em;
    margin-bottom: 28px; text-transform: uppercase;
  }
  .hero-eyebrow-dot {
    width: 6px; height: 6px; border-radius: 50%;
    background: var(--cyan);
    animation: pulse-dot 2s ease-in-out infinite;
  }
  @keyframes pulse-dot {
    0%,100% { opacity: 1; transform: scale(1); }
    50% { opacity: .4; transform: scale(.7); }
  }
  h1.hero-title {
    font-size: clamp(42px, 6vw, 72px);
    font-weight: 700; letter-spacing: -0.04em; line-height: 1.05;
    margin-bottom: 24px;
    background: linear-gradient(135deg, #e2eaf5 0%, #7ab3f5 50%, #06b6d4 100%);
    -webkit-background-clip: text; -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  .hero-sub {
    font-size: 18px; color: var(--text2); max-width: 580px;
    margin: 0 auto 36px; line-height: 1.65; font-weight: 400;
  }
  .hero-badges {
    display: flex; flex-wrap: wrap; gap: 8px; justify-content: center;
    margin-bottom: 40px;
  }
  .badge {
    display: inline-flex; align-items: center; gap: 6px;
    padding: 5px 12px; border-radius: 6px; font-size: 12px; font-weight: 500;
    border: 1px solid var(--border);
    background: rgba(255,255,255,0.03);
    color: var(--text2); transition: border-color .2s, background .2s;
  }
  .badge:hover { border-color: var(--border2); background: rgba(255,255,255,0.06); }
  .badge-go { border-color: rgba(0,173,216,0.35); color: #00add8; background: rgba(0,173,216,0.06); }
  .badge-k8s { border-color: rgba(50,108,229,0.35); color: #326ce5; background: rgba(50,108,229,0.06); }
  .badge-docker { border-color: rgba(29,99,237,0.35); color: #1d63ed; background: rgba(29,99,237,0.06); }
  .badge-ci { border-color: rgba(34,197,94,0.35); color: var(--green); background: rgba(34,197,94,0.06); }
  .badge-mit { border-color: rgba(139,92,246,0.35); color: var(--purple); background: rgba(139,92,246,0.06); }
  .hero-ctas {
    display: flex; gap: 12px; flex-wrap: wrap; justify-content: center;
  }
  .btn-primary {
    background: linear-gradient(135deg, var(--blue), var(--cyan));
    color: #fff; padding: 11px 24px; border-radius: 8px;
    font-size: 14px; font-weight: 600; border: none; cursor: pointer;
    transition: opacity .2s, transform .15s;
  }
  .btn-primary:hover { opacity: .88; transform: translateY(-1px); color:#fff; }
  .btn-outline {
    background: transparent; color: var(--text2);
    padding: 11px 24px; border-radius: 8px;
    font-size: 14px; font-weight: 500;
    border: 1px solid var(--border2); cursor: pointer;
    transition: border-color .2s, color .2s, background .2s;
  }
  .btn-outline:hover { border-color: var(--blue); color: var(--text); background: var(--glow-blue); }

  /* ANIMATED FLOW DIAGRAM */
  #flow-visual {
    margin: 70px auto 0; max-width: 800px; width: 100%;
    position: relative;
  }

  /* SECTION HEADERS */
  .section-label {
    font-size: 11px; letter-spacing: .1em; text-transform: uppercase;
    color: var(--cyan); font-weight: 600; margin-bottom: 12px;
  }
  h2.section-title {
    font-size: clamp(26px, 3.5vw, 38px); font-weight: 700;
    letter-spacing: -0.03em; line-height: 1.15; margin-bottom: 16px;
    color: var(--text);
  }
  .section-desc {
    color: var(--text2); font-size: 16px; max-width: 560px; line-height: 1.7;
  }

  /* HOW IT WORKS */
  #how-it-works { background: var(--bg); }
  .how-grid {
    display: grid; grid-template-columns: 1fr 1fr;
    gap: 60px; align-items: start; margin-top: 60px;
  }
  .step-list { display: flex; flex-direction: column; gap: 0; }
  .step {
    display: flex; gap: 20px; padding: 20px 0;
    border-bottom: 1px solid var(--border);
    cursor: pointer; transition: background .2s; border-radius: 8px;
    padding: 16px 12px;
  }
  .step:hover { background: rgba(59,130,246,0.04); }
  .step.active { background: rgba(59,130,246,0.07); }
  .step-num {
    width: 32px; height: 32px; border-radius: 8px; flex-shrink: 0;
    display: flex; align-items: center; justify-content: center;
    font-size: 12px; font-weight: 700;
    background: rgba(59,130,246,0.15);
    color: var(--blue); border: 1px solid rgba(59,130,246,0.3);
    margin-top: 2px;
  }
  .step-content h4 { font-size: 14px; font-weight: 600; margin-bottom: 4px; }
  .step-content p { font-size: 13px; color: var(--text2); line-height: 1.6; }

  /* ARCHITECTURE */
  #architecture { background: var(--bg1); }
  .arch-diagram-wrap {
    margin-top: 60px; border-radius: 16px;
    border: 1px solid var(--border); overflow: hidden;
    background: var(--bg);
  }
  .arch-diagram-header {
    padding: 14px 20px; border-bottom: 1px solid var(--border);
    display: flex; align-items: center; gap: 8px;
  }
  .arch-dot { width: 10px; height: 10px; border-radius: 50%; }
  .arch-dot-r { background: #ef4444; }
  .arch-dot-y { background: #f59e0b; }
  .arch-dot-g { background: #22c55e; }
  .arch-diagram-body { padding: 32px; }

  /* COMPONENTS TABLE */
  .comp-table {
    width: 100%; border-collapse: collapse; margin-top: 40px;
    font-size: 14px;
  }
  .comp-table th {
    text-align: left; padding: 10px 16px;
    font-size: 11px; text-transform: uppercase; letter-spacing: .08em;
    color: var(--text3); border-bottom: 1px solid var(--border);
    font-weight: 600;
  }
  .comp-table td {
    padding: 14px 16px; border-bottom: 1px solid var(--border);
    vertical-align: top; line-height: 1.55;
  }
  .comp-table tr:last-child td { border-bottom: none; }
  .comp-table tr:hover td { background: rgba(59,130,246,0.03); }
  .comp-name { font-weight: 600; color: var(--cyan); white-space: nowrap; }
  .comp-resp { color: var(--text2); }

  /* CONNECTORS */
  #connectors { background: var(--bg); }
  .connector-grid {
    display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 16px; margin-top: 48px;
  }
  .connector-card {
    background: var(--bg1); border: 1px solid var(--border);
    border-radius: 12px; padding: 24px;
    transition: border-color .2s, transform .2s;
    position: relative; overflow: hidden;
  }
  .connector-card:hover { border-color: var(--border2); transform: translateY(-2px); }
  .connector-card.active { border-color: rgba(34,197,94,0.4); }
  .connector-card.planned { border-color: var(--border); opacity: .7; }
  .connector-icon {
    width: 40px; height: 40px; border-radius: 10px;
    display: flex; align-items: center; justify-content: center;
    font-size: 20px; margin-bottom: 14px;
  }
  .connector-card h4 { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
  .connector-card p { font-size: 13px; color: var(--text2); line-height: 1.55; }
  .connector-status {
    display: inline-flex; align-items: center; gap: 5px;
    font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: .06em;
    margin-top: 14px; padding: 3px 9px; border-radius: 4px;
  }
  .status-active { background: rgba(34,197,94,0.1); color: var(--green); }
  .status-planned { background: rgba(90,122,154,0.1); color: var(--text3); }
  .status-dot { width: 5px; height: 5px; border-radius: 50%; background: currentColor; }

  /* YAML */
  #quickstart { background: var(--bg1); }
  .qs-grid {
    display: grid; grid-template-columns: 1fr 1fr; gap: 32px; margin-top: 60px;
  }
  .qs-step {
    display: flex; align-items: flex-start; gap: 12px;
    margin-bottom: 24px; padding-bottom: 24px;
    border-bottom: 1px solid var(--border);
  }
  .qs-step:last-child { border-bottom: none; margin-bottom: 0; padding-bottom: 0; }
  .qs-step-n {
    width: 24px; height: 24px; border-radius: 6px; flex-shrink: 0;
    background: rgba(59,130,246,0.12); color: var(--blue);
    font-size: 11px; font-weight: 700;
    display: flex; align-items: center; justify-content: center;
    margin-top: 2px;
  }
  .qs-step h5 { font-size: 13px; font-weight: 600; margin-bottom: 4px; }
  .qs-step code {
    display: block; background: var(--bg); padding: 8px 12px;
    border-radius: 6px; border: 1px solid var(--border);
    color: var(--green); font-size: 12px; margin-top: 8px;
    line-height: 1.6;
  }

  /* POW */
  #proof { background: var(--bg); }
  .pow-grid {
    display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 12px; margin-top: 48px;
  }
  .pow-card {
    background: var(--bg1); border: 1px solid var(--border);
    border-radius: 10px; padding: 20px;
  }
  .pow-card h5 { font-size: 13px; font-weight: 600; margin-bottom: 4px; }
  .pow-card p { font-size: 12px; color: var(--text2); }
  .pow-check {
    width: 22px; height: 22px; border-radius: 6px;
    background: rgba(34,197,94,0.12); color: var(--green);
    display: flex; align-items: center; justify-content: center;
    font-size: 12px; margin-bottom: 12px;
  }
  .pow-check svg { width: 13px; height: 13px; }
  .val-table { width: 100%; border-collapse: collapse; margin-top: 48px; font-size: 14px; }
  .val-table th {
    text-align: left; padding: 10px 16px;
    font-size: 11px; text-transform: uppercase; letter-spacing: .08em;
    color: var(--text3); border-bottom: 1px solid var(--border); font-weight: 600;
  }
  .val-table td { padding: 13px 16px; border-bottom: 1px solid var(--border); color: var(--text2); }
  .val-table tr:last-child td { border-bottom: none; }
  .val-pass { color: var(--green); font-weight: 600; }

  /* ROADMAP */
  #roadmap { background: var(--bg1); }
  .roadmap-cols {
    display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 24px; margin-top: 56px;
  }
  .roadmap-col { }
  .roadmap-col-header {
    font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: .08em;
    padding: 8px 14px; border-radius: 6px; margin-bottom: 16px;
    display: inline-block;
  }
  .rc-done { background: rgba(34,197,94,0.1); color: var(--green); }
  .rc-soon { background: rgba(59,130,246,0.1); color: var(--blue); }
  .rc-later { background: rgba(90,122,154,0.1); color: var(--text3); }
  .roadmap-item {
    display: flex; align-items: flex-start; gap: 10px;
    padding: 10px 0; border-bottom: 1px solid var(--border);
    font-size: 13.5px; color: var(--text2);
  }
  .roadmap-item:last-child { border-bottom: none; }
  .roadmap-item svg { flex-shrink: 0; margin-top: 3px; }

  /* REPO STRUCTURE */
  #structure { background: var(--bg); }
  .repo-tree {
    background: var(--bg1); border: 1px solid var(--border);
    border-radius: 12px; overflow: hidden; margin-top: 48px;
  }
  .repo-tree-head {
    padding: 14px 20px; border-bottom: 1px solid var(--border);
    display: flex; align-items: center; gap: 8px;
    font-size: 12px; color: var(--text2);
  }
  .repo-tree-body { padding: 24px 28px; }
  .tree-line { display: flex; gap: 12px; align-items: flex-start; margin-bottom: 6px; }
  .tree-path { font-family: var(--mono); font-size: 13px; color: var(--text2); }
  .tree-path .dir { color: var(--blue); font-weight: 500; }
  .tree-path .file { color: var(--text); }
  .tree-desc { font-size: 12px; color: var(--text3); padding-top: 1px; }

  /* SECURITY */
  #security { background: var(--bg1); }
  .security-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px,1fr)); gap: 16px; margin-top: 48px; }
  .sec-card { background: var(--bg); border: 1px solid var(--border); border-radius: 10px; padding: 20px; }
  .sec-icon { font-size: 22px; margin-bottom: 12px; }
  .sec-card h5 { font-size: 14px; font-weight: 600; margin-bottom: 6px; }
  .sec-card p { font-size: 13px; color: var(--text2); line-height: 1.6; }

  /* FOOTER */
  footer {
    background: var(--bg); border-top: 1px solid var(--border);
    padding: 48px 40px; text-align: center;
  }
  .footer-inner { max-width: 1100px; margin: 0 auto; }
  .footer-logo { font-size: 18px; font-weight: 700; letter-spacing: -0.02em; margin-bottom: 12px; }
  .footer-tagline { font-size: 13px; color: var(--text3); margin-bottom: 28px; }
  .footer-links { display: flex; gap: 24px; justify-content: center; flex-wrap: wrap; margin-bottom: 28px; }
  .footer-links a { color: var(--text3); font-size: 13px; transition: color .2s; }
  .footer-links a:hover { color: var(--text2); }
  .footer-copy { font-size: 12px; color: var(--text3); }

  /* DIVIDERS */
  .divider { height: 1px; background: var(--border); margin: 0; }

  /* ANIMATIONS */
  @keyframes fadeInUp {
    from { opacity: 0; transform: translateY(24px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .anim-in { animation: fadeInUp .6s ease both; }
  .anim-in-d1 { animation-delay: .1s; }
  .anim-in-d2 { animation-delay: .2s; }
  .anim-in-d3 { animation-delay: .3s; }
  .anim-in-d4 { animation-delay: .4s; }
  .anim-in-d5 { animation-delay: .55s; }

  /* RESPONSIVE */
  @media (max-width: 768px) {
    .container { padding: 0 20px; }
    nav { padding: 0 20px; }
    .nav-links { display: none; }
    .how-grid, .qs-grid, .roadmap-cols { grid-template-columns: 1fr; }
    h1.hero-title { font-size: 36px; }
    .section-pad { padding: 70px 0; }
  }
</style>
</head>
<body>

<!-- NAV -->
<nav>
  <div class="nav-logo">
    <svg width="22" height="22" viewBox="0 0 22 22" fill="none">
      <rect x="1" y="1" width="20" height="20" rx="5" stroke="#3b82f6" stroke-width="1.5"/>
      <path d="M6 11h10M11 6v10" stroke="#06b6d4" stroke-width="1.5" stroke-linecap="round"/>
      <circle cx="6" cy="6" r="1.5" fill="#3b82f6"/>
      <circle cx="16" cy="6" r="1.5" fill="#06b6d4"/>
      <circle cx="6" cy="16" r="1.5" fill="#06b6d4"/>
      <circle cx="16" cy="16" r="1.5" fill="#3b82f6"/>
    </svg>
    KubeBridge
  </div>
  <ul class="nav-links">
    <li><a href="#how-it-works">How it works</a></li>
    <li><a href="#architecture">Architecture</a></li>
    <li><a href="#connectors">Connectors</a></li>
    <li><a href="#quickstart">Quick start</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
  </ul>
  <a href="https://github.com" class="nav-star">Star on GitHub</a>
</nav>

<!-- HERO -->
<section id="hero">
  <div class="hero-grid-bg">
    <svg viewBox="0 0 1400 800" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="xMidYMid slice">
      <defs>
        <radialGradient id="rg1" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stop-color="#3b82f6" stop-opacity="0.08"/>
          <stop offset="100%" stop-color="#080d14" stop-opacity="0"/>
        </radialGradient>
        <radialGradient id="rg2" cx="20%" cy="30%" r="40%">
          <stop offset="0%" stop-color="#06b6d4" stop-opacity="0.06"/>
          <stop offset="100%" stop-color="#080d14" stop-opacity="0"/>
        </radialGradient>
      </defs>
      <rect width="1400" height="800" fill="url(#rg1)"/>
      <rect width="1400" height="800" fill="url(#rg2)"/>
      <!-- grid lines -->
      <g stroke="rgba(59,130,246,0.06)" stroke-width="0.5">
        <line x1="0" y1="100" x2="1400" y2="100"/><line x1="0" y1="200" x2="1400" y2="200"/>
        <line x1="0" y1="300" x2="1400" y2="300"/><line x1="0" y1="400" x2="1400" y2="400"/>
        <line x1="0" y1="500" x2="1400" y2="500"/><line x1="0" y1="600" x2="1400" y2="600"/>
        <line x1="0" y1="700" x2="1400" y2="700"/>
        <line x1="100" y1="0" x2="100" y2="800"/><line x1="200" y1="0" x2="200" y2="800"/>
        <line x1="300" y1="0" x2="300" y2="800"/><line x1="400" y1="0" x2="400" y2="800"/>
        <line x1="500" y1="0" x2="500" y2="800"/><line x1="600" y1="0" x2="600" y2="800"/>
        <line x1="700" y1="0" x2="700" y2="800"/><line x1="800" y1="0" x2="800" y2="800"/>
        <line x1="900" y1="0" x2="900" y2="800"/><line x1="1000" y1="0" x2="1000" y2="800"/>
        <line x1="1100" y1="0" x2="1100" y2="800"/><line x1="1200" y1="0" x2="1200" y2="800"/>
        <line x1="1300" y1="0" x2="1300" y2="800"/>
      </g>
      <!-- floating dots -->
      <g fill="rgba(59,130,246,0.5)">
        <circle cx="100" cy="100" r="2"><animate attributeName="opacity" values="0.5;0.1;0.5" dur="3s" repeatCount="indefinite"/></circle>
        <circle cx="300" cy="200" r="1.5"><animate attributeName="opacity" values="0.3;0.8;0.3" dur="4s" repeatCount="indefinite"/></circle>
        <circle cx="700" cy="100" r="2"><animate attributeName="opacity" values="0.6;0.1;0.6" dur="2.5s" repeatCount="indefinite"/></circle>
        <circle cx="1100" cy="200" r="1.5"><animate attributeName="opacity" values="0.2;0.7;0.2" dur="3.5s" repeatCount="indefinite"/></circle>
        <circle cx="1300" cy="100" r="2"><animate attributeName="opacity" values="0.5;0.2;0.5" dur="5s" repeatCount="indefinite"/></circle>
      </g>
      <g fill="rgba(6,182,212,0.4)">
        <circle cx="200" cy="400" r="1.5"><animate attributeName="opacity" values="0.4;0;0.4" dur="4s" repeatCount="indefinite"/></circle>
        <circle cx="900" cy="300" r="2"><animate attributeName="opacity" values="0.6;0.2;0.6" dur="3s" repeatCount="indefinite"/></circle>
        <circle cx="1200" cy="500" r="1.5"><animate attributeName="opacity" values="0.3;0.7;0.3" dur="5s" repeatCount="indefinite"/></circle>
      </g>
    </svg>
  </div>

  <div class="hero-eyebrow anim-in">
    <span class="hero-eyebrow-dot"></span>
    Open Source · Kubernetes Native · Go
  </div>
  <h1 class="hero-title anim-in anim-in-d1">KubeBridge</h1>
  <p class="hero-sub anim-in anim-in-d2">
    Kubernetes-native integration orchestration. Define external service connections as first-class Kubernetes resources — reconciled, observable, and declarative.
  </p>

  <div class="hero-badges anim-in anim-in-d3">
    <span class="badge badge-go">
      <svg width="14" height="14" viewBox="0 0 40 15" fill="none"><text x="0" y="12" font-size="12" fill="#00add8" font-weight="700">Go</text></svg>
      Go 1.21+
    </span>
    <span class="badge badge-k8s">⎈ Kubernetes 1.26+</span>
    <span class="badge badge-docker">🐳 Docker</span>
    <span class="badge badge-ci">
      <svg width="8" height="8" viewBox="0 0 8 8"><circle cx="4" cy="4" r="4" fill="currentColor"/></svg>
      CI Passing
    </span>
    <span class="badge badge-mit">
      <svg width="10" height="10" viewBox="0 0 10 10" fill="none"><path d="M5 1L6.2 3.5H9L6.9 5.3L7.6 8L5 6.5L2.4 8L3.1 5.3L1 3.5H3.8L5 1Z" fill="currentColor"/></svg>
      MIT License
    </span>
    <span class="badge">controller-runtime</span>
    <span class="badge">CRD-based</span>
    <span class="badge">RBAC Scoped</span>
  </div>

  <div class="hero-ctas anim-in anim-in-d4">
    <a href="#quickstart" class="btn-primary">Quick Start</a>
    <a href="#architecture" class="btn-outline">View Architecture</a>
    <a href="https://github.com" class="btn-outline" style="gap:6px; display:inline-flex; align-items:center;">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12c0 4.42 2.87 8.17 6.84 9.49.5.09.68-.22.68-.48v-1.7c-2.78.6-3.37-1.34-3.37-1.34-.46-1.16-1.11-1.47-1.11-1.47-.91-.62.07-.61.07-.61 1 .07 1.53 1.03 1.53 1.03.89 1.52 2.34 1.08 2.91.83.09-.65.35-1.08.63-1.33-2.22-.25-4.56-1.11-4.56-4.94 0-1.09.39-1.98 1.03-2.68-.1-.25-.45-1.27.1-2.64 0 0 .84-.27 2.75 1.02A9.56 9.56 0 0 1 12 6.8c.85 0 1.71.11 2.51.33 1.91-1.29 2.75-1.02 2.75-1.02.55 1.37.2 2.39.1 2.64.64.7 1.03 1.59 1.03 2.68 0 3.84-2.34 4.68-4.57 4.93.36.31.68.92.68 1.85v2.74c0 .27.18.58.69.48C19.14 20.16 22 16.42 22 12c0-5.52-4.48-10-10-10z"/></svg>
      GitHub
    </a>
  </div>

  <!-- Animated Flow Visual -->
  <div id="flow-visual" class="anim-in anim-in-d5">
    <svg id="hero-flow" width="100%" viewBox="0 0 780 220" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <filter id="glow-blue" x="-30%" y="-30%" width="160%" height="160%">
          <feGaussianBlur stdDeviation="3" result="blur"/>
          <feComposite in="SourceGraphic" in2="blur" operator="over"/>
        </filter>
        <marker id="arr" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto-start-reverse">
          <path d="M2 1L8 5L2 9" fill="none" stroke="#3b82f6" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
        </marker>
        <marker id="arr-cyan" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto-start-reverse">
          <path d="M2 1L8 5L2 9" fill="none" stroke="#06b6d4" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
        </marker>
        <marker id="arr-green" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto-start-reverse">
          <path d="M2 1L8 5L2 9" fill="none" stroke="#22c55e" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
        </marker>
      </defs>

      <!-- Connecting lines with animation -->
      <!-- K8s API -> Connector CRD -->
      <line x1="155" y1="110" x2="235" y2="110" stroke="#3b82f6" stroke-width="1.5" marker-end="url(#arr)" opacity="0.7">
        <animate attributeName="opacity" values="0.4;0.9;0.4" dur="2s" repeatCount="indefinite"/>
      </line>
      <!-- Connector CRD -> KubeBridge -->
      <line x1="355" y1="110" x2="435" y2="110" stroke="#06b6d4" stroke-width="1.5" marker-end="url(#arr-cyan)" opacity="0.7">
        <animate attributeName="opacity" values="0.4;0.9;0.4" dur="2s" begin="0.5s" repeatCount="indefinite"/>
      </line>
      <!-- KubeBridge -> External API -->
      <line x1="555" y1="110" x2="635" y2="110" stroke="#22c55e" stroke-width="1.5" marker-end="url(#arr-green)" opacity="0.7">
        <animate attributeName="opacity" values="0.4;0.9;0.4" dur="2s" begin="1s" repeatCount="indefinite"/>
      </line>

      <!-- Animated data packets on lines -->
      <circle r="4" fill="#3b82f6" opacity="0.9">
        <animateMotion dur="2s" repeatCount="indefinite" begin="0s">
          <mpath href="#path1"/>
        </animateMotion>
      </circle>
      <path id="path1" d="M155 110 L235 110" fill="none" stroke="none"/>

      <circle r="4" fill="#06b6d4" opacity="0.9">
        <animateMotion dur="2s" repeatCount="indefinite" begin="0.7s">
          <mpath href="#path2"/>
        </animateMotion>
      </circle>
      <path id="path2" d="M355 110 L435 110" fill="none" stroke="none"/>

      <circle r="4" fill="#22c55e" opacity="0.9">
        <animateMotion dur="2s" repeatCount="indefinite" begin="1.4s">
          <mpath href="#path3"/>
        </animateMotion>
      </circle>
      <path id="path3" d="M555 110 L635 110" fill="none" stroke="none"/>

      <!-- NODE 1: Kubernetes API Server -->
      <g>
        <rect x="10" y="72" width="145" height="76" rx="10" fill="#0d1520" stroke="rgba(59,130,246,0.4)" stroke-width="1"/>
        <rect x="10" y="72" width="145" height="76" rx="10" fill="url(#ng1)"/>
        <!-- K8s wheel icon simplified -->
        <circle cx="82.5" cy="95" r="10" fill="none" stroke="#326ce5" stroke-width="1.5"/>
        <circle cx="82.5" cy="95" r="3.5" fill="#326ce5"/>
        <line x1="82.5" y1="85" x2="82.5" y2="83" stroke="#326ce5" stroke-width="1.5"/>
        <line x1="82.5" y1="105" x2="82.5" y2="107" stroke="#326ce5" stroke-width="1.5"/>
        <line x1="72.5" y1="95" x2="70.5" y2="95" stroke="#326ce5" stroke-width="1.5"/>
        <line x1="92.5" y1="95" x2="94.5" y2="95" stroke="#326ce5" stroke-width="1.5"/>
        <text x="82.5" y="120" text-anchor="middle" font-size="11" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="600">Kubernetes API</text>
        <text x="82.5" y="135" text-anchor="middle" font-size="10" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Desired State Store</text>
      </g>

      <!-- NODE 2: Connector CRD -->
      <g>
        <rect x="235" y="72" width="120" height="76" rx="10" fill="#0d1520" stroke="rgba(6,182,212,0.4)" stroke-width="1"/>
        <text x="295" y="97" text-anchor="middle" font-size="16" fill="#06b6d4">⚙</text>
        <text x="295" y="118" text-anchor="middle" font-size="11" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="600">Connector CRD</text>
        <text x="295" y="133" text-anchor="middle" font-size="10" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Custom Resource</text>
      </g>

      <!-- NODE 3: KubeBridge Controller (central, bigger) -->
      <g>
        <rect x="435" y="58" width="120" height="104" rx="10" fill="#0f1e33" stroke="rgba(59,130,246,0.6)" stroke-width="1.5"/>
        <rect x="435" y="58" width="120" height="5" rx="2.5" fill="#3b82f6" opacity="0.6"/>
        <text x="495" y="88" text-anchor="middle" font-size="20" fill="#3b82f6">⬡</text>
        <text x="495" y="111" text-anchor="middle" font-size="11" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="700">KubeBridge</text>
        <text x="495" y="126" text-anchor="middle" font-size="10" fill="#6090b8" font-family="Inter,system-ui,sans-serif">Controller</text>
        <text x="495" y="143" text-anchor="middle" font-size="9.5" fill="#3b82f6" font-family="Inter,system-ui,sans-serif">controller-runtime</text>
        <!-- Pulse ring -->
        <circle cx="495" cy="88" r="18" fill="none" stroke="#3b82f6" stroke-width="0.5" opacity="0.3">
          <animate attributeName="r" values="18;26;18" dur="2.5s" repeatCount="indefinite"/>
          <animate attributeName="opacity" values="0.4;0;0.4" dur="2.5s" repeatCount="indefinite"/>
        </circle>
      </g>

      <!-- NODE 4: External API -->
      <g>
        <rect x="635" y="72" width="130" height="76" rx="10" fill="#0d1520" stroke="rgba(34,197,94,0.35)" stroke-width="1"/>
        <text x="700" y="97" text-anchor="middle" font-size="16" fill="#22c55e">⬡</text>
        <text x="700" y="118" text-anchor="middle" font-size="11" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="600">External API</text>
        <text x="700" y="133" text-anchor="middle" font-size="10" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">GitHub · and more</text>
      </g>

      <!-- Labels below nodes -->
      <text x="82.5" y="170" text-anchor="middle" font-size="9" fill="rgba(59,130,246,0.6)" font-family="Inter,system-ui,sans-serif" letter-spacing="0.06em">WATCH</text>
      <text x="295" y="170" text-anchor="middle" font-size="9" fill="rgba(6,182,212,0.6)" font-family="Inter,system-ui,sans-serif" letter-spacing="0.06em">DECLARE</text>
      <text x="495" y="178" text-anchor="middle" font-size="9" fill="rgba(59,130,246,0.6)" font-family="Inter,system-ui,sans-serif" letter-spacing="0.06em">RECONCILE</text>
      <text x="700" y="170" text-anchor="middle" font-size="9" fill="rgba(34,197,94,0.6)" font-family="Inter,system-ui,sans-serif" letter-spacing="0.06em">INTEGRATE</text>
    </svg>
  </div>
</section>

<!-- HOW IT WORKS -->
<section id="how-it-works" class="section-pad">
  <div class="container">
    <div class="how-grid">
      <div>
        <p class="section-label">How it works</p>
        <h2 class="section-title">Declarative integrations, continuous reconciliation</h2>
        <p class="section-desc">KubeBridge follows the Kubernetes operator pattern. You define what connections you want — the controller ensures they exist.</p>
        <div class="step-list" style="margin-top:36px">
          <div class="step active" id="step1">
            <div class="step-num">01</div>
            <div class="step-content">
              <h4>Declare a Connector resource</h4>
              <p>Write a <code>Connector</code> YAML manifest specifying the provider, credentials reference, and configuration. Apply it with <code>kubectl apply</code>.</p>
            </div>
          </div>
          <div class="step" id="step2">
            <div class="step-num">02</div>
            <div class="step-content">
              <h4>Controller watches the resource</h4>
              <p>The KubeBridge controller receives a watch event from the Kubernetes API Server and queues a reconciliation request.</p>
            </div>
          </div>
          <div class="step" id="step3">
            <div class="step-num">03</div>
            <div class="step-content">
              <h4>Credentials retrieved from Secrets</h4>
              <p>The controller reads the referenced Kubernetes Secret to obtain API credentials — never storing them in the CRD spec itself.</p>
            </div>
          </div>
          <div class="step" id="step4">
            <div class="step-num">04</div>
            <div class="step-content">
              <h4>Connector adapter executes</h4>
              <p>Provider-specific adapter logic handles the external API call. The result is reflected back into the resource's <code>.status</code> field.</p>
            </div>
          </div>
          <div class="step" id="step5">
            <div class="step-num">05</div>
            <div class="step-content">
              <h4>Status observed, retries scheduled</h4>
              <p>On failure the controller re-queues with backoff. On success, the status transitions to <code>Ready</code> and continues watching for drift.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Reconciliation SVG -->
      <div style="position:sticky;top:80px;align-self:start;">
        <svg width="100%" viewBox="0 0 340 520" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <marker id="ra" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto-start-reverse">
              <path d="M2 1L8 5L2 9" fill="none" stroke="#3b82f6" stroke-width="1.8" stroke-linecap="round"/>
            </marker>
          </defs>

          <!-- Title -->
          <text x="170" y="22" text-anchor="middle" font-size="11" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif" letter-spacing="0.08em" text-transform="uppercase">RECONCILIATION LOOP</text>

          <!-- Nodes -->
          <!-- Resource Created -->
          <rect x="60" y="38" width="220" height="44" rx="8" fill="#0d1520" stroke="rgba(59,130,246,0.35)" stroke-width="1"/>
          <text x="170" y="61" text-anchor="middle" font-size="12.5" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="500" dominant-baseline="central">Resource Created / Updated</text>

          <!-- Arrow -->
          <line x1="170" y1="82" x2="170" y2="100" stroke="#3b82f6" stroke-width="1.5" marker-end="url(#ra)" opacity="0.7"/>
          <circle r="3" fill="#3b82f6" opacity="0.7">
            <animateMotion dur="3s" repeatCount="indefinite" begin="0s">
              <mpath href="#rp1"/>
            </animateMotion>
          </circle>
          <path id="rp1" d="M170 82 L170 100" fill="none"/>

          <!-- Watch event -->
          <rect x="60" y="100" width="220" height="44" rx="8" fill="#0d1520" stroke="rgba(6,182,212,0.35)" stroke-width="1"/>
          <text x="170" y="123" text-anchor="middle" font-size="12.5" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="500" dominant-baseline="central">Watch Event Received</text>

          <line x1="170" y1="144" x2="170" y2="162" stroke="#3b82f6" stroke-width="1.5" marker-end="url(#ra)" opacity="0.7"/>
          <circle r="3" fill="#06b6d4" opacity="0.7">
            <animateMotion dur="3s" repeatCount="indefinite" begin="0.6s">
              <mpath href="#rp2"/>
            </animateMotion>
          </circle>
          <path id="rp2" d="M170 144 L170 162" fill="none"/>

          <!-- Read Config -->
          <rect x="60" y="162" width="220" height="44" rx="8" fill="#0d1520" stroke="rgba(59,130,246,0.35)" stroke-width="1"/>
          <text x="170" y="185" text-anchor="middle" font-size="12.5" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="500" dominant-baseline="central">Read Connector Config</text>

          <line x1="170" y1="206" x2="170" y2="224" stroke="#3b82f6" stroke-width="1.5" marker-end="url(#ra)" opacity="0.7"/>
          <circle r="3" fill="#3b82f6" opacity="0.7">
            <animateMotion dur="3s" repeatCount="indefinite" begin="1.0s">
              <mpath href="#rp3"/>
            </animateMotion>
          </circle>
          <path id="rp3" d="M170 206 L170 224" fill="none"/>

          <!-- Secret Fetch -->
          <rect x="60" y="224" width="220" height="44" rx="8" fill="#0d1520" stroke="rgba(139,92,246,0.35)" stroke-width="1"/>
          <text x="170" y="247" text-anchor="middle" font-size="12.5" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="500" dominant-baseline="central">Fetch Secret Credentials</text>

          <line x1="170" y1="268" x2="170" y2="286" stroke="#3b82f6" stroke-width="1.5" marker-end="url(#ra)" opacity="0.7"/>
          <circle r="3" fill="#8b5cf6" opacity="0.7">
            <animateMotion dur="3s" repeatCount="indefinite" begin="1.4s">
              <mpath href="#rp4"/>
            </animateMotion>
          </circle>
          <path id="rp4" d="M170 268 L170 286" fill="none"/>

          <!-- Execute -->
          <rect x="60" y="286" width="220" height="44" rx="8" fill="#0d1520" stroke="rgba(6,182,212,0.4)" stroke-width="1.5"/>
          <text x="170" y="309" text-anchor="middle" font-size="12.5" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">Execute Connector Adapter</text>

          <line x1="170" y1="330" x2="170" y2="348" stroke="#3b82f6" stroke-width="1.5" marker-end="url(#ra)" opacity="0.7"/>
          <circle r="3" fill="#06b6d4" opacity="0.7">
            <animateMotion dur="3s" repeatCount="indefinite" begin="1.8s">
              <mpath href="#rp5"/>
            </animateMotion>
          </circle>
          <path id="rp5" d="M170 330 L170 348" fill="none"/>

          <!-- Update Status -->
          <rect x="60" y="348" width="220" height="44" rx="8" fill="#0d1520" stroke="rgba(34,197,94,0.35)" stroke-width="1"/>
          <text x="170" y="371" text-anchor="middle" font-size="12.5" fill="#90a8c8" font-family="Inter,system-ui,sans-serif" font-weight="500" dominant-baseline="central">Update .status</text>

          <!-- Retry / Continue loop arrow -->
          <path d="M280 368 Q320 368 320 430 Q320 490 170 490" stroke="#3b82f6" stroke-width="1" stroke-dasharray="4 3" fill="none" opacity="0.4" marker-end="url(#ra)"/>
          <text x="318" y="415" text-anchor="middle" font-size="10" fill="#3b4e6a" font-family="Inter,system-ui,sans-serif" transform="rotate(90,318,415)">requeue / watch</text>

          <!-- Bottom label -->
          <rect x="60" y="470" width="220" height="34" rx="8" fill="rgba(59,130,246,0.07)" stroke="rgba(59,130,246,0.2)" stroke-width="1"/>
          <text x="170" y="488" text-anchor="middle" font-size="11.5" fill="#5a8abd" font-family="Inter,system-ui,sans-serif" dominant-baseline="central">Idempotent — safe to re-run</text>
        </svg>
      </div>
    </div>
  </div>
</section>

<!-- ARCHITECTURE -->
<section id="architecture" class="section-pad">
  <div class="container">
    <p class="section-label">Architecture</p>
    <h2 class="section-title">Components and responsibilities</h2>
    <p class="section-desc">KubeBridge follows the standard Kubernetes controller architecture — API types declared as CRDs, reconciliation logic in Go, and provider adapters isolated from controller logic.</p>

    <div class="arch-diagram-wrap">
      <div class="arch-diagram-header">
        <div class="arch-dot arch-dot-r"></div>
        <div class="arch-dot arch-dot-y"></div>
        <div class="arch-dot arch-dot-g"></div>
        <span style="font-size:12px;color:var(--text3);margin-left:4px;">KubeBridge Architecture</span>
      </div>
      <div class="arch-diagram-body">
        <svg width="100%" viewBox="0 0 960 380" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <marker id="ba" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto-start-reverse">
              <path d="M2 1L8 5L2 9" fill="none" stroke="#3b82f6" stroke-width="1.8" stroke-linecap="round"/>
            </marker>
            <marker id="ba2" viewBox="0 0 10 10" refX="8" refY="5" markerWidth="5" markerHeight="5" orient="auto-start-reverse">
              <path d="M2 1L8 5L2 9" fill="none" stroke="#5a7a9a" stroke-width="1.8" stroke-linecap="round"/>
            </marker>
          </defs>

          <!-- Layer: Kubernetes Platform -->
          <rect x="10" y="10" width="940" height="110" rx="12" fill="rgba(50,108,229,0.04)" stroke="rgba(50,108,229,0.15)" stroke-width="1" stroke-dasharray="5 3"/>
          <text x="24" y="30" font-size="10" fill="#3264a0" font-family="Inter,system-ui,sans-serif" letter-spacing="0.08em" font-weight="600">KUBERNETES PLATFORM</text>
          <!-- API Server -->
          <rect x="30" y="42" width="160" height="60" rx="8" fill="#0a1525" stroke="rgba(50,108,229,0.4)" stroke-width="1"/>
          <text x="110" y="69" text-anchor="middle" font-size="12" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">API Server</text>
          <text x="110" y="86" text-anchor="middle" font-size="10.5" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Stores desired state</text>
          <!-- RBAC -->
          <rect x="210" y="42" width="130" height="60" rx="8" fill="#0a1525" stroke="rgba(50,108,229,0.3)" stroke-width="1"/>
          <text x="275" y="69" text-anchor="middle" font-size="12" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">RBAC</text>
          <text x="275" y="86" text-anchor="middle" font-size="10.5" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Least privilege</text>
          <!-- Connector CRD -->
          <rect x="360" y="42" width="160" height="60" rx="8" fill="#0a1525" stroke="rgba(6,182,212,0.4)" stroke-width="1"/>
          <text x="440" y="69" text-anchor="middle" font-size="12" fill="#06b6d4" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">Connector CRD</text>
          <text x="440" y="86" text-anchor="middle" font-size="10.5" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Integration config</text>
          <!-- Secret -->
          <rect x="540" y="42" width="130" height="60" rx="8" fill="#0a1525" stroke="rgba(139,92,246,0.35)" stroke-width="1"/>
          <text x="605" y="69" text-anchor="middle" font-size="12" fill="#a78bfa" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">Secret</text>
          <text x="605" y="86" text-anchor="middle" font-size="10.5" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">API credentials</text>

          <!-- KubeBridge layer -->
          <rect x="10" y="140" width="580" height="110" rx="12" fill="rgba(59,130,246,0.05)" stroke="rgba(59,130,246,0.2)" stroke-width="1" stroke-dasharray="5 3"/>
          <text x="24" y="160" font-size="10" fill="#3b5a90" font-family="Inter,system-ui,sans-serif" letter-spacing="0.08em" font-weight="600">KUBEBRIDGE CONTROLLER</text>
          <!-- Reconciler -->
          <rect x="30" y="172" width="170" height="60" rx="8" fill="#0a1525" stroke="rgba(59,130,246,0.45)" stroke-width="1.5"/>
          <text x="115" y="199" text-anchor="middle" font-size="12" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="700" dominant-baseline="central">Reconciler</text>
          <text x="115" y="216" text-anchor="middle" font-size="10.5" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Watch · queue · loop</text>
          <!-- Factory -->
          <rect x="220" y="172" width="160" height="60" rx="8" fill="#0a1525" stroke="rgba(59,130,246,0.35)" stroke-width="1"/>
          <text x="300" y="199" text-anchor="middle" font-size="12" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">Connector Factory</text>
          <text x="300" y="216" text-anchor="middle" font-size="10.5" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Adapter selection</text>
          <!-- Status writer -->
          <rect x="400" y="172" width="160" height="60" rx="8" fill="#0a1525" stroke="rgba(59,130,246,0.3)" stroke-width="1"/>
          <text x="480" y="199" text-anchor="middle" font-size="12" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">Status Writer</text>
          <text x="480" y="216" text-anchor="middle" font-size="10.5" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Observed state updates</text>

          <!-- Connectors layer -->
          <rect x="10" y="270" width="580" height="100" rx="12" fill="rgba(6,182,212,0.04)" stroke="rgba(6,182,212,0.15)" stroke-width="1" stroke-dasharray="5 3"/>
          <text x="24" y="290" font-size="10" fill="#0e7a8a" font-family="Inter,system-ui,sans-serif" letter-spacing="0.08em" font-weight="600">CONNECTOR ADAPTERS</text>
          <!-- GitHub connector -->
          <rect x="30" y="300" width="160" height="52" rx="8" fill="#0a1525" stroke="rgba(34,197,94,0.4)" stroke-width="1"/>
          <text x="110" y="322" text-anchor="middle" font-size="12" fill="#22c55e" font-family="Inter,system-ui,sans-serif" font-weight="600" dominant-baseline="central">GitHub</text>
          <text x="110" y="338" text-anchor="middle" font-size="10" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">Implemented</text>
          <!-- Placeholder 2 -->
          <rect x="210" y="300" width="130" height="52" rx="8" fill="#0a1525" stroke="rgba(90,122,154,0.2)" stroke-width="1" stroke-dasharray="4 3"/>
          <text x="275" y="322" text-anchor="middle" font-size="12" fill="#4a6a8a" font-family="Inter,system-ui,sans-serif" font-weight="500" dominant-baseline="central">Future</text>
          <text x="275" y="338" text-anchor="middle" font-size="10" fill="#3a5a7a" font-family="Inter,system-ui,sans-serif">Planned</text>
          <!-- Placeholder 3 -->
          <rect x="360" y="300" width="130" height="52" rx="8" fill="#0a1525" stroke="rgba(90,122,154,0.2)" stroke-width="1" stroke-dasharray="4 3"/>
          <text x="425" y="322" text-anchor="middle" font-size="12" fill="#4a6a8a" font-family="Inter,system-ui,sans-serif" font-weight="500" dominant-baseline="central">Future</text>
          <text x="425" y="338" text-anchor="middle" font-size="10" fill="#3a5a7a" font-family="Inter,system-ui,sans-serif">Planned</text>

          <!-- External API box -->
          <rect x="730" y="260" width="210" height="100" rx="12" fill="#0a1525" stroke="rgba(34,197,94,0.3)" stroke-width="1"/>
          <text x="835" y="295" text-anchor="middle" font-size="13" fill="#e2eaf5" font-family="Inter,system-ui,sans-serif" font-weight="600">External Services</text>
          <text x="835" y="314" text-anchor="middle" font-size="11" fill="#5a7a9a" font-family="Inter,system-ui,sans-serif">GitHub API</text>
          <text x="835" y="330" text-anchor="middle" font-size="11" fill="#3a5a7a" font-family="Inter,system-ui,sans-serif">Future integrations</text>

          <!-- Arrows: API Server → Reconciler -->
          <path d="M110 102 L115 172" stroke="#3b82f6" stroke-width="1" fill="none" marker-end="url(#ba)" opacity="0.5"/>
          <!-- CRD → Reconciler -->
          <path d="M440 102 L300 172" stroke="#06b6d4" stroke-width="1" fill="none" marker-end="url(#ba)" opacity="0.5"/>
          <!-- Secret → Reconciler -->
          <path d="M605 102 L480 172" stroke="#8b5cf6" stroke-width="1" fill="none" marker-end="url(#ba)" opacity="0.5"/>
          <!-- Reconciler → Factory -->
          <line x1="200" y1="202" x2="218" y2="202" stroke="#3b82f6" stroke-width="1" fill="none" marker-end="url(#ba)" opacity="0.6"/>
          <!-- Factory → Status -->
          <line x1="380" y1="202" x2="398" y2="202" stroke="#3b82f6" stroke-width="1" fill="none" marker-end="url(#ba)" opacity="0.6"/>
          <!-- Factory → Connectors -->
          <path d="M300 232 L300 270" stroke="#06b6d4" stroke-width="1" fill="none" marker-end="url(#ba)" opacity="0.5"/>
          <!-- GitHub → External -->
          <path d="M590 326 L728 310" stroke="#22c55e" stroke-width="1" fill="none" marker-end="url(#ba)" opacity="0.6">
            <animate attributeName="opacity" values="0.3;0.7;0.3" dur="2.5s" repeatCount="indefinite"/>
          </path>
        </svg>
      </div>
    </div>

    <table class="comp-table" style="margin-top:40px">
      <thead>
        <tr>
          <th>Component</th>
          <th>Responsibility</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr><td class="comp-name">API Server</td><td class="comp-resp">Stores and serves the desired state of all Connector resources</td><td><span style="color:var(--blue);font-size:12px">Kubernetes built-in</span></td></tr>
        <tr><td class="comp-name">Connector CRD</td><td class="comp-resp">Custom Resource Definition that defines the schema for integration configurations</td><td><span style="color:var(--green);font-size:12px">Implemented</span></td></tr>
        <tr><td class="comp-name">Controller</td><td class="comp-resp">Watches Connector resources, reconciles desired state with observed state via controller-runtime</td><td><span style="color:var(--green);font-size:12px">Implemented</span></td></tr>
        <tr><td class="comp-name">Secret</td><td class="comp-resp">Kubernetes Secret holding provider credentials; referenced by name from the Connector spec</td><td><span style="color:var(--green);font-size:12px">Implemented</span></td></tr>
        <tr><td class="comp-name">Connector Factory</td><td class="comp-resp">Selects and instantiates the correct provider adapter based on the connector type field</td><td><span style="color:var(--green);font-size:12px">Implemented</span></td></tr>
        <tr><td class="comp-name">GitHub Adapter</td><td class="comp-resp">Provider-specific logic for communicating with the GitHub API; isolated from controller core</td><td><span style="color:var(--green);font-size:12px">Implemented</span></td></tr>
        <tr><td class="comp-name">Status Writer</td><td class="comp-resp">Updates the <code>.status</code> subresource with observed state, conditions, and error messages</td><td><span style="color:var(--green);font-size:12px">Implemented</span></td></tr>
        <tr><td class="comp-name">RBAC</td><td class="comp-resp">ClusterRole and bindings scoped to minimum required permissions for controller operation</td><td><span style="color:var(--green);font-size:12px">Implemented</span></td></tr>
      </tbody>
    </table>
  </div>
</section>

<!-- CONNECTORS -->
<section id="connectors" class="section-pad">
  <div class="container">
    <p class="section-label">Connectors</p>
    <h2 class="section-title">Provider adapter architecture</h2>
    <p class="section-desc">Provider-specific logic lives in isolated adapters that implement a common interface. The controller core remains unchanged when adding new providers.</p>

    <div class="connector-grid">
      <div class="connector-card active">
        <div class="connector-icon" style="background:rgba(34,197,94,0.1);">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="#22c55e"><path d="M12 2C6.48 2 2 6.48 2 12c0 4.42 2.87 8.17 6.84 9.49.5.09.68-.22.68-.48v-1.7c-2.78.6-3.37-1.34-3.37-1.34-.46-1.16-1.11-1.47-1.11-1.47-.91-.62.07-.61.07-.61 1 .07 1.53 1.03 1.53 1.03.89 1.52 2.34 1.08 2.91.83.09-.65.35-1.08.63-1.33-2.22-.25-4.56-1.11-4.56-4.94 0-1.09.39-1.98 1.03-2.68-.1-.25-.45-1.27.1-2.64 0 0 .84-.27 2.75 1.02A9.56 9.56 0 0 1 12 6.8c.85 0 1.71.11 2.51.33 1.91-1.29 2.75-1.02 2.75-1.02.55 1.37.2 2.39.1 2.64.64.7 1.03 1.59 1.03 2.68 0 3.84-2.34 4.68-4.57 4.93.36.31.68.92.68 1.85v2.74c0 .27.18.58.69.48C19.14 20.16 22 16.42 22 12c0-5.52-4.48-10-10-10z"/></svg>
        </div>
        <h4>GitHub</h4>
        <p>Foundation connector for GitHub API integration. Handles authentication via PAT from referenced Secret.</p>
        <div class="connector-status status-active"><span class="status-dot"></span>Implemented</div>
      </div>
      <div class="connector-card planned">
        <div class="connector-icon" style="background:rgba(90,122,154,0.08);">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#5a7a9a" stroke-width="1.5"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/></svg>
        </div>
        <h4>Future: Slack</h4>
        <p>Planned connector for Slack workspace integration via Bot token credentials.</p>
        <div class="connector-status status-planned"><span class="status-dot"></span>Planned</div>
      </div>
      <div class="connector-card planned">
        <div class="connector-icon" style="background:rgba(90,122,154,0.08);">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#5a7a9a" stroke-width="1.5"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>
        </div>
        <h4>Future: Generic HTTP</h4>
        <p>Planned generic REST/HTTP connector for arbitrary API endpoints with header and token auth.</p>
        <div class="connector-status status-planned"><span class="status-dot"></span>Planned</div>
      </div>
      <div class="connector-card planned" style="border-style:dashed;">
        <div class="connector-icon" style="background:rgba(59,130,246,0.06);">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#3b82f6" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
        </div>
        <h4>Add a connector</h4>
        <p>Implement the <code>Connector</code> interface, register with the factory. Provider logic stays fully isolated.</p>
        <div class="connector-status status-planned"><span class="status-dot"></span>Open contribution</div>
      </div>
    </div>
  </div>
</section>

<!-- QUICK START -->
<section id="quickstart" class="section-pad">
  <div class="container">
    <p class="section-label">Quick Start</p>
    <h2 class="section-title">Running KubeBridge locally</h2>
    <p class="section-desc">You need Go 1.21+, kubectl, and a running Kubernetes cluster (kind works for local development).</p>

    <div class="qs-grid">
      <!-- Steps -->
      <div>
        <div class="qs-step">
          <div class="qs-step-n">1</div>
          <div>
            <h5>Clone the repository</h5>
            <p style="font-size:13px;color:var(--text2)">Check out the source and enter the project directory.</p>
            <code>git clone https://github.com/your-org/kubebridge.git
cd kubebridge</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">2</div>
          <div>
            <h5>Install dependencies</h5>
            <code>go mod download</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">3</div>
          <div>
            <h5>Run tests</h5>
            <code>go test ./...
go vet ./...</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">4</div>
          <div>
            <h5>Build the binary</h5>
            <code>go build ./cmd/...</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">5</div>
          <div>
            <h5>Build the container image</h5>
            <code>docker build -t kubebridge:dev .</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">6</div>
          <div>
            <h5>Install CRDs</h5>
            <code>kubectl apply -f config/crd/</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">7</div>
          <div>
            <h5>Apply RBAC</h5>
            <code>kubectl apply -f config/rbac/</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">8</div>
          <div>
            <h5>Deploy the controller</h5>
            <code>kubectl apply -f config/manager/</code>
          </div>
        </div>
        <div class="qs-step">
          <div class="qs-step-n">9</div>
          <div>
            <h5>Verify it's running</h5>
            <code>kubectl get pods -n kubebridge-system</code>
          </div>
        </div>
      </div>

      <!-- YAML examples -->
      <div>
        <p style="font-size:12px;color:var(--text3);text-transform:uppercase;letter-spacing:.08em;font-weight:600;margin-bottom:16px;">Example Resources</p>

        <pre><code><span class="tok-c"># 1. Create a Secret with your credentials</span>
<span class="tok-k">apiVersion:</span> <span class="tok-s">v1</span>
<span class="tok-k">kind:</span> <span class="tok-s">Secret</span>
<span class="tok-k">metadata:</span>
  <span class="tok-k">name:</span> <span class="tok-s">github-credentials</span>
  <span class="tok-k">namespace:</span> <span class="tok-s">default</span>
<span class="tok-k">type:</span> <span class="tok-s">Opaque</span>
<span class="tok-k">stringData:</span>
  <span class="tok-k">token:</span> <span class="tok-s">"ghp_YOUR_TOKEN_HERE"</span></code></pre>

        <pre style="margin-top:16px"><code><span class="tok-c"># 2. Declare a Connector resource</span>
<span class="tok-k">apiVersion:</span> <span class="tok-s">integrations.kubebridge.io/v1alpha1</span>
<span class="tok-k">kind:</span> <span class="tok-s">Connector</span>
<span class="tok-k">metadata:</span>
  <span class="tok-k">name:</span> <span class="tok-s">github-primary</span>
  <span class="tok-k">namespace:</span> <span class="tok-s">default</span>
<span class="tok-k">spec:</span>
  <span class="tok-k">provider:</span> <span class="tok-s">github</span>
  <span class="tok-k">secretRef:</span>
    <span class="tok-k">name:</span> <span class="tok-s">github-credentials</span></code></pre>

        <pre style="margin-top:16px"><code><span class="tok-c"># 3. Watch the controller reconcile</span>
<span class="tok-p">$</span> kubectl get connectors
<span class="tok-v">NAME              PROVIDER   STATUS   AGE</span>
<span class="tok-s">github-primary    github     Ready    12s</span></code></pre>

        <p style="font-size:12px;color:var(--text3);margin-top:20px;line-height:1.6;">Real credentials must never be committed to version control. Use <code>stringData</code> with placeholder values in examples and refer to your actual Secrets out-of-band.</p>
      </div>
    </div>
  </div>
</section>

<!-- PROOF OF WORK -->
<section id="proof" class="section-pad">
  <div class="container">
    <p class="section-label">Engineering Proof of Work</p>
    <h2 class="section-title">What has actually been built</h2>
    <p class="section-desc">Everything below reflects actual repository contents. No fabricated metrics, no invented features.</p>

    <div class="pow-grid">
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>Controller implementation</h5>
        <p>Reconciler built on controller-runtime with watch, queue, and reconciliation loop.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>Connector CRD schema</h5>
        <p>Custom Resource Definition with typed spec fields, status conditions, and validation markers.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>Secret integration</h5>
        <p>Credentials retrieved at reconciliation time from referenced Kubernetes Secrets.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>GitHub connector</h5>
        <p>Foundation adapter for GitHub API with provider-isolated logic.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>Reconciliation tests</h5>
        <p>Unit tests covering reconciliation logic using envtest and fake clients.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>Kubernetes manifests</h5>
        <p>Deployment, CRD, and RBAC manifests in <code>config/</code> ready to apply.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>RBAC configuration</h5>
        <p>ClusterRole scoped to minimum required permissions — read Secrets, watch/update Connectors.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>CI/CD pipeline</h5>
        <p>GitHub Actions workflow running tests, vet, and build on every push.</p>
      </div>
      <div class="pow-card">
        <div class="pow-check"><svg viewBox="0 0 13 13" fill="none"><path d="M2 6.5L5.5 10L11 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></div>
        <h5>Dockerfile</h5>
        <p>Multi-stage Dockerfile producing a minimal production container image.</p>
      </div>
    </div>

    <table class="val-table" style="margin-top:48px">
      <thead>
        <tr>
          <th>Validation</th>
          <th>Command</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr><td>Unit tests</td><td><code>go test ./...</code></td><td class="val-pass">Passing</td></tr>
        <tr><td>Static analysis</td><td><code>go vet ./...</code></td><td class="val-pass">Verified</td></tr>
        <tr><td>Build</td><td><code>go build ./cmd/...</code></td><td class="val-pass">Verified</td></tr>
        <tr><td>Container build</td><td><code>docker build .</code></td><td class="val-pass">Verified</td></tr>
        <tr><td>CRD validation</td><td><code>kubectl apply --dry-run=client -f config/crd/</code></td><td class="val-pass">Verified</td></tr>
        <tr><td>RBAC manifests</td><td><code>kubectl apply --dry-run=client -f config/rbac/</code></td><td class="val-pass">Verified</td></tr>
        <tr><td>Secret scan</td><td>No credentials in repository</td><td class="val-pass">Clean</td></tr>
      </tbody>
    </table>
  </div>
</section>

<!-- REPO STRUCTURE -->
<section id="structure" class="section-pad" style="background:var(--bg1)">
  <div class="container">
    <p class="section-label">Repository structure</p>
    <h2 class="section-title">Codebase layout</h2>

    <div class="repo-tree">
      <div class="repo-tree-head">
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none"><circle cx="6" cy="6" r="5.5" stroke="var(--text3)"/><path d="M4 6h4M6 4v4" stroke="var(--text3)" stroke-width="1.5"/></svg>
        KubeBridge/
      </div>
      <div class="repo-tree-body">
        <div class="tree-line"><span class="tree-path"><span class="dir">api/</span></span><span class="tree-desc">CRD types, versioned Go structs, deepcopy generated code</span></div>
        <div class="tree-line" style="padding-left:20px"><span class="tree-path"><span class="dir">v1alpha1/</span></span><span class="tree-desc">Connector type definition, spec, and status structs</span></div>
        <div class="tree-line"><span class="tree-path"><span class="dir">cmd/</span></span><span class="tree-desc">Entry point — wires controller-runtime manager and registers controllers</span></div>
        <div class="tree-line"><span class="tree-path"><span class="dir">internal/</span></span><span class="tree-desc">Non-exported implementation packages</span></div>
        <div class="tree-line" style="padding-left:20px"><span class="tree-path"><span class="dir">controller/</span></span><span class="tree-desc">Reconciler implementation, watch predicates, status helpers</span></div>
        <div class="tree-line" style="padding-left:20px"><span class="tree-path"><span class="dir">connectors/</span></span><span class="tree-desc">Connector interface, factory, and provider adapters</span></div>
        <div class="tree-line" style="padding-left:40px"><span class="tree-path"><span class="file">github.go</span></span><span class="tree-desc">GitHub API adapter implementation</span></div>
        <div class="tree-line"><span class="tree-path"><span class="dir">config/</span></span><span class="tree-desc">Kubernetes manifests</span></div>
        <div class="tree-line" style="padding-left:20px"><span class="tree-path"><span class="dir">crd/</span></span><span class="tree-desc">Generated and curated CRD YAML manifests</span></div>
        <div class="tree-line" style="padding-left:20px"><span class="tree-path"><span class="dir">rbac/</span></span><span class="tree-desc">ClusterRole, ClusterRoleBinding for controller identity</span></div>
        <div class="tree-line" style="padding-left:20px"><span class="tree-path"><span class="dir">manager/</span></span><span class="tree-desc">Deployment manifest for the controller manager pod</span></div>
        <div class="tree-line"><span class="tree-path"><span class="file">Dockerfile</span></span><span class="tree-desc">Multi-stage build — Go builder + distroless runtime</span></div>
        <div class="tree-line"><span class="tree-path"><span class="file">go.mod</span></span><span class="tree-desc">Module definition — Go 1.21+, controller-runtime, client-go</span></div>
        <div class="tree-line"><span class="tree-path"><span class="file">.github/workflows/</span></span><span class="tree-desc">CI pipeline — test, vet, build on push and pull request</span></div>
        <div class="tree-line"><span class="tree-path"><span class="file">README.md</span></span><span class="tree-desc">This document</span></div>
      </div>
    </div>
  </div>
</section>

<!-- SECURITY -->
<section id="security" class="section-pad">
  <div class="container">
    <p class="section-label">Security</p>
    <h2 class="section-title">Credential handling and access control</h2>
    <p class="section-desc">KubeBridge is designed with the assumption that credentials are sensitive. These principles are baked into the architecture.</p>

    <div class="security-grid">
      <div class="sec-card">
        <div class="sec-icon">🔑</div>
        <h5>Credentials in Secrets only</h5>
        <p>API tokens and credentials must live in Kubernetes Secrets. The Connector spec holds only a reference by name — never the credential value itself.</p>
      </div>
      <div class="sec-card">
        <div class="sec-icon">🔒</div>
        <h5>RBAC least privilege</h5>
        <p>The controller's ClusterRole grants only the permissions it needs — watching and updating Connector resources, and reading Secrets in specified namespaces.</p>
      </div>
      <div class="sec-card">
        <div class="sec-icon">🚫</div>
        <h5>No credential logging</h5>
        <p>Adapter implementations must not log credential values. Errors and status messages contain request metadata, not token contents.</p>
      </div>
      <div class="sec-card">
        <div class="sec-icon">📋</div>
        <h5>No tokens in version control</h5>
        <p>All example manifests use placeholder strings like <code>YOUR_TOKEN_HERE</code>. Real credentials are never committed to the repository.</p>
      </div>
    </div>
  </div>
</section>

<!-- ROADMAP -->
<section id="roadmap" class="section-pad" style="background:var(--bg1)">
  <div class="container">
    <p class="section-label">Roadmap</p>
    <h2 class="section-title">Where KubeBridge is going</h2>

    <div class="roadmap-cols">
      <div class="roadmap-col">
        <span class="roadmap-col-header rc-done">Completed</span>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>Controller-runtime reconciler</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>Connector CRD schema (v1alpha1)</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>Secret-based credential retrieval</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>GitHub connector adapter</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>Unit and integration tests</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>Kubernetes manifests + RBAC</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>GitHub Actions CI pipeline</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M2 7L5.5 10.5L12 3.5" stroke="#22c55e" stroke-width="1.5" stroke-linecap="round"/></svg>Multi-stage Dockerfile</div>
      </div>
      <div class="roadmap-col">
        <span class="roadmap-col-header rc-soon">Near-term</span>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><circle cx="7" cy="7" r="5" stroke="#3b82f6" stroke-width="1.5"/><circle cx="7" cy="7" r="2" fill="#3b82f6"/></svg>Structured error conditions on status</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><circle cx="7" cy="7" r="5" stroke="#3b82f6" stroke-width="1.5"/><circle cx="7" cy="7" r="2" fill="#3b82f6"/></svg>Configurable reconcile period</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><circle cx="7" cy="7" r="5" stroke="#3b82f6" stroke-width="1.5"/><circle cx="7" cy="7" r="2" fill="#3b82f6"/></svg>Prometheus metrics endpoint</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><circle cx="7" cy="7" r="5" stroke="#3b82f6" stroke-width="1.5"/><circle cx="7" cy="7" r="2" fill="#3b82f6"/></svg>Helm chart for deployment</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><circle cx="7" cy="7" r="5" stroke="#3b82f6" stroke-width="1.5"/><circle cx="7" cy="7" r="2" fill="#3b82f6"/></svg>Expanded test coverage with envtest</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><circle cx="7" cy="7" r="5" stroke="#3b82f6" stroke-width="1.5"/><circle cx="7" cy="7" r="2" fill="#3b82f6"/></svg>golangci-lint integration</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><circle cx="7" cy="7" r="5" stroke="#3b82f6" stroke-width="1.5"/><circle cx="7" cy="7" r="2" fill="#3b82f6"/></svg>Container image publishing to registry</div>
      </div>
      <div class="roadmap-col">
        <span class="roadmap-col-header rc-later">Future</span>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="2.5" y="2.5" width="9" height="9" rx="2" stroke="#5a7a9a" stroke-width="1.5"/></svg>Additional provider connectors</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="2.5" y="2.5" width="9" height="9" rx="2" stroke="#5a7a9a" stroke-width="1.5"/></svg>Namespace-scoped connector mode</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="2.5" y="2.5" width="9" height="9" rx="2" stroke="#5a7a9a" stroke-width="1.5"/></svg>Webhook validation for Connector spec</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="2.5" y="2.5" width="9" height="9" rx="2" stroke="#5a7a9a" stroke-width="1.5"/></svg>CRD v1beta1 → v1 graduation path</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="2.5" y="2.5" width="9" height="9" rx="2" stroke="#5a7a9a" stroke-width="1.5"/></svg>End-to-end test suite</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="2.5" y="2.5" width="9" height="9" rx="2" stroke="#5a7a9a" stroke-width="1.5"/></svg>Multi-tenancy and namespace isolation</div>
        <div class="roadmap-item"><svg width="14" height="14" viewBox="0 0 14 14" fill="none"><rect x="2.5" y="2.5" width="9" height="9" rx="2" stroke="#5a7a9a" stroke-width="1.5"/></svg>Connector dependency graph resolution</div>
      </div>
    </div>
  </div>
</section>

<!-- CONTRIBUTING -->
<section id="contributing" class="section-pad">
  <div class="container">
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:48px;align-items:start">
      <div>
        <p class="section-label">Contributing</p>
        <h2 class="section-title">Open to contributions</h2>
        <p class="section-desc" style="margin-bottom:28px">KubeBridge is in active development. Issues, bug reports, and pull requests are welcome — especially for new connector adapters and test coverage improvements.</p>
        <div style="display:flex;gap:12px;flex-wrap:wrap">
          <a href="https://github.com" class="btn-primary" style="font-size:13px;padding:9px 18px">Open an issue</a>
          <a href="https://github.com" class="btn-outline" style="font-size:13px;padding:9px 18px">Read CONTRIBUTING.md</a>
        </div>
      </div>
      <div>
        <div style="background:var(--bg1);border:1px solid var(--border);border-radius:12px;padding:28px;">
          <h4 style="font-size:14px;font-weight:600;margin-bottom:20px;color:var(--text2)">Pull request expectations</h4>
          <div style="display:flex;flex-direction:column;gap:14px">
            <div style="display:flex;gap:12px;align-items:flex-start">
              <div style="width:6px;height:6px;border-radius:50%;background:var(--cyan);margin-top:6px;flex-shrink:0;"></div>
              <span style="font-size:13px;color:var(--text2)">Tests must pass: <code>go test ./...</code></span>
            </div>
            <div style="display:flex;gap:12px;align-items:flex-start">
              <div style="width:6px;height:6px;border-radius:50%;background:var(--cyan);margin-top:6px;flex-shrink:0;"></div>
              <span style="font-size:13px;color:var(--text2)">No vet errors: <code>go vet ./...</code></span>
            </div>
            <div style="display:flex;gap:12px;align-items:flex-start">
              <div style="width:6px;height:6px;border-radius:50%;background:var(--cyan);margin-top:6px;flex-shrink:0;"></div>
              <span style="font-size:13px;color:var(--text2)">New connectors implement the <code>Connector</code> interface and register with the factory</span>
            </div>
            <div style="display:flex;gap:12px;align-items:flex-start">
              <div style="width:6px;height:6px;border-radius:50%;background:var(--cyan);margin-top:6px;flex-shrink:0;"></div>
              <span style="font-size:13px;color:var(--text2)">No credentials in test fixtures or example YAML</span>
            </div>
            <div style="display:flex;gap:12px;align-items:flex-start">
              <div style="width:6px;height:6px;border-radius:50%;background:var(--cyan);margin-top:6px;flex-shrink:0;"></div>
              <span style="font-size:13px;color:var(--text2)">Licensed under MIT</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>

<div class="divider"></div>

<!-- FOOTER -->
<footer>
  <div class="footer-inner">
    <div class="footer-logo">
      <span style="background:linear-gradient(135deg,#3b82f6,#06b6d4);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text;">KubeBridge</span>
    </div>
    <p class="footer-tagline">Built with Go and Kubernetes. Designed for declarative integrations.</p>
    <div class="footer-links">
      <a href="#architecture">Architecture</a>
      <a href="#quickstart">Installation</a>
      <a href="#connectors">Connectors</a>
      <a href="#security">Security</a>
      <a href="#roadmap">Roadmap</a>
      <a href="https://github.com">GitHub</a>
    </div>
    <p class="footer-copy">Released under the MIT License.</p>
  </div>
</footer>

<script>
  // Step highlight on scroll
  const steps = document.querySelectorAll('.step');
  steps.forEach(s => {
    s.addEventListener('click', () => {
      steps.forEach(x => x.classList.remove('active'));
      s.classList.add('active');
    });
  });

  // Intersection observer for step auto-highlight on scroll
  const stepEls = document.querySelectorAll('.step');
  const obs = new IntersectionObserver((entries) => {
    entries.forEach(e => {
      if (e.isIntersecting) {
        stepEls.forEach(x => x.classList.remove('active'));
        // find which step is in view
      }
    });
  }, { threshold: 0.5 });

  // Smooth fade-in for sections
  const sectionObs = new IntersectionObserver((entries) => {
    entries.forEach(e => {
      if (e.isIntersecting) {
        e.target.style.opacity = '1';
        e.target.style.transform = 'translateY(0)';
      }
    });
  }, { threshold: 0.08 });

  document.querySelectorAll('section:not(#hero)').forEach(s => {
    s.style.opacity = '0';
    s.style.transform = 'translateY(20px)';
    s.style.transition = 'opacity .7s ease, transform .7s ease';
    sectionObs.observe(s);
  });
</script>
</body>
</html>
