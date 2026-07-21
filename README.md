# Techboot Reno


Trealla prolog

```
# 1. Install standard system build tools and development libraries
sudo apt update
sudo apt install build-essential git libedit-dev libffi-dev libssl-dev -y

# 2. Clone the official Trealla source code repository
git clone https://github.com/trealla-prolog/trealla.git
cd trealla

# 3. Compile the binary and run tests to verify local environment safety
make
make test

# 4. Install globally to your system path (/usr/local/bin/tpl)
sudo make install
```

Run:
```
tpl
```
