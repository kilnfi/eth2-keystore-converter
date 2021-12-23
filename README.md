# eth2-keystore-converter

Converts Eth2 EIP-2335 scrypt keystores to pbkdf2 keystores (and vice-versa).

## Usage

Converting a scrypt keystore to pbkdf2 using stdin/stdout

```
cat keystore-scrypt.json \
  | eth2-keystore-converter -p $(cat keystore-scrypt.txt) \
  > keystore-pbkdf2.json
```

Converting a pbkdf2 keystore to scrypt using options

```
eth2-keystore-converter \
  -f keystore-pbkdf2.json \
  -p $(cat keystore-pbkdf2.txt) \
  -c scrypt \
  -o keystore-scrypt.json
```

## Motivation

The scrypt formatted keystores load slower and require a lot more memory than pbkdf2 formatted keystores which load almost instantly.

Thus, converting on the fly all keystores to pbkdf2 makes it easier to operate eth2 validators at scale and reduce computing resources.

## Credits

This utility is mostly a wrapper for the [go-eth2-wallet-encryptor-keystorev4](https://github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4]) library.
