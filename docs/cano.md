# CanoKey Documentation

## Glossary

- **FIDO**: Fast Identity Online, an open industry alliance dedicated to setting and promoting passwordless authentication standards.
- **WebAuthn**: Web Authentication, an authentication protocol defined by W3C. Its device-side implementations include CTAP 2.1/2.0 protocols.
- **Passkey**: A Passkey is an implementation of WebAuthn. CanoKey can store Passkey keys through the CTAP protocol.
- **OpenPGP**: Open Pretty Good Privacy, an open standard protocol used for encrypting and signing data, widely used for securing communications such as emails. CanoKey supports all mandatory features in OpenPGP 3.4.1.
- **GnuPG**: Gnu Privacy Guard, a free software tool implementing the OpenPGP standard, used for encrypting and signing data.
- **PIV**: Personal Identity Verification, a U.S. government identity authentication standard designed to enhance personal identity verification security. CanoKey supports all mandatory and some extended features of the PIV standard.
- **NDEF**: NFC Data Exchange Format, a standard format for exchanging data between NFC devices. CanoKey supports emulation of NFC tags compliant with the NDEF format.
- **OTP**: One-Time Password, a single-use password used for one-time authentication, enhancing security and often used in two-factor authentication. Common OTPs include HOTP and TOTP. CanoKey supports storing keys for both OTP types.
- **HOTP**: HMAC-based One-Time Password, a one-time password generated based on the HMAC (Hash-based Message Authentication Code) algorithm.
- **TOTP**: Time-based One-Time Password, a dynamic one-time password based on time, generated using the current time and a shared key.
- **WebUSB**: A protocol that allows web pages to directly communicate with USB devices, simplifying device connection and data transfer. CanoKey supports configuration via WebUSB.

## Default PINs

During the use of CanoKey, users will encounter various PINs. This page lists all default passwords and descriptions, aiming to help users distinguish different PINs.

| PIN Name          |  Default Value   | Description                                                                                                                                                                                                                                                               |
| :---------------- | :--------------: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Admin PIN         |      123456      | Used to manage different applications on CanoKey, such as resetting applications, modifying NDEF configurations, etc.                                                                                                                                                     |
| FIDO2 PIN         | No default value | Some FIDO2 applications that enforce PIN usage will ask for this password. FIDO2 PIN has no preset value. Users will be prompted to set the FIDO2 PIN when they first use a FIDO2 application that forces PIN usage, at which point the user can set this PIN themselves. |
| OpenPGP PIN       |      123456      | Used for regular OpenPGP operations, such as OpenPGP signing, etc.                                                                                                                                                                                                        |
| OpenPGP Admin PIN |     12345678     | Used for management operations of the OpenPGP application. For example, it is required when generating an OpenPGP key pair on CanoKey or modifying OpenPGP key attributes.                                                                                                |
| PIV PIN           |      123456      | Used for regular PIV operations, such as PIV identity verification, signing through PKCS#11 call to PIV, etc.                                                                                                                                                             |
| PIV PUK           |     12345678     | Used to unlock the PIV PIN if the PIV PIN gets locked.                                                                                                                                                                                                                    |

## Setup

This section introduces the necessary settings required to use CanoKey.

### 1. Windows

It can be used without configuration.

### 2. macOS

#### 2.1 Big Sur and earlier

