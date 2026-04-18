# 🤖 AI \u0026 Integration Context: MLC MarkItDown MCP Server

## 1. Identität \u0026 Zweck
- **Kernaufgabe:** Robuster, performanter MCP-Server zur Konvertierung verschiedener Dokumentformate (PDF, Word, Excel, PowerPoint, HTML, CSV, Bilder, Audio) in Markdown.
- **Technischer Stack:** Go (Server-Wrapper), Python (MarkItDown-Bibliothek via Shim), gRPC (Integration mit Artifact-Server).
- **Hoster/Infrastruktur:** Docker (empfohlen), Python 3.11-slim Runtime, Go-Binary.

## 2. Die \"Nachbarschaft\" (System-Kontext)
- **Upstream (Wovon hänge ich ab?):**
  - **Microsoft MarkItDown (Python)** -\u003e https://github.com/microsoft/markitdown -\u003e Basis-Bibliothek für die Konvertierung.
  - **mlcartifact (gRPC)** -\u003e github.com/hmsoft0815/mlcartifact -\u003e Optionaler Speicher für große Konvertierungsergebnisse (> 10KB).
  - **LLM Provider (OpenAI / Ollama)** -\u003e Für Vision-basierte Bildbeschreibungen und Audio-Transkriptionen.
- **Downstream (Wer nutzt mich?):**
  - **MCP Clients (z.B. Claude, Gemini, IDE-Plugins)** -\u003e Nutzen die Tools zur Dokumentenverarbeitung in LLM-Workflows.
- **Shared Resources:**
  - **Artifact Store:** Nutzt den gleichen gRPC-Endpunkt wie andere MLC-Services zur persistenten Speicherung von Dokumenten.

## 3. Schnittstellen-Vertrag
- **Primäre API:** MCP (Model Context Protocol) via `stdio`, `sse` oder `streamable` HTTP.
- **Auth-Mechanismus:** Konfiguration via Umgebungsvariablen (z.B. `DEFAULT_LLM_AUTH_KEY` für OpenAI).
- **Wichtige Datenmodelle:**
  - `markitdown__convert__mlc`: Konvertiert URI zu Markdown.
  - `markitdown__convert_artifact__mlc`: Konvertiert existierende Artefakte.
  - `markitdown__quick_inspect__mlc`: Metadaten-Extraktion ohne volle Konvertierung.
- **API-Doku-Link:** Siehe `README.md` (Abschnitt \"Tools\") für detaillierte Argument-Listen und JSON-Strukturen.

## 4. Leitplanken \u0026 Regeln
- **Naming:** Go-Konventionen (PascalCase für Exporte, camelCase für interne Variablen).
- **Testing:** Unit-Tests via `go test ./...`, Integrationstests via `mcp-tester` (siehe `TESTING.md`).
- **Sicherheit:** Sensible API-Keys niemals in den Code, nur via Environment-Variablen injizieren.

## 5. Aktueller Fokus (Status)
- **Besonderheit:** Smart Artifact Integration (automatisches Speichern großer Outputs als Artefakt).
- **Transport:** Volle Unterstützung für `stdio` und `sse`.
- **Nächste Schritte:** Stabilisierung der Ollama-Vision-Integration.

---

## 📋 Meta

- **Zuletzt aktualisiert:** 2026-04-18
- **Aktualisiert von:** gemini-2.0-flash-exp
- **Status:** Initialversion
