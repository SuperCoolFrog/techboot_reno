# Techboot Reno

## Puzzles

Logical Gates Puzzles

Commands: `And` and `OR`

Example:

```
    [  END  ]
        .
        .
        .
        .
        .
  1:[   ?   ]
      .   |
      .   |
      .   |
      .   |
      .   |
      .   |
      .   |
    [ START ]
```

To complete user must enter command:

```
set 1 or=
```

Future:

- They can set multiple at once:

```
set 1 or, 2 and, 3 not=
```

- Some they would need to send current to fill "powerbar" then redirect to end 
- Maybe some looping


- Boss battle will hide in locations and can be "zapped" by current
- Boss Battle sometimes you have to avoid sending current

## Trealla prolog

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


