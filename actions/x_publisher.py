import os
import argparse
import sys
import requests
from requests_oauthlib import OAuth1
from dotenv import load_dotenv

load_dotenv()

def main():
    parser = argparse.ArgumentParser(description="X Publisher for filelauncher")
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
        return

    print(f"Marker found in {file_path}. Initiating X publish...")

    api_key = os.environ.get("TWITTER_API_KEY")
    api_secret = os.environ.get("TWITTER_API_SECRET")
    access_token = os.environ.get("TWITTER_ACCESS_TOKEN")
    access_token_secret = os.environ.get("TWITTER_ACCESS_SECRET")

    if not all([api_key, api_secret, access_token, access_token_secret]):
        print("Error: TWITTER credentials not fully set in .env", file=sys.stderr)
        sys.exit(1)

    # API Logic
    url = "https://api.twitter.com/2/tweets"
    auth = OAuth1(api_key, api_secret, access_token, access_token_secret)

    # Simple tweet content: remove marker, take first 280 chars (naive)
    cleaned_content = content.replace("#publish", "").strip()
    tweet_text = cleaned_content[:280]

    payload = {"text": tweet_text}

    try:
        response = requests.post(url, auth=auth, json=payload)
        response.raise_for_status()
        print(f"Successfully published to X: {response.json()}")
    except Exception as e:
        print(f"Failed to publish to X: {e}", file=sys.stderr)
        if hasattr(e, 'response') and e.response is not None:
             print(f"Response: {e.response.text}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
