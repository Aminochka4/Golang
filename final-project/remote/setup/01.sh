#!/bin/bash
set -eu
# ==================================================================================== #
# VARIABLES
# ==================================================================================== #
# Set the timezone for the server. A full list of available timezones can be found by
# running timedatectl list-timezones.
TIMEZONE=Asia/Almaty
# Set the name of the new user to create.
USERNAME=aminochka
# Prompt to enter a password for the PostgreSQL user (rather than hard-coding
# a password in this script).
read -p "Enter password for postgres DB user: " DB_PASSWORD
# Force all output to be presented in en_US for the duration of this script. This avoids
# any "setting locale failed" errors while this script is running, before we have
# installed support for all locales. Do not change this setting!
export LC_ALL=en_US.UTF-8
# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #
# Enable the "universe" repository.
add-apt-repository --yes universe
# Update all software packages. Using the --force-confnew flag means that configuration
# files will be replaced if newer ones are available.
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade
# Set the system timezone and install all locales.
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all
# Add the new user (and give them sudo privileges).
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
# Force a password to be set for the new user the first time they log in.
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"
# Copy the SSH keys from the root user to the new user.
rsync --archive --chown=${USERNAME}:${USERNAME} /Victus/.ssh /home/${USERNAME}
# Configure the firewall to allow SSH, HTTP and HTTPS traffic.
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable
# Install fail2ban.
apt --yes install fail2ban
# Install the migrate CLI tool.
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate
# Install PostgreSQL.
apt --yes install postgresql
# Set up the DB and create a user account with the password entered earlier.
sudo -i -u postgres psql -c "CREATE DATABASE mydatabase"
sudo -i -u postgres psql -d mydatabase -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgres psql -d mydatabase -c "CREATE ROLE myuser WITH LOGIN PASSWORD '${DB_PASSWORD}'"
# Add a DSN for connecting to the database to the system-wide environment
# variables in the /etc/environment file.
echo "DB_DSN='postgres://myuser:${DB_PASSWORD}@localhost/mydatabase'" >> /etc/environment
# Install Caddy (see https://caddyserver.com/docs/install#debian-ubuntu-raspbian).
apt --yes install -y debian-keyring debian-archive-keyring apt-transport-https
curl -L https://dl.cloudsmith.io/public/caddy/stable/gpg.key | sudo apt-key add -
curl -L https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt | sudo tee -a /etc/apt/sources.list.d/caddy-stable.list
apt update
apt --yes install caddy
# Clone your project from GitHub
git clone https://github.com/Aminochka4/Golang/tree/main/final-project
cd final-project
# Install Go
wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
sudo tar -xvf go1.14.2.linux-amd64.tar.gz
sudo mv go /usr/local
# Set environment variables for Go
echo "export GOROOT=/usr/local/go" >> ~/.profile
echo "export GOPATH=$HOME/go" >> ~/.profile
echo "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH" >> ~/.profile
source ~/.profile
# Install dependencies
go mod download
# Build the project
go build
# Run the project
./cmd/my-project
echo "Script complete! Rebooting..."
reboot
