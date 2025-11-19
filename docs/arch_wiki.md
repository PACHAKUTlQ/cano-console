# Arch Wiki

The [YubiKey](https://www.yubico.com) is a small USB Security token. Depending on the model, it can:

- Act as a smartcard (using the CCID protocol) - allowing storage of both [PGP](https://developers.yubico.com/PGP/) and [PIV](https://developers.yubico.com/PIV/) secret keys.
- Handle Universal 2nd Factor (U2F) requests.
- Store and query approximately 30 Initiative for Open Authentication (OATH) credentials.
- Handle challenge-response requests, in either the Yubico OTP mode or the HMAC-SHA1 mode.
- Generate One-time passwords (OTP) - [Yubico's AES based standard](https://developers.yubico.com/OTP/).
- "Type" a static password up to 63 characters.

While offering many features, newer versions of the YubiKey are [not released as open source](https://www.yubico.com/blog/secure-hardware-vs-open-source/). Alternatives are the Solo, TKey or Nitrokey.

## Installation

### Management tools

- **[YubiKey Manager](https://developers.yubico.com/yubikey-manager/)**: Python library and command-line tool (`ykman`) for configuring and querying a YubiKey over USB. Has optional GUI. (Package: `yubikey-manager`, `yubikey-manager-qt` (AUR))

  > [!note]
  >
  > After installation, enable `pcscd.service`, alternatively it might be useful to enable `pcscd.socket`.

  > [!warning]
  >
  > `yubikey-manager-qt` (AUR) is [completely broken](https://github.com/Yubico/yubikey-manager-qt/issues/361) on recent Arch/python/yubikey-manager versions, and Yubico [have deprecated this tool](https://github.com/Yubico/yubikey-manager-qt/commit/28dc02d11b081683b59d16d12043aaa3c0a6c75f).

- **[YubiKey Personalization](https://developers.yubico.com/yubikey-personalization/)**: Library and tool for configuring and querying a YubiKey over the OTP USB connection. More powerful than ykman, but harder to use. Has optional GUI. (Package: `yubikey-personalization`, `yubikey-personalization-gui`)

  > [!warning]
  >
  > Yubico officially announced the End of Life of the YubiKey Personalization package on February 19, 2025. It will not updated after February 19, 2026

### Authentication tools

- **[Yubico PAM](https://developers.yubico.com/yubico-pam/)**: PAM user authentication with either Yubico OTP or challenge-response. (Package: `yubico-pam`. The https://github.com/Yubico/yubico-pam GitHub repository is archived since February 2025. (Dead link))
- **[Yubico PAM-U2F](https://developers.yubico.com/pam-u2f/)**: PAM user authentication with U2F. (Package: `pam-u2f`)
- **[Yubico Authenticator for Desktop](https://developers.yubico.com/OATH/YubiKey_OATH_software.html)**: GUI to read OATH codes from your YubiKey over USB. Support the newer OATH implementation (YubiKey NEO and 4) as well as the older slot-based implementation (YubiKey Standard and Edge). Note Issue: `archlinux/packaging/packages/yubioath-desktop` #2. (Package: `yubioath-desktop` (AUR) (package not found))
- **[Yubico Authenticator 6.0+ for Desktop](https://developers.yubico.com/yubioath-flutter/)**: Yubico Authenticator 6.0+ (Version 6.0 and higher) is an application for managing the YubiKey's second factor credentials. Starting with Version 6.0, the codebase has been completely rewritten using the Flutter framework. (Package: `yubico-authenticator-bin` (AUR), `yubico-authenticator` (AUR))
- **[libfido2](https://developers.yubico.com/libfido2/)**: Client-side U2F support. Enables web browsers to use the U2F protocol for authentication with your YubiKey. (Package: `libfido2`)
- **[YubiKey Full Disk Encryption](https://github.com/agherzan/yubikey-full-disk-encryption)**: Use challenge-response mode to create strong LUKS passphrases. Supports full disk encryption. (Package: `yubikey-full-disk-encryption`)

## Inputs

The YubiKey takes inputs in the form of API calls over USB and button presses.

The button is very sensitive. Depending on the context, touching it does one of these things:

- Trigger a static password or one-time password (OTP) (Short press for slot 1, long press for slot 2). This is the default behavior, and easy to trigger inadvertently.
- Confirm / allow a function or access. The LED will illuminate to prompt the user.
- Insert / eject the smartcard

## Outputs

The YubiKey transforms these inputs into outputs:

- Keystrokes (emulating a USB keyboard), used to type static passwords and OTPs. (Note that static passwords are vulnerable to keyloggers.)
- The built-in LED:
  - Blinks once when plugged in, useful for troubleshooting.
  - Blinks steadily when a button press is required to permit an API response.
- API responses over USB. This is used for:
  - Challenge-Response requests (calculated using either Yubico OTP mode or HMAC-SHA1 mode)
  - U2F Challenge-Response requests
  - CCID Smartcard related requests

## USB connection modes

Depending on the YubiKey model, the device provides up to three different USB interfaces. Two of the interfaces implement the USB HID (Human Interface Device) device class; the third is a smart card interface (CCID). All three can be enabled or disabled independently, allowing control of their associated protocols.

The following table shows which protocols use which interfaces:

| Protocol | Interface    |
| -------- | ------------ |
| OTP      | Keyboard HID |
| FIDO     | Other HID    |
| PIV      | CCID         |
| OpenPGP  | CCID         |
| OATH     | CCID         |

`ykman` uses the term "modes", named OTP, FIDO, and CCID.

> [!note]
>
> ykman renamed "U2F" to "FIDO" in release 0.6.1 (2018-04-16). https://developers.yubico.com/yubikey-manager/Release_Notes.html

### Get enabled modes

For YubiKey prior to version 5:

```bash
$ ykman config mode
Current connection mode is: OTP+FIDO+CCID
```

> [!note]
>
> The command `ykman mode` has been deprecated and may be removed later.

For YubiKey version 5:

```bash
$ ykman info
Device type: YubiKey 5 NFC
Serial number: XXXXXXXXX
Firmware version: 5.4.3
Form factor: Keychain (USB-A)
Enabled USB interfaces: OTP, FIDO, CCID
NFC transport is enabled.

Applications    USB     NFC
FIDO2           Enabled Enabled
OTP             Enabled Enabled
FIDO U2F        Enabled Enabled
OATH            Enabled Enabled
YubiHSM Auth    Enabled Enabled
OpenPGP         Enabled Enabled
PIV             Enabled Enabled
```

### Set modes

All modes are enabled from the factory. To change them:

```bash
$ ykman mode [OPTIONS] MODE
```

- `MODE` can be a string, such as `OTP+FIDO+CCID`, or a shortened form `o+f+c`.
- `MODE` can be a mode-number, which encodes several enabled modes.

Here is a table of mode-numbers, if you care to use them:

| Mode numbers | Devices                             |
| ------------ | ----------------------------------- |
| 0            | OTP device only.                    |
| 1            | CCID device only.                   |
| 2            | OTP/CCID composite device.          |
| 3            | U2F device only.                    |
| 4            | OTP/U2F composite device.           |
| 5            | U2F/CCID composite device.          |
| 6            | OTP/U2F/CCID composite device.      |
| 81           | CCID device only, with touch-eject. |

> [!note]
>
> Some examples use mode number 86, which is [invalid](https://github.com/Yubico/yubikey-manager/issues/20#issuecomment-326496204). The 80 will be ignored, and it will behave like 6.

Options:

- `--touch-eject` - The button will insert and eject the smart card. This only works if the mode is CCID only; FIDO and OTP must be disabled.
- `--autoeject-timeout SECONDS` - Automatically eject the smart card after some time. Same restrictions as `--touch-eject`.
- `--chalresp-timeout SECONDS` - Set the challenge-response timeout.

For more information, see `ykman mode --help`.

## One-time password

This feature has a somewhat misleading name, because it also encompasses the static password and challenge-response functions.

2 slots are provided for this feature, accessible by short and long button presses respectively. Each can be configured with **one** of the following:

- Yubico OTP
- OATH-HOTP
- OATH-TOTP
- Challenge-response
- Static Password

Each function has several configuration options provided at the time of creation, but once set they cannot be read back. It is possible to swap slots 1 and 2, with `ykman otp swap`.

### Factory configuration

On a new YubiKey, Yubico OTP is preconfigured on slot 1. This initial AES symmetric key is stored in the YubiKey and on the Yubico Authentication server. This allows validating against YubiCloud, allowing the use of Yubico OTP in combination with the Yubico Forum website for instance or on https://demo.yubico.com).

> [!warning]
>
> If you ever overwrite the factory key in slot 1, you cannot create a new key of the same trust level. Factory generated Yubico OTP credentials begin with a `CC` prefix, while user generated credentials begin with `VV`. There is no fundamental difference in security or functionality, though some services only trust `CC` credentials. More information can be found in this [forum thread](https://forum.yubico.com/viewtopic12ca.html?f%3D16&t%3D1960).

### Yubico OTP

The [Yubico OTP](https://developers.yubico.com/OTP/) is based on symmetric cryptography. More specifically, each YubiKey contains a 128-bit AES key unique to that device, which is also stored on a validation server. When asked for a password, the YubiKey will create a token by concatenating different fields such as the ID of the key, a counter, and a random number, and encrypting the result.

This OTP is sent to the target system, which passes it to a validation server. The validation server (also in posession of the secret key) decrypts it and verifies the information inside. The result is returned to the target system, which can then decide whether to grant access.

#### YubiCloud and validation servers

Yubico provides a validation server with free unlimited access, called YubiCloud. YubiCloud knows the factory configuration of all YubiKeys, and is the "default" validation service used by (for example) `yubico-pam`. Yubico also offers [open-source implementations] (https://developers.yubico.com/Software_Projects/Yubico_OTP/YubiCloud_Validation_Servers/) (Dead link) of the server.

> [!note]
>
> To authenticate the Yubico validation server, you can use:
>
> - **HMAC**: use https://upgrade.yubico.com/getapikey/ to get an HMAC key and ID
> - **HTTPS**: the validation server's certificate is signed by GoDaddy, and is thus trusted by default in Arch installs (at least if you have `ca-certificates` installed)

#### Configuration and usage

Generate a new key in slot 2, and upload it to YubiCloud (opens in a browser):

```bash
$ ykman otp yubiotp --generate-key --upload 2
```

For more information, see `ykman otp yubiotp --help`.

#### Security risks

##### AES key compromise

As you can imagine, the AES key should be kept secret. It cannot be retrieved from the YubiKey itself (or it should not, at least not with software). It is also present in the validation server, so the security of this server is very important.

##### Validation requests/responses tampering

Since the target system relies on a validation server, a possible attack would be to impersonate it. To prevent this, the target system needs to authenticate the validation server, either using HMAC or HTTPS.

### Challenge-response

A challenge is sent to the YubiKey, which calculates a response based on some secret. The same challenge always results in the same response. Without the secret this calculation is not feasible, even with many challenge-response pairs.

This can be used for

- True 2-factor authentication: The user is provided a challenge, they must provide the correct response in addition to a password. Both parties must have the secret key.
- "Semi" 2-factor authentication: the challenge acts as a password, and the server stores the correct response. This is not an OTP, and if anyone can obtain the response they will gain access, but it is simpler as the server does not need the secret key.

There are two Challenge-Response algorithms:

- HMAC-SHA1
- Yubico OTP

You can set them up with a GUI using the `yubikey-personalization-gui`, or with the following instructions:

#### HMAC-SHA1 algorithm

Set up slot 2 in challenge response mode with a generated key:

```bash
$ ykman otp chalresp --generate 2
```

You can omit the `--generate` flag in order to provide a key, see `ykman otp chalresp --help`. A main advantage of providing a key is that it can be used to setup a second device as a backup. The command `openssl rand -hex 20` generates a suitable key, for example.

#### Yubico OTP algorithm

`ykman` Does not appear to support setting the chal-yubico algorithm, but you can use `ykpersonalize`. Generate a random key in slot 2:

```bash
$ ykpersonalize -2 -ochal-resp -ochal-yubico
```

For more information, see `ykpersonalize(1)`.

#### Sending a challenge

To send a challenge and get a response, the `ykchalresp -slot challenge` command can be used. For example,

```bash
$ ykchalresp -2 archie
12a19763be77d75af46fb76f0b737c117fa47205
```

returns a 40-byte SHA1-hash unique to the programmed slot 2. A different challenge produces another unique response.

### Static password

You can either generate a static password:

```bash
$ ykman otp static --generate slot
```

or provide one:

```bash
$ ykman otp static slot password
```

You have several options; you can set the length and character set of the generated password, and whether or not to send an Enter keystroke. See `ykman otp static --help` for more.

> [!tip]
>
> Most YubiKeys provide only limited (2) slots to save static passwords. A setup challenge-response slot provides static hash responses for unlimited challenges. While these are numeric and lower alphabet only, the response length provides considerable entropy.

### Emulated USB keyboard limitations, or "Why does my password look so weak?"

In order for the YubiKey to work with most keyboard layouts, passwords are by default limited to the ModHex alphabet (`cbdefghijklnrtuv`), digits `0-9`, and `!`. These characters use the same scan codes across a very large number of keyboard layouts, ensuring compatibility with most computers.

Yubico has provided a [whitepaper](https://resources.yubico.com/53ZDUYE6/as/9hccqgx9bwwqq96mhkk8jb4h/Static_Password_Function.pdf) on the subject.

## OATH

The YubiKey offers 2 OATH implementations:

**OATH API**: Newer method, can store approximately 30 credentials depending on the model. (YubiKey 4, NEO, and newer)
**OTP slot**: Older method, both OTP slots can store a single credential. (All models which support challenge-response)

### OATH API

If you prefer a GUI, you can use `yubioath-desktop` (AUR) (package not found).

`ykman` can add codes in the URI format with `ykman oath uri`. Here is a one-liner that will add a credential from an image of a QR code:

```bash
$ zbarimg qr_code.png --quiet --raw | xargs ykman oath accounts uri
```

You can also do things manually. Program a TOTP key, requiring a button touch to generate a code:

```bash
$ ykman oath accounts add --touch name secret
```

Program an HOTP key:

```bash
$ ykman oath accounts add --oath-type HOTP name secret
```

List credentials:

```bash
$ ykman oath accounts list
```

Generate codes:

```bash
$ ykman oath accounts code query
```

To see all available subcommands see `ykman oath --help`. To see information about each, use `ykman oath subcommand --help`.

### OTP slot implementation

Program an HOTP in slot 2:

```bash
$ ykman otp hotp 2 key
```

Program a TOTP:

```bash
$ ykman otp chalresp --totp slot key
```

Generate an HOTP:

```bash
$ ykman otp calculate slot
```

Generate a TOTP:

```bash
$ ykman otp calculate --totp slot
```

See also: `ykman otp --help` and https://developers.yubico.com/OATH/

## U2F

Universal 2nd Factor (U2F) with a YubiKey is very simple, requiring no configuration for the key itself. Note that this mode is also referred to as 'FIDO' in some documentation and utilities. You have a few limited management options through the `ykman` utility:

- Set a PIN: `ykman fido access change-pin`
- delete individual credentials: `ykman fido credentials delete QUERY`
- Reset all credentials and PIN: `ykman fido reset`

To use U2F for authentication, see the instructions in U2F.

Also see WebAuthn.

## CCID smartcard

CCID (Chip Card Interface Device) is a USB standard device class for use by USB devices that act as smart card readers or with security tokens that connect directly via USB, like the YubiKey. HID (Human Interface Device) and CCID are both USB device classes, i.e. they are in the same category of USB specifications. HID is a specification for computer peripherals, like keyboards. The YubiKey works like a USB (HID) keyboard when used in the OTP and FIDO modes, but switches to the CCID protocol when using the PIV application, or as an OpenPGP device.

CCID mode should be enabled by default on all YubiKeys shipped since November 2015 [https://www.yubico.com/support/knowledge-base/categories/articles/use-yubikey-yubikey-windows-hello-app/] (Dead link). Enable at least the CCID mode. Please see Get enabled modes.

### PIV

Starting with the YubiKey NEO, the YubiKeys contain a PIV (Personal Identity Verification) application on the chip. PIV is a US government standard (FIPS 201) that specifies how a token using RSA or ECC (Elliptic Curve Cryptography) is used for personal electronic identification. The YubiKey NEO only supports RSA encryption, later models (YubiKey 4 and 5) support both RSA and ECC. The exact algorithms supported depends on the firmware. For example, only YubiKeys with firmware 5.7 and up support RSA 3072, RSA 4096, Ed25519, and X25519 keys [https://developers.yubico.com/PIV/Introduction/YubiKey_and_PIV.html]. The distinguishing characteristic of a PIV token is that it is built to protect private keys and operate on-chip. A private key never leaves the token after it has been installed on it. Optionally, the private key can even be generated on-chip with the aid of an on-chip random number generator. If generated on-chip, the private key is never handled outside of the chip, and there is no way to recover it from the token. When using the PIV mechanism, the YubiKey functions as a CCID device.

### OpenPGP smartcards

The YubiKey can act as a standard OpenPGP smartcard; see GnuPG#Smartcards for instructions on how to set up and use it with GnuPG. Yubico also provides some documentation in https://developers.yubico.com/PGP/.

If you do not want to use the other features (U2F and OTP), the button can be configured to insert and eject it, and an auto-eject timeout can be set as well. See USB connection modes for more.

The default user pin is `123456` and the default admin pin is `12345678`. The default PUK is also `12345678`. Remember to change all 3.

## Use cases

This section details how to use your YubiKey for various authentication purposes. It is by no means an exhaustive list.

### Full disk encryption with LUKS

You have several options:

- **Challenge-Response:** the response to some challenge is used as a LUKS key. The challenge can act as a password for true 2-factor authentication, or stored in plain-text for one-factor authentication.
- **GnuPG:** Uses the yubikey's PGP smartcard functionality. Offers strong 2-factor authentication without needing a huge passphrase.
- **FIDO HMAC Secret:** If your YubiKey supports U2F, it can be configured to return a symmetric secret.

> [!note]
>
> A disk's encryption is only as strong as its weakest key. Once you configure one of these tools, consider removing the initial passphrase, or replacing it with a very long one.

#### Common prerequisites

- A bootable LUKS encrypted system, using the `encrypt` mkinitcpio hook, with at least one free keyslot.
  - With the exception of `mkinitcpio-ykfde` (AUR), the `sd-encrypt` hook is not supported by any of these tools.
- Backed up LUKS header (Optional, though advisable)

#### Challenge-response

See `yubikey-full-disk-encryption`'s [official documentation](https://github.com/agherzan/yubikey-full-disk-encryption#usage) for complete instructions. Broadly:

1. Install `yubikey-full-disk-encryption`.
2. Configure `/etc/ykfde.conf`.
3. Enroll the disk: `# ykfde-enroll -d /dev/DISK -s LUKS_SLOT`
4. Add the `ykfde` mkinitcpio hook before the `encrypt` hook.
5. Regenerate the initramfs.

   > [!note]
   >
   > Plymouth users: replace the `plymouth-encrypt` hook with the `ykfde` hook.

There are a few variations available:

- 2FA: default behavior. You must provide the challenge as a password when enrolling the device, and upon boot.
- 1FA: Set `YKFDE_CHALLENGE` in `ykfde.conf`. Note that this is stored in plaintext. Consider disabling non-root read permissions to this file.
- [NFC support](https://github.com/agherzan/yubikey-full-disk-encryption#enable-nfc-support-in-ykfde-initramfs-hook-experimental) (Experimental)
- [Suspend & Resume support](https://github.com/agherzan/yubikey-full-disk-encryption#enable-ykfde-suspend-service-experimental) (Experimental) Automatically lock encrypted volumes on suspend, unlock them on resume.

You must regenerate the initramfs for any configuration changes to take effect.

#### systemd-based initramfs

Users of the `sd-encrypt` hook may install `mkinitcpio-ykfde` (AUR) or `mkinitcpio-ykfde-git` (AUR) and follow the instruction in the [project documentation](https://github.com/eworm-de/mkinitcpio-ykfde/blob/master/README-mkinitcpio.md). The procedure is broadly similar to `yubikey-full-disk-encryption`.

#### GnuPG encrypted keyfile

One tool to accomplish this is [initramfs-scencrypt](https://github.com/fuhry/initramfs-scencrypt); see its docs for complete instructions. Note that as of October 2022 this package is not in the AUR and is not thoroughly tested, though the GitHub repository offers a PKGBUILD.

The dm-crypt pages offer a few alternatives, though they are mostly links to old forum posts.

#### HMAC secret extension of FIDO2 protocol

Yet another way of using YubiKey for full disk encryption is to utilize HMAC Secret Extension to retrieve the LUKS password from YubiKey. This can be protected by a passphrase. This functionality requires at least [YubiKey 5 with firmware 5.2.3+](https://support.yubico.com/hc/en-us/articles/360016649319-YubiKey-5-2-3-Enhancements-to-FIDO-2-Support).
For a passphrase protected solution, install `khefin` (AUR) and follow instructions available in [project documentation](https://github.com/mjec/khefin/wiki/Arch-Linux-LUKS-configuration).
For single factor (optionally PIN-protected) solution and starting with systemd 248, it is possible to use your FIDO2 key as LUKS2 keyslot. Instructions available in the [author's blog post](https://0pointer.net/blog/unlocking-luks2-volumes-with-tpm2-fido2-pkcs11-security-hardware-on-systemd-248.html).

### KeePass

KeePass can be configured for YubiKey support; see the YubiKey section for instructions.

### SSH keys

#### CCID

If your YubiKey supports CCID smartcards, you can use it as a hardware-backed SSH key, either based on GPG or PIV keys. Yubico offers good documentation:

- An [overview of both](https://developers.yubico.com/PIV/Guides/Securing_SSH_with_OpenPGP_or_PIV.html) possibilities, giving their advantages and disadvantages
- Instructions for [PGP authentication](https://developers.yubico.com/PGP/SSH_authentication/index.html)
- Instructions for [PIV authentication through user certificates](https://developers.yubico.com/PIV/Guides/SSH_user_certificates.html)
- Instructions for [PIV authentication through #PKCS11](https://developers.yubico.com/PIV/Guides/SSH_with_PIV_and_PKCS11.html)

> [!note]
>
> The default PIN code of the PIV application on the YubiKey is `123456`; you may want to change it, as well as the default management key. See the [device setup instructions](https://developers.yubico.com/PIV/Guides/Device_setup.html) for more.

#### U2F

You may also use the U2F feature of the YubiKey to create hardware-backed SSH keys. See SSH keys#FIDO/U2F for instructions.

#### PIV

`yubikey-agent` (AUR) stores the SSH key as PIV token. See [https://github.com/FiloSottile/yubikey-agent#readme] for a setup guide.

### Linux user authentication with PAM

PAM, and therefore anything which uses PAM for user authentication, can be configured to use a YubiKey as a factor of its user authentication process. This includes sudo, su, ssh, screen lockers, display managers, and nearly every other instance where a Linux system needs to authenticate a user. Its flexible configuration allows you to set whichever authentication requirements fit your needs, for the entire system, a specific application, or for groups of applications. For example, you could accept the YubiKey as an alternative to a password for local sessions, while requiring both for remote sessions. In addition to the Arch Wiki, You are encouraged to read `pam(8)` and `pam.conf(5)` to understand how it works and how to configure it.

There are several modules available which integrate YubiKey-supported protocols into PAM:

- `pam-u2f` - Supports U2F via the FIDO2 standard. If you are not sure which method to use, this one is a good choice.
  - Universal 2nd Factor#Authentication for user sessions
  - [Yubico's official docs](https://developers.yubico.com/pam-u2f/), including a list of supported module parameters.
  - Man Pages: `pam_u2f(8)`, `pamu2fcfg(1)`
- `oath-toolkit` - Supports OATH one-time passwords (either HOTP or TOTP)
  - pam_oath
- `yubico-pam` - Supports Yubico OTP and challenge-response OTPs. Note that Yubico OTP mode requires a network connection to a validation server, while challenge-response mode does not.
  - [Yubico's official docs](https://developers.yubico.com/yubico-pam/) (Dead link)
  - `pam_yubico(8)` - Take note of the `mode` parameter, used to set challenge-response mode.

> [!warning]
>
> Exercise caution when modifying PAM configuration files. Mistakes can render a system completely insecure, or so secure that no authentication is possible.

PAM configuration is beyond the scope of this article, but for a brief overview:

- Create file(s) containing authorized keys, either in users' home directories or centrally.
- Add a line in the appropriate place in the appropriate PAM configuration file which follows this format:
  ```
  	auth [required|sufficient] [module_name].so [module arguments]
  ```
- `auth required` for multifactor, `auth sufficient` for single factor.
- `module_name` - Example: `pam_u2f.so`. See a list of installed modules: `ls /usr/lib/security`
- Module configuration arguments are for things like the location of the keyfile, or which method the module should use for authentication.

#### SSH notes

- Yubico has provided [additional guidance](https://developers.yubico.com/yubico-pam/Yubikey_and_SSH_via_PAM.html) (Dead link). It is written for an old version of Ubuntu, but much of it still applies to an updated Arch system.
- If you are configuring a distant server to use YubiKey, you should open at least one additional, rescue SSH session, so that you are not locked out if the configuration fails.
- Check that `/etc/ssh/sshd_config` contains the following settings. The `sshd_config` shipped with `openssh` has these set correctly by default.
  ```
  ChallengeResponseAuthentication no
  UsePAM yes
  ```

### Browser/web integration

Many web services are beginning to support FIDO hardware tokens. See the U2F and WebAuthn pages for more information, but usually the only thing you need to do is to install `libfido2` and [try it](https://demo.yubico.com/webauthn-technical/registration).

## Tips and tricks

### Executing actions on insertion/removal of YubiKey device

For example, you want to perform an action when you pull your YubiKey out of the USB slot, create `/etc/udev/rules.d/80-yubikey-actions.rules` and add the following contents:

```
ACTION=="remove", ENV{ID_VENDOR}=="Yubico", ENV{ID_VENDOR_ID}=="1050", ENV{ID_MODEL_ID}=="0010|0111|0112|0113|0114|0115|0116|0401|0402|0403|0404|0405|0406|0407|0410", RUN+="/usr/local/bin/script args"
```

Please note, most keys are covered within this example but it may not work for all versions of YubiKey. You will have to look at the output of lsusb to get the vendor and model ID's, along with the description of the device or you could use udevadm to get information. Of course, to execute a script on insertion, you would change the action to 'add' instead of remove.

### Start Yubico Authenticator on insertion

The authenticator is a long-running GUI process. If run directly in a udev rule, the process would block udev's processing. If forked, udev would unconditionally kill the process after the event handling finishes. Thus you cannot start the authenticator from udev rules. However, systemd.device may be used to handle this case.

Similar to above, create `/etc/udev/rules.d/80-yubikey-actions.rules` and add the following contents:

```
ENV{ID_VENDOR}=="Yubico", ENV{ID_VENDOR_ID}=="1050", ENV{ID_MODEL_ID}=="0010|0111|0112|0113|0114|0115|0116|0401|0402|0403|0404|0405|0406|0407|0410", SYMLINK+="yubikey", TAG+="systemd"
```

Then create a new systemd user unit:

```ini
# ~/.config/systemd/user/yubioath-desktop.service
[Unit]
Description=Autostart Yubico Authenticator
# Uncomment if you want to stop the authenticator when unplugged.
#StopPropagatedFrom=dev-yubikey.device

[Install]
WantedBy=dev-yubikey.device

[Service]
Type=oneshot
ExecStart=/usr/bin/yubioath-desktop
```

and enable it. _systemctl_ would warn that it is added as a dependency to a non-existent unit `dev-yubikey.device`. But it is okay. Such unit will start existing once the YubiKey is plugged in.

## Maintenance / upgrades

### Installing the OATH Applet for a YubiKey NEO

These steps will allow you to install the OATH applet onto your YubiKey NEO. This allows the use of Yubico Authenticator in the Google Play Store.

> [!note]
>
> These steps are only for NEOs with a firmware version <= 3.1.2. The current generation NEOs (with U2F) come with the OpenPGP applet already installed)

#### Configure the NEO as a CCID device

1. Install `yubikey-personalization-gui` (`yubikey-personalization-gui-git` (AUR) (package not found)).
2. Add the udev rules and reboot so you can manage the YubiKey without needing to be root
3. Run `ykpersonalize -m82`, enter `y`, and hit enter.

#### Install the applet

1. Install `gpshell` (AUR), `gppcscconnectionplugin` (AUR), `globalplatform` (AUR), and `pcsclite`.
2. Start `pcscd.service`.
3. Download the most recent CAP file from the [ykneo-oath](https://developers.yubico.com/ykneo-oath/Releases/) site.
4. Download `gpinstall.txt` from [GitHub](https://github.com/Yubico/ykneo-oath/blob/master/gpinstall.txt).
5. Edit the line in gpinstall.txt beginning with `install -file` to reflect the path where the CAP file is located.
6. Open a terminal and run `gpshell path/to/gpinstall.txt`.
7. Ideally, a bunch of text will scroll by and it ends saying something like
   ```text
   Command --> 80E88013D7C000C400BE00C700CA00CA00B400BE00CE00D200D500D700B000DB00C700DF00BEFFFF00BE00E400AC00AE00AE00DB00E700A
   A00EA00ED00ED00ED00BE00EF00F100F400F100F700FA00FF00BE00F700AA01010103010700CA00C400B400AA00F700B400AA00B600C7010C
   010C00AA0140012001B0056810B0013005600000056810E0011006B4B44304B44404B44106B44B4405B443400343B002410636810E06B4B44
   407326810B004B43103441003334002B102B404B3B403BB4003B440076820A4100221024405B4341008B44600000231066820A15D848CB77
   27D0EDA00
   Response <-- 009000
   Command --> 80E60C002107A000000527210108A00000052721010108A000000527210101010003C901000000
   Wrapped command --> 84E60C002907A000000527210108A00000052721010108A000000527210101010003C9010000B4648127914A4C7C00
   Response <-- 009000
   card_disconnect
   release_context
   ```
8. Unplug the NEO and try it with the Yubico Authenticator app.

#### (Optional) Install the Yubico Authenticator desktop client

You can get the desktop version of the Yubico Authenticator by installing `yubioath-desktop` (AUR) (package not found).

While `pcscd.service` is running, run `yubioath-desktop` and insert your YubiKey when prompted.

## Troubleshooting

Restart, especially if you have completed updates since your YubiKey last worked. Do this even if some functions appear to be functioning.

### YubiKey not acting as HID device

> [!note]
>
> These steps are no longer necessary after [systemd since v244](https://github.com/systemd/systemd/commit/d45ee2f31a8358db0accde2e7c81777cedadc3c2) added native support for this functionality.

Add udev rule as described in [this article](https://michaelheap.com/yubikey-on-arch/):

```bash
# /etc/udev/rules.d/10-security-key.rules
KERNEL=="hidraw*", SUBSYSTEM=="hidraw", MODE="0664", GROUP="users", ATTRS{idVendor}=="2581", ATTRS{idProduct}=="f1d0"
```

Run `udevadm trigger` afterwards.

### ykman fails to connect to the YubiKey

If the manager fails to connect to the YubiKey, make sure you have started `pcscd.service` or `pcscd.socket`.

### Error: Failed connecting to YubiKey 5 [OTP+FIDO+CCID]. Make sure the application have the required permissions.

This can occur when using `ykman` to access the oath credentials on the device if `scdaemon` has already taken exclusive control of the device. [https://github.com/Yubico/yubikey-manager/issues/35]

To fix this you can set the `reader-port` option with the correct value for your device in `~/.gnupg/scdaemon.conf`. [https://support.yubico.com/hc/en-us/articles/360013714479-Troubleshooting-Issues-with-GPG]

> [!note]
>
> This will cause the gpgagent to re-prompt you to unlock the YubiKey after each time you access the YubiKey through ykman.

For YubiKey NEO and YubiKey 4:

```
reader-port Yubico Yubikey
```

or for YubiKey 5:

```
reader-port Yubico Yubi
```

### YubiKey fails to bind within a guest VM

Assuming the YubiKey is available to the guest, the issue results from a driver binding to the device on the host. To unbind the device, the bus and port information is needed from dmesg on the host:

```bash
# dmesg | grep -B1 Yubico | tail -n 2 | head -n 1 | sed -E 's/^\[[^]]+\] usb ([^:]*):.*/\1/'
```

The resulting USB id should be of the form `X-Y.Z` or `X-Y`. Then, on the host, use `find` to search `/sys/bus/usb/drivers` for which driver the YubiKey is binded to (e.g. `usbhid` or `usbfs`).

```bash
$ find /sys/bus/usb/drivers -name "*X-Y.Z*"
```

To unbind the device, use the result from the previous command (i.e. `/sys/bus/usb/drivers/DRIVER/X-Y.Z:1.0`):

```bash
# echo 'X-Y.Z:1.0' > /sys/bus/usb/drivers/DRIVER/unbind
```

### Error: [key] could not be locally signed or gpg: No default secret key: No public key

Occurs when attempting to sign keys on a non-standard keyring while a YubiKey is plugged in, e.g. as Pacman does in `pacman-key --populate`. The solution is to remove the offending YubiKey and start over.

### YubiKey disappears and reappears in Yubico Authenticator

This happens when the CCID driver is not installed. You may need to install the `ccid` package.

### YubiKey core error: timeout

You are probably using the wrong slot. Try the other one.

### gpg: no such device

Because gpg(scdaemon) tries to acquire exclusive access to the yubikey. It needs to be configured to use pscs and use shared access.[https://blog.apdu.fr/posts/2024/12/gnupg-and-pcsc-conflicts-episode-3/][https://github.com/LudovicRousseau/PCSC/issues/65]

Your file `~/.gnupg/scdaemon.conf` should contain:

```
disable-ccid
pcsc-shared
```

For old versions of GnuPG, the `pcsc-shared` option is not available. Only keep `disable-ccid` and restart `pcscd.service` as a workaround.
