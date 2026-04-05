# Servify Website

Static product landing page for Servify. The website should present Servify primarily as an independently deployed intelligent customer service product, not as an implementation status dashboard.

Structure:
- index.html            Landing page
- assets/css/style.css  Stylesheet
- assets/js/main.js     Small client-side interactions
- assets/img/           Images (logo, illustrations)

Local preview:
- Simple Python server:
  - python3 -m http.server -d apps/website 8081
  - Open http://localhost:8081

Production hosting:
- Any static hosting (Nginx/Apache/CDN/Pages)
- Set cache headers for /assets/**
