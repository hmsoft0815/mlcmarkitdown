import sys
from markitdown import MarkItDown
import os

def main():
    if len(sys.argv) < 2:
        print("Usage: python shim.py <uri>")
        sys.exit(1)

    uri = sys.argv[1]
    
    # Check if local file exists if it's not a URL
    if not uri.startswith(("http://", "https://")) and not os.path.exists(uri):
        print(f"Error: File not found: {uri}", file=sys.stderr)
        sys.exit(1)

    try:
        md = MarkItDown()
        result = md.convert(uri)
        # We output the result directly to stdout
        print(result.text_content)
    except Exception as e:
        print(f"Error converting document: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
