#!/bin/bash

# create-project.sh - A helper script to create a SpringWell project and navigate to it

# Default values
PROJECT_NAME="my-project"
TEMPLATE="basic"
DB="h2"
CLI_PATH="./build/springwell"

# Display help
function show_help {
    echo "SpringWell Project Creator"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME       Project name (default: my-project)"
    echo "  -t, --template NAME   Project template (default: basic, options: basic, aws-temporal-auth0)"
    echo "  -d, --db TYPE         Database type (default: h2, options: h2, postgres, mysql)"
    echo "  -p, --path PATH       Path to the springwell CLI (default: ./build/springwell)"
    echo "  -h, --help            Show this help message"
    echo ""
    echo "Example:"
    echo "  $0 --name my-awesome-app --template aws-temporal-auth0 --db postgres"
    echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -n|--name)
            PROJECT_NAME="$2"
            shift # past argument
            shift # past value
            ;;
        -t|--template)
            TEMPLATE="$2"
            shift # past argument
            shift # past value
            ;;
        -d|--db)
            DB="$2"
            shift # past argument
            shift # past value
            ;;
        -p|--path)
            CLI_PATH="$2"
            shift # past argument
            shift # past value
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Check if the CLI exists
if [ ! -f "$CLI_PATH" ]; then
    echo "Error: SpringWell CLI not found at $CLI_PATH"
    echo "Please build the CLI first with 'make build' or specify the correct path with --path"
    exit 1
fi

# Create the project
echo "Creating new SpringWell project: $PROJECT_NAME"
echo "  Template: $TEMPLATE"
echo "  Database: $DB"
echo ""

"$CLI_PATH" new "$PROJECT_NAME" --template "$TEMPLATE" --db "$DB"

if [ $? -ne 0 ]; then
    echo "Error: Failed to create project"
    exit 1
fi

echo ""
echo "Project created successfully!"
echo "Navigating to project directory..."

# Navigate to the project directory
cd "$PROJECT_NAME" || { 
    echo "Error: Failed to navigate to project directory"
    exit 1
}

echo "You are now in the $PROJECT_NAME directory"
echo "Ready to start development!"

# Execute a new shell in the project directory
exec $SHELL 