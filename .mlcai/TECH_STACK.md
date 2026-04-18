# 🛠 Tech Stack \u0026 Constraints: MLC MarkItDown MCP Server

## Kern-Versionen
- **Sprache:** Go 1.24.2 (Server-Wrapper), Python 3.11 (Konvertierungs-Runtime).
- **Framework:** MCP (Model Context Protocol) via `mcp-go`.
- **Infrastruktur:** Docker-basiertes Deployment (Stage 1: Build Go, Stage 2: Python Runtime).

## Bibliotheken (Erlaubt/Fixiert)
- **Go:**
  - `github.com/mark3labs/mcp-go v0.45.0` (Core MCP Funktionalität).
  - `github.com/hmsoft0815/mlcartifact v0.4.0` (gRPC-Client für Artifact-Server).
  - `google.golang.org/grpc` (gRPC-Kommunikation).
- **Python:**
  - `markitdown` (Microsoft Bibliothek zur Dokumentenkonvertierung).
  - `ffmpeg`, `libmagic1`, `libxml2`, `libxslt1-dev` (System-Dependencies für MarkItDown).

## Einschränkungen (Constraints)
- **Dateigröße:** Maximal 250 Zeilen pro Go-Datei (Style Guide).
- **Dokumentation:** Jedes exportierte Struct und jede exportierte Funktion MUSS GoDoc-Kommentare haben.
- **Konvertierung:** Erfolgt über einen Python-Shim (`shim.py`), der von Go als Subprozess aufgerufen wird.
- **Keine direkten Datei-Writes:** Ergebnisse > 10KB werden automatisch an den Artifact-Server (via gRPC) delegiert.

## Styling-Regeln (Go)
- **Naming:** CamelCase für Variablen, PascalCase für exportierte Typen/Funktionen.
- **Error Handling:** Explizit, mit Context-Wrapping (`fmt.Errorf(\"context: %w\", err)`).
- **Guard Clauses:** Frühe Returns statt tiefer Verschachtelung werden bevorzugt.
- **Receiver:** Kurze Receiver-Namen (z.B. `h *ConvertHandler`).

---

## 📋 Meta

- **Zuletzt aktualisiert:** 2026-04-18
- **Aktualisiert von:** gemini-2.0-flash-exp
- **Status:** Initialversion