Please compile and install the latest version of the [ccid driver](https://ccid.apdu.fr/).

#### 2.2 Monterey and Ventura

It can be used without configuration.

#### 2.3 Sonoma and later

CanoKey Pigeon users need to compile and install the latest version of the [ccid driver](https://ccid.apdu.fr/).

CanoKey Canary users can use it without configuration.

#### 2.4 pcsc-tools

pcsc-tools provides the `pcsc_scan` command for quickly troubleshooting PC/SC driver issues.

```bash
brew install pcsc-lite
```

### 3. Linux

Linux users can perform the following configuration for easier use.

#### 3.1 udev

The aim of udev rules are to allow users without root privileges to use it. Please create `/etc/udev/rules.d/69-canokeys.rules` and fill in the following content.

```
# GnuPG/pcsclite
SUBSYSTEM!="usb", GOTO="canokeys_rules_end"
ACTION!="add|change", GOTO="canokeys_rules_end"
ATTRS{idVendor}=="20a0", ATTRS{idProduct}=="42d4", ENV{ID_SMARTCARD_READER}="1"
LABEL="canokeys_rules_end"

# FIDO2
# note that if you find this line in 70-fido.rules, you can ignore it
KERNEL=="hidraw*", SUBSYSTEM=="hidraw", ATTRS{idVendor}=="20a0", ATTRS{idProduct}=="42d4", TAG+="uaccess", GROUP="plugdev", MODE="0660"

# make this usb device accessible for users, used in WebUSB
# change the mode so unprivileged users can access it, insecure rule, though
SUBSYSTEMS=="usb", ATTR{idVendor}=="20a0", ATTR{idProduct}=="42d4", MODE:="0666"
# if the above works for WebUSB (web console), you may change into a more secure way
# choose one of the following rules
# note if you use "plugdev", make sure you have this group and the wanted user is in that group
#SUBSYSTEMS=="usb", ATTR{idVendor}=="20a0", ATTR{idProduct}=="42d4", GROUP="plugdev", MODE="0660"
#SUBSYSTEMS=="usb", ATTR{idVendor}=="20a0", ATTR{idProduct}=="42d4", TAG+="uaccess"
```

`TAG+="uaccess"` is more inclined to systemd, while `GROUP="plugdev", MODE="0660"` is more traditional. You can choose either solution.

After adding this file, run the following command to apply the changes.

```
udevadm control --reload-rules && udevadm trigger
```

#### 3.2 CCID

Please use ccid version 1.4.34 or later.

#### 3.3 pcsc-tools

pcsc-tools provides the `pcsc_scan` command for quickly troubleshooting PC/SC driver issues.

## Settings

### 1. Main Functions

The settings of CanoKey include important operations for managing and configuring the device. Currently, you can manage CanoKey through the following methods:

- Web Console (Supports only Chrome and Chromium-based browsers): <https://console.canokeys.org>
- [iOS App](https://apps.apple.com/app/canokey-console/id6476454147)
- [Android App](https://play.google.com/store/apps/details?id=org.canokeys.console)

CanoKey Pigeon only supports NFC access in the iOS App.

The main configuration contents of the settings include:

- LED Status: Enabled by default
- NDEF (NFC Data Exchange Format): Enabled by default
- WebUSB Login Page: Enabled by default
- Reset

### 2. Reset

#### 2.1 Reset Application Individually

If you know the Admin PIN, you can reset each application individually through the console.

**The data of the reset application will be erased**

#### 2.2 Reset CanoKey

If you do not remember the Admin PIN, you can reset CanoKey when the Admin PIN is completely locked out.
In the settings application of the console, click "Reset", and the LED indicator of CanoKey will flash. When it flashes, please touch the key; repeat the above operation until it stops flashing.

**All data will be erased**.

## WebAuthn (Passkey)

### 1. Features

CanoKey's WebAuthn functionality adheres to the [CTAP 2.1](https://fidoalliance.org/specs/fido-v2.1-ps-20210615/fido-client-to-authenticator-protocol-v2.1-ps-20210615.html) and [CTAP 2.0](https://fidoalliance.org/specs/fido-v2.0-ps-20190130/fido-client-to-authenticator-protocol-v2.0-ps-20190130.html).

Supported features include:

- Up to 64 sets of discoverable credentials (resident keys)
- HMAC extension support
- Ed25519 algorithm

From firmware version 2.0.0, CanoKey also supports the following features:

- Discoverable Credentials management
- PIN Protocol 2
- Credential Blob
- Large Blob

From firmware version 3.0.0, CanoKey experimentally supports the SM2 algorithm.

> [!note]
>
> CanoKey with firmware version 3.0.0 does not support U2F, WebAuthn on iOS 17.4 and 18 via USB, or WebAuthn on macOS (including Safari, Firefox, and applications relying on Apple's CTAP stack).

### 2. Primary Uses

#### 2.1 Multi-Factor Authentication

CanoKey can be used for two-factor authentication on many [websites](https://2fa.directory/int/).

> [!note]
>
> By default, CanoKey does not set a PIN. Some websites and certain features (such as Discoverable Credentials management) require you to set a PIN. Please set it when prompted.

#### 2.2 SSH

##### 2.2.1 OpenSSH Version Requirement

To use FIDO keys for SSH authentication, ensure that the installed OpenSSH version supports this feature. The minimum version requirements are as follows:

- OpenSSH 8.2 and above

You can check the current `ssh` client and `sshd` service versions using the following commands:

```bash
ssh -V
sshd -V
```

##### 2.2.2 Using ECDSA-SK and ED25519-SK Key Pairs

###### Creating ECDSA-SK and ED25519-SK Key Pairs

1. Create an ECDSA-SK key pair:

```bash
ssh-keygen -t ecdsa-sk -f ~/.ssh/id_ecdsa_sk
```

2. Create an ED25519-SK key pair:

```bash
ssh-keygen -t ed25519-sk -f ~/.ssh/id_ed25519_sk
```

###### Adding the Public Key to the Target Server

Add the content of the generated public key file (`~/.ssh/id_ecdsa_sk.pub` or `~/.ssh/id_ed25519_sk.pub`) to the `~/.ssh/authorized_keys` file on the target server.

You can use the following command to copy the public key to the remote server:

```bash
ssh-copy-id -i ~/.ssh/id_ecdsa_sk.pub username@remote_host
```

Or

```bash
ssh-copy-id -i ~/.ssh/id_ed25519_sk.pub username@remote_host
```

###### Using on Other Machines

Copy the generated private key file (`~/.ssh/id_ecdsa_sk` or `~/.ssh/id_ed25519_sk`) to other machines where it needs to be used. Ensure the correct file permissions:

```bash
chmod 600 ~/.ssh/id_ecdsa_sk
chmod 600 ~/.ssh/id_ed25519_sk
```

##### 2.2.3 Using Discoverable Credential (Resident Key)

> [!note]
>
> CanoKey firmware version must be at least 2.0.0.

###### Creating RK Keys

1. Create an ECDSA-SK RK:

```bash
ssh-keygen -t ecdsa-sk -O resident -f ~/.ssh/id_ecdsa_sk
```

2. Create an ED25519-SK RK:

```bash
ssh-keygen -t ed25519-sk -O resident -f ~/.ssh/id_ed25519_sk
```

###### Adding the Public Key to the Target Server

Similar to the non-RK keys above, add the content of the generated public key file to the server's `~/.ssh/authorized_keys` file.

```bash
ssh-copy-id -i ~/.ssh/id_ecdsa_sk.pub username@remote_host
```

Or

```bash
ssh-copy-id -i ~/.ssh/id_ed25519_sk.pub username@remote_host
```

#### 2.3 PAM

Please refer to [pam-u2f](https://developers.yubico.com/pam-u2f/).

#### 2.4 HMAC-secret Extension

- [systemd-cryptenroll](http://0pointer.net/blog/unlocking-luks2-volumes-with-tpm2-fido2-pkcs11-security-hardware-on-systemd-248.html), used for LUKS full-disk encryption

> [!note]
>
> Due to a [bug](https://github.com/Yubico/libfido2/issues/322#issuecomment-817174671) in the CTAP implementation, CanoKey firmware version ≤ 1.3 is incompatible with libfido2 1.7.0, and thus cannot be used with `systemd-cryptenroll`. Affected users should use libfido2 1.6.0.

## OpenPGP

[OpenPGP](https://www.openpgp.org/) is a signature and encryption standard specified by [RFC4880](https://tools.ietf.org/html/rfc4880). This standard achieves information and file signing/encryption through the use of private keys. One of the commonly used OpenPGP tools is GNU Privacy Guard, often abbreviated as GnuPG or GPG. In Windows, you can also use [Kleopatra](https://www.openpgp.org/software/kleopatra/).

### 1. Basic Information

#### 1.1 Supported Algorithms

- RSA2048
- RSA3072
- RSA4096
- X25519
- Ed25519
- NIST P-256 (secp256r1, prime256v1)
- NIST P-384 (secp384r1)
- secp256k1

#### 1.2 Defaults

- PIN: Default is 123456, minimum length is 6, maximum length is 64
- Admin PIN: Default is 12345678, minimum length is 8, maximum length is 64
- Reset Code: Default is empty, minimum length is 8, maximum length is 64
- Signature PIN: forced (PIN verification required for each signature)
- Touch Policy: SIG, DEC, AUT are off
- Touch Cache Time: 0

> [!note]
>
> Firmware versions 1.6.1 and earlier only support RSA public keys with e = 65537.
> Firmware versions 2.0.0 and above support RSA3072 / RSA4096 key generation.

#### 1.3 Touch Policy

> [!note]
>
> Touch policy is only effective when using the USB interface.

OpenPGP supports up to 3 keys: signature key (SIG), encryption key (DEC), and authentication key (AUT). Depending on the firmware version, you can set the touch policy for SIG, DEC, and AUT in the CanoKey Console or via the `gpg` command. The value of touch cache time ranges from 0 to 255 seconds (0 means no cache).

##### Firmware Version <= 1.4

Please use the “Settings” application in the CanoKey Console to modify the touch policy.

##### Firmware Version >= 1.5

Please use GnuPG to modify the touch policy.

#### 1.4 PIN Policy

For DEC and AUT keys, after the PIN verification is successful, verification will not be required again until CanoKey is disconnected and reinserted.

For SIG, if `forcesig` is on, a PIN is required for each signature; otherwise, a PIN is only required for the first signature after power-on.

### 2. Common Operations

Please refer to the [GNU Privacy Handbook](https://gnupg.org/gph/en/manual.html).

### 3. FAQs

#### 3.1 GnuPG and PC/SC Conflict

GnuPG, by default, uses its own implementation ([scdaemon](https://www.gnupg.org/documentation/manuals/gnupg/Invoking-SCDAEMON.html)) to access smart cards including CanoKey, which conflicts with PC/SC. For details, see: <https://ludovicrousseau.blogspot.com/2019/06/gnupg-and-pcsc-conflicts.html>.

To avoid conflicts, we recommend using the PC/SC interface to access CanoKey by adding the following to `scdaemon.conf`:

```
disable-ccid
```

In Linux and macOS, this file is usually located at `~/.gnupg/scdaemon.conf`.

In Windows, this issue is typically not encountered. If necessary, please modify the `scdaemon.conf` file under the GnuPG installation directory.

#### 3.2 PC/SC Occupancy

Since PC/SC access to smart cards may be exclusive (depending on the application access mode), even if configured correctly, GnuPG may still fail to access CanoKey. If you encounter this issue, simply re-plug CanoKey.

Programs commonly occupying PC/SC include:

- Firefox: You can unload "OpenSC Smartcard framework" in “Preferences > Privacy & Security > Certificates”.

## PIV

PIV (Personal Identity Verification) is defined by the US federal government [FIPS 201](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.201-2.pdf) standard. PIV can store keys and certificates for signing and encryption, enabling functions such as digital signatures and file encryption.

### 1. Basic Information

#### 1.1 Supported Algorithms

- RSA2048
- NIST P-256
- NIST P-384

Starting from CanoKey Canary, the following extended algorithms are also supported:

| Algorithm Name | Algorithm ID |
| :------------- | :----------- |
| RSA3072        | 05           |
| RSA4096        | 16           |
| secp256k1      | 53           |
| Ed25519        | E0           |
| X25519         | E1           |
| SM2            | 54           |

> [!note]
>
> CanoKey firmware version 3.0.0 only supports signing 32-byte data using the Ed25519 algorithm and only supports using internally generated X25519 keys.

#### 1.2 Default Values

- PIN: 123456
- PUK: 12345678
- Management Key: `010203040506070801020304050607080102030405060708`

#### 1.3 Key Slots

CanoKey supports the following key slots:

- 9A: PIV Authentication
- 9E: Card Authentication
- 9C: Digital Signature
- 9D: Key Management

Starting from firmware version 2.0.0, CanoKey also supports the following key slots:

- 82, 83

#### 1.4 PIN and Touch Policies

##### PIN Policy

- Never: Never verify PIN
- Always: Verify PIN for every use
- Once: Verify PIN once per session

##### Touch Policy

- Never: Never require touch
- Always: Require touch for every use
- Cached: No touch required if touched within the last 15 seconds, otherwise touch is required

##### Default Policies

> [!note]
>
> Starting from firmware version 2.0.0, CanoKey supports configuring PIV PIN and touch policies.

| Key Slot | Default PIN Policy | Default Touch Policy |
| :------- | :----------------- | :------------------- |
| 9E       | Never              | Never                |
| Others   | Once               | Never                |

#### 1.5 Data Size Limitations

- Certificate:
  - Firmware version 1.5 or earlier: 1000 bytes
  - Firmware version 1.6 or later: 3000 bytes
- Card Capability Container: 287 bytes
- Card Holder Unique Identifier: 2916 bytes
- Printed Information: 245 bytes

#### 1.6 Other Features

Starting from firmware version 2.0.0, CanoKey supports viewing PIV metadata.

### 2. Common Operations

> [!note]
>
> As PIV is typically issued by system administrators and used by regular users, please review the documentation to understand the following content before proceeding.

#### 2.1 Tools

It is recommended to use [yubico-piv-tool](https://developers.yubico.com/yubico-piv-tool/Releases/) for related operations.

#### 2.2 Importing Keys and Certificates Separately

If the key and certificate are in two separate files, they need to be imported separately.

Importing the private key:

```bash
yubico-piv-tool -r canokey -a import-key -s 9a -i private-key.pem
```

Importing the certificate:

```bash
yubico-piv-tool -r canokey -a import-certificate -s 9a -i certificate.pem
```

Here, `-s 9a` indicates using the 9A key slot, which can be changed as needed.

#### 2.3 Importing PKCS#12 File

To import a PKCS#12 file (.p12 or .pfx) containing both the private key and certificate, execute:

```bash
yubico-piv-tool -r canokey -a import-key -a import-certificate -K PKCS12 -s 9a -i certificate.p12
```

#### 2.4 Generating Key and Self-signing

Generate a new private key and self-sign it:

```bash
yubico-piv-tool -r canokey -a generate -s 9a -A RSA2048 -o public-key.pem
yubico-piv-tool -r canokey -a verify-pin -a selfsign -s 9a -S "/CN=Test Certificate" -i public-key.pem -o certificate.pem
yubico-piv-tool -r canokey -a import-certificate -s 9a -i certificate.pem
```

#### 2.5 Additional Steps for Windows

Since Windows caches certificate information based on CHUID, you need to update the CHUID after certificate import on Windows:

```bash
yubico-piv-tool -r canokey -a set-chuid
```

## NDEF (NFC Tag)

NDEF (NFC Data Exchange Format) is a lightweight, extensible binary data format that typically contains data record formats such as URLs and text.

### Default Values

- Mode: Default is read/write mode, which can be modified in the console
- Content: Default is URL, with value "https://canokeys.org"
- Maximum length: 1022 bytes

### How to Use

1. Download and install the NFC Tools application.
2. Open the NFC Tools application, select the "Read" option to read existing NDEF tag content.
3. Select the "Write" option to write new content to the NDEF tag. Note that NDEF tags can be rewritten.
4. Bring the top of your iPhone or the back of your Android device close to the CanoKey to complete the read/write operation.

Please note that the NDEF function has no encryption protection, so the transmitted information is stored and transmitted in plain text. Please use it cautiously according to your security requirements and avoid storing and transmitting sensitive information.

## OTP

OATH is [an organization](https://openauthentication.org/) who provides open authentication standards: Time-based One Time Password (TOTP) and HMAC-based One Time Password (HOTP).

HOTP and TOTP are both implemented in CanoKey Pigeon and Canokey epoxy editions. CanoKey can hold up to 100 OATH tokens.

### Firmware version 1.5 and newer (For example, CanoKey Pigeon).

You should use `ykman` command version 4.0 or above to configure OATH and read OATH token.

#### Setting up

If your authentication provider provides you with a URI `otpauth://totp/username@EXAMPLE.COM:12345678-90ab-cdef-1234-567890abcdef?digits=6&secret=SOMESECRET&period=30&algorithm=SHA1&issuer=username%40EXAMPLE.COM`, you should configure it by

```
ykman -r "Canokeys" oath accounts uri "otpauth://totp/username@EXAMPLE.COM:12345678-90ab-cdef-1234-567890abcdef?digits=6&secret=SOMESECRET&period=30&algorithm=SHA1&issuer=username%40EXAMPLE.COM"
```

Or if you are offered the fields separately, including the base32 encoded secret (in this example, it's `SOMESECRET`), the algorithm (in this example, it's SHA1), the length of digits, you can setup by

```
ykman -r "Canokeys" oath accounts add -o TOTP -d 6 -a SHA1 -i 'username@EXAMPLE.COM' -P 30 USERNAME SOMESECRET
```

You can use `ykman oath accounts add --help` to find all the available options.

#### Read OATH token

You can use [Yubico Authenticator](https://www.yubico.com/products/yubico-authenticator/) to read the OTP.

- Open Yubico Authenticator, click the button on the top left corner to fire up the side menu.
- Click 'Settings' and then click 'Custom reader'
- Select 'Enable custom reader' and fill in 'Canokey' in the 'Custom reader filter'
- Click 'Save'
- Click the '<' on the top to go back to Settings.
- Unplug and plug in your CanoKey
- Click the top left button to switch to 'Authenticator'

Then you'll see all your OATH tokens listed.

### Firmware version older than 1.5 (for example, CanoKey epoxy edition)

You should use [CanoKey Web Console](https://console.canokeys.org) on a Chromium-based web browser to configure OATH.

#### Setup

- Go to [OATH Applet](https://console.canokeys.org/oath) of the web console and connect your CanoKey.
- Click 'CONNECT' on the top right corner
- Select your CanoKey from the prompt dialog
- Click the '+' sign on the down left of the box, then the 'Add credential to OATH Applet' section will show up on the web page.
- Fill in the details, **or** copy the URI you got and click 'IMPORT OTPAUTH FROM CLIPBOARD'.
- Click ADD
- If there is no error, there will be a green box with 'Add OATH credential success' shown in the bottom left of your web page.

#### Read OATH TOTP token

- Go to [OATH Applet](https://console.canokeys.org/oath) of the web console and connect your CanoKey.
- Click 'CONNECT' on the top right corner
- Select your CanoKey from the prompt dialog
- Click the 'Calculate TOTP' button on the line of TOTP you want to calculate.
- If there is no error, there will be a green box with 'TOTP code is xxxxxx' shown in the bottom left of your web page.

### Optionally: Enable touch to input for HOTP

If you are using HOTP and you want to make your CanoKey input your HOTP token every time you touch the key, you should

- Go to [Admin Applet](https://console.canokeys.org/admin) of the web console and connect your CanoKey.
- If you haven't connected your CanoKey to the web console, click 'CONNECT' on the top right corner and select your CanoKey from the prompt dialog.
- Click AUTHENTICATE and input your admin applet password to authenticate as admin user
- Enable 'HOTP on touch' in the Config section, then you'll see a green box with 'HOTP on touch is on' shown in the bottom left of the web page.
- Go to [OATH Applet](https://console.canokeys.org/oath) of the web console
- Click the star icon '٭' on the line of HOTP you want to use. This will make the HOTP token default.
- Click DISCONNECT on the top right corner
- Unplug and plug in your CanoKey again.
  Now you can press the touch area to input your default HOTP token.
