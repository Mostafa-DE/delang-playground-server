#!/bin/bash

install_pre_dependencies() {
    sudo apt-get update
    sudo apt-get install -y curl unzip jq
}

install_go() {
    wget https://dl.google.com/go/go1.18.linux-amd64.tar.gz
    sudo tar -xvf go1.18.linux-amd64.tar.gz
    sudo mv go /usr/local
    rm go1.18.linux-amd64.tar.gz

    echo "export GOROOT=/usr/local/go" >> ~/.profile
    echo "export PATH=\$PATH:\$GOROOT/bin" >> ~/.profile
}

install_de() {
    # Fetch the latest tag from GitHub
    latest_tag=$(curl -s https://api.github.com/repos/Mostafa-DE/delang/tags | jq -r '.[0].name')

    # Remove the 'v' from the tag if present
    formatted_tag=${latest_tag#v}

    download_url="https://github.com/Mostafa-DE/delang/archive/refs/tags/$latest_tag.zip"

    curl -L $download_url -o "delang-$latest_tag.zip"
    unzip "delang-$latest_tag.zip" -d delang-latest
    rm "delang-$latest_tag.zip"

    # Move the contents to ~/delang
    mkdir -p "$HOME/delang"
    mv delang-latest/delang-$formatted_tag/* "$HOME/delang/"
    rm -rf delang-latest

}

shell_profile=""
update_shell_profile() {
    if [ -f "$HOME/.zshrc" ]; then
        shell_profile="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        shell_profile="$HOME/.bashrc"
    else
        echo "Unsupported shell."
        exit 1
    fi

    # Append the de() function to the shell profile
    cat << 'EOF' >> "$shell_profile"

de() {
    local curr_dir="$PWD"
    if [ -z "$1" ]; then
        # No file name provided, enter the REPL
        cd ~/delang/ && go run main.go && cd "$curr_dir"
    else
        local file_name="$1"
        if [[ $file_name == /* ]]; then
            local full_file_path="$file_name"
        else
            local file_dir="$PWD"
            local full_file_path="$file_dir/$file_name"
        fi
        cd ~/delang/ && go run main.go "$full_file_path" && cd "$curr_dir"
    fi
}
EOF
}

# Check if Go is installed, if not, install it
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go..."
    install_go
fi

if [ -d "$HOME/delang" ]; then
    echo "Delang is already installed in your machine under $HOME/delang"
    exit 1
else
    echo "Installing Pre dependencies"
    install_pre_dependencies

    echo "Installing DE..."
    install_de

    echo "Updating Shell Profile"
    update_shell_profile

fi

rm de_install.sh

echo -e "\033[32mInstallation complete. Please restart your terminal or run 'source $shell_profile' to apply changes.\033[0m"

