import os
import argparse
import sys
import requests
from dotenv import load_dotenv

load_dotenv()

def main():
    parser = argparse.ArgumentParser(description="Medium Publisher for filelauncher")
    parser.add_argument("--input", required=True, help="Input file path")
    args = parser.parse_args()

    file_path = args.input

    if not os.path.exists(file_path):
        print(f"Error: Input file '{file_path}' not found.", file=sys.stderr)
        sys.exit(1)

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        print(f"Error reading input file: {e}", file=sys.stderr)
        sys.exit(1)

    if "#publish" not in content:
        # Marker not found, skip silently or with strict logging
        # print(f"Skipping {file_path}: #publish marker not found.")
        return

    # Marker found, proceed to publish
    print(f"Marker found in {file_path}. Initiating Medium publish...")

    token = os.environ.get("MEDIUM_INTEGRATION_TOKEN")
    user_id = os.environ.get("MEDIUM_USER_ID")

    if not token or not user_id:
        print("Error: MEDIUM_INTEGRATION_TOKEN or MEDIUM_USER_ID not set in .env", file=sys.stderr)
        sys.exit(1)

    # API Logic
    url = f"https://api.medium.com/v1/users/{user_id}/posts"
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
        "Accept": "application/json",
    }

    # Simple title extraction (first line #)
    lines = content.split('\n')
    title = "Untitled"
    for line in lines:
        if line.startswith("# "):
            title = line.replace("# ", "").strip()
            break

    # Remove the #publish marker from content before posting
    cleaned_content = content.replace("#publish", "").strip()

    data = {
        "title": title,
        "contentFormat": "markdown",
        "content": cleaned_content,
        "publishStatus": "draft" # Default to draft for safety
    }

    try:
        response = requests.post(url, headers=headers, json=data)
        response.raise_for_status()
        print(f"Successfully published to Medium: {response.json().get('data', {}).get('url')}")
    except Exception as e:
        print(f"Failed to publish to Medium: {e}", file=sys.stderr)
        if hasattr(e, 'response') and e.response is not None:
             print(f"Response: {e.response.text}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
