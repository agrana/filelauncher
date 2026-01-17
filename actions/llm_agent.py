
import os
import argparse
import sys
import time
from openai import OpenAI
from dotenv import load_dotenv

load_dotenv()

def get_prompt_for_output(output_type, content, base_path):
    """
    Reads the prompt from a file in the 'prompts' directory.
    Falls back to a default if the file is not found.
    """
    prompt_file = os.path.join(base_path, "prompts", f"{output_type}.txt")
    
    if os.path.exists(prompt_file):
        try:
            with open(prompt_file, 'r', encoding='utf-8') as f:
                template = f.read()
                return template.replace("{content}", content)
        except Exception as e:
            print(f"Warning: Failed to read prompt file {prompt_file}: {e}", file=sys.stderr)

    # Fallback default
    return f"You are a helpful assistant. Please process the following text for the '{output_type}' format:\n\n{content}"

def main():
    parser = argparse.ArgumentParser(description="LLM Agent for filelauncher")
    parser.add_argument("--input", required=True, help="Input file path")
    parser.add_argument("--outputs", required=True, help="Comma-separated list of output suffixes")
    parser.add_argument("--model", default="gpt-4o", help="OpenAI model to use (default: gpt-4o)")
    parser.add_argument("--temperature", type=float, default=0.7, help="Sampling temperature (default: 0.7)")
    parser.add_argument("--system-prompt", default="default", help="Path to system prompt file or 'default'")
    
    args = parser.parse_args()
    
    input_path = args.input
    outputs = args.outputs.split(",")
    model_name = args.model
    temperature = args.temperature
    system_prompt_arg = args.system_prompt
    
    if not os.path.exists(input_path):
        print(f"Error: Input file '{input_path}' not found.", file=sys.stderr)
        sys.exit(1)
        
    # Load system prompt
    system_prompt_content = "You are a helpful content creator assistant."
    if system_prompt_arg != "default":
        # Check if it's a file relative to project root or absolute
        # Assuming script_dir/project_root logic below needs to happen earlier or we use absolute path
        pass 
        
    try:
        with open(input_path, 'r', encoding='utf-8') as f:
            content = f.read()
            # Safety: Remove #publish marker from input to prevent propagation
            content = content.replace("#publish", "")
    except Exception as e:
        print(f"Error reading input file: {e}", file=sys.stderr)
        sys.exit(1)

    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        print("Error: OPENAI_API_KEY environment variable not set.", file=sys.stderr)
        sys.exit(1)

    # Determine the project root (assuming this script is in actions/)
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(script_dir)
    
    client = OpenAI(api_key=api_key)

    # Resolve system prompt fully
    if system_prompt_arg != "default":
        sys_prompt_path = system_prompt_arg
        if not os.path.isabs(sys_prompt_path):
             sys_prompt_path = os.path.join(project_root, sys_prompt_path)
             
        if os.path.exists(sys_prompt_path):
             try:
                 with open(sys_prompt_path, 'r', encoding='utf-8') as f:
                     system_prompt_content = f.read().strip()
             except Exception as e:
                 print(f"Warning: Failed to read system prompt file {sys_prompt_path}: {e}", file=sys.stderr)
        else:
             print(f"Warning: System prompt file {sys_prompt_path} not found. Using default.", file=sys.stderr)

    for output_type in outputs:
        output_type = output_type.strip()
        if not output_type:
            continue
            
        print(f"Processing output: {output_type} with model {model_name}, temp {temperature}...")
        
        prompt = get_prompt_for_output(output_type, content, project_root)
        
        try:
            start_time = time.time()
            response = client.chat.completions.create(
                model=model_name,
                temperature=temperature,
                messages=[
                    {"role": "system", "content": system_prompt_content},
                    {"role": "user", "content": prompt}
                ]
            )
            duration = time.time() - start_time
            
            generated_content = response.choices[0].message.content
            
            # Construct output filename: input.output_type.ext
            # e.g., note.md -> note.medium.md
            base, ext = os.path.splitext(input_path)
            output_filename = f"{base}.{output_type}{ext}"
            
            with open(output_filename, 'w', encoding='utf-8') as f:
                f.write(generated_content)
                
            print(f"Generated: {output_filename} (took {duration:.2f}s)")
            
        except Exception as e:
            print(f"Error generating content for {output_type}: {e}", file=sys.stderr)
            # We don't exit here, so other outputs might succeed

if __name__ == "__main__":
    main()
