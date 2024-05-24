#!/bin/bash

install_pre_dependencies() {
    sudo apt-get update
    sudo apt-get install -y curl unzip jq
}

install_go() {
    GO_VERSION="1.22.3"
    wget https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz
    sudo tar -xvf go$GO_VERSION.linux-amd64.tar.gz
    sudo mv go /usr/local
    rm go$GO_VERSION.linux-amd64.tar.gz

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
    cat << 'EOF' >> "$HOME/de"
#!/bin/bash
de() {
    local delang_dir="$HOME/delang"
    local curr_dir="$PWD"

    finish() {
        cd "$curr_dir"
    }

    trap finish EXIT

    if [ -z "$1" ]; then
        # No file name provided, enter the REPL
        cd "$delang_dir"
        go run main.go
    else
        local file_name="$1"

        if [[ $file_name != /* ]]; then
            file_name="$PWD/$file_name"
        fi

        cd "$delang_dir"
        go run main.go "$file_name"
    fi

    finish
}
de "$@"
EOF

chmod +x "$HOME/de"
sudo mv "$HOME/de" /usr/local/bin/
}

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

# Check if Go is installed, if not, install it
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go..."
    install_go
fi

rm de_install.sh

echo -e "\033[32mInstallation complete. Please restart your terminal or run 'source $shell_profile' to apply changes.\033[0m"

