import sys
import argparse
import os
from markitdown import MarkItDown

def main():
    parser = argparse.ArgumentParser(description='MarkItDown Shim')
    parser.add_argument('uri', help='The URI of the document to convert')
    parser.add_argument('--llm-model', help='The LLM model to use for vision/audio descriptions', default=None)
    parser.add_argument('--llm-base-url', help='The base URL for the LLM API (e.g. for Ollama)', default=None)
    parser.add_argument('--openai-key', help='API Key for the LLM (OpenAI or other)', default="no-key-required")

    args = parser.parse_args()

    # Check if local file exists if it's not a URL
    if not args.uri.startswith(("http://", "https://")) and not os.path.exists(args.uri):
        print(f"Error: File not found: {args.uri}", file=sys.stderr)
        sys.exit(1)

    try:
        md_args = {}
        if args.llm_model:
            try:
                from openai import OpenAI
                
                # Setup client with optional base_url (e.g. for Ollama)
                client_args = {"api_key": args.openai_key}
                if args.llm_base_url:
                    client_args["base_url"] = args.llm_base_url
                
                client = OpenAI(**client_args)
                md_args['llm_client'] = client
                md_args['llm_model'] = args.llm_model
            except ImportError:
                print("Warning: openai package not installed. Vision/Audio features might not work.", file=sys.stderr)

        md = MarkItDown(**md_args)
        result = md.convert(args.uri)
        # We output the result directly to stdout
        print(result.text_content)
    except Exception as e:
        print(f"Error converting document: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
